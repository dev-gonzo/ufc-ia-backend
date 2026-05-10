package tapology

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/websocket"
)

var ErrMissingScrapingBrowserWS = errors.New("missing scraping browser websocket url")

func ScrapePromotionEvents(
	ctx context.Context,
	wsURL string,
	promotionURL string,
) ([]Event, error) {
	log.Printf("tapology: starting scrape promotion events url=%s", promotionURL)

	html, err := fetchRenderedHTML(
		ctx,
		wsURL,
		promotionURL,
	)
	if err != nil {
		log.Printf("tapology: fetch rendered html failed url=%s err=%s", promotionURL, redactSecrets(err.Error()))
		return nil, err
	}

	log.Printf("tapology: fetched html bytes=%d url=%s", len(html), promotionURL)

	events, err := parsePromotionEvents(html)
	if err != nil {
		log.Printf("tapology: parse failed url=%s err=%s", promotionURL, redactSecrets(err.Error()))
		return nil, err
	}

	if len(events) == 0 {
		snippet := html
		if len(snippet) > 400 {
			snippet = snippet[:400]
		}
		log.Printf("tapology: parsed 0 events url=%s html_prefix=%q", promotionURL, snippet)
	} else {
		first := events[0]
		log.Printf("tapology: parsed events count=%d first_name=%q first_url=%s", len(events), first.Name, first.URL)
	}

	return events, nil
}

func fetchRenderedHTML(
	ctx context.Context,
	wsURL string,
	targetURL string,
) (string, error) {
	wsURL = strings.TrimSpace(wsURL)
	if wsURL == "" {
		return "", ErrMissingScrapingBrowserWS
	}

	parsed, err := url.Parse(wsURL)
	if err != nil {
		log.Printf("tapology: invalid scraping browser url err=%s", redactSecrets(err.Error()))
		return "", err
	}

	log.Printf("tapology: allocator url scheme=%s host=%s path=%s", parsed.Scheme, parsed.Host, parsed.Path)
	log.Printf("tapology: connecting to remote browser (ws redacted) target=%s", targetURL)

	ctx, cancelTimeout := context.WithTimeout(ctx, 4*time.Minute)
	defer cancelTimeout()

	client, err := dialCDP(ctx, parsed)
	if err != nil {
		return "", err
	}
	defer client.Close()

	target, err := client.createTarget(ctx, "about:blank")
	if err != nil {
		return "", err
	}

	sessionID, err := client.attachToTarget(ctx, target, true)
	if err != nil {
		return "", err
	}

	_ = client.call(ctx, sessionID, "Page.enable", nil, nil)

	if err := client.call(ctx, sessionID, "Page.navigate", map[string]any{"url": targetURL}, nil); err != nil {
		return "", err
	}

	var captcha struct {
		Status string `json:"status"`
	}
	if err := client.call(
		ctx,
		sessionID,
		"Captcha.waitForSolve",
		map[string]any{"detectTimeout": 30 * 1000},
		&captcha,
	); err != nil {
		log.Printf("tapology: Captcha.waitForSolve failed err=%s", redactSecrets(err.Error()))
	} else if captcha.Status != "" {
		log.Printf("tapology: captcha status=%q", captcha.Status)
	}

	deadline := time.Now().Add(90 * time.Second)
	var html string
	for {
		if time.Now().After(deadline) {
			if len(html) > 0 {
				snippet := html
				if len(snippet) > 400 {
					snippet = snippet[:400]
				}
				log.Printf("tapology: html prefix=%q", snippet)
			}
			return "", context.DeadlineExceeded
		}

		var eval struct {
			Result struct {
				Value string `json:"value"`
			} `json:"result"`
		}
		if err := client.call(
			ctx,
			sessionID,
			"Runtime.evaluate",
			map[string]any{"expression": "document.documentElement.outerHTML", "returnByValue": true},
			&eval,
		); err != nil {
			return "", err
		}

		html = eval.Result.Value
		if strings.Contains(html, "/fightcenter/events/") {
			break
		}

		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(3 * time.Second):
		}
	}

	return html, nil
}

func parsePromotionEvents(
	html string,
) ([]Event, error) {
	document, err := goquery.NewDocumentFromReader(
		strings.NewReader(html),
	)
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	var events []Event

	document.Find("a").Each(func(_ int, selection *goquery.Selection) {
		href, ok := selection.Attr("href")
		if !ok {
			return
		}

		if !strings.HasPrefix(href, "/fightcenter/events/") && !strings.Contains(href, "/fightcenter/events/") {
			return
		}

		name := normalizeWhitespace(selection.Text())
		if name == "" {
			return
		}

		fullURL := href
		if strings.HasPrefix(fullURL, "/") {
			fullURL = "https://www.tapology.com" + fullURL
		}

		if seen[fullURL] {
			return
		}
		seen[fullURL] = true

		container := selection.ParentsFiltered("li, tr, div").First()
		details := ""
		if container.Length() > 0 {
			details = normalizeWhitespace(container.Text())
		}
		if details == "" {
			details = name
		}

		events = append(events, Event{
			Name:    name,
			URL:     fullURL,
			Details: details,
		})
	})

	return events, nil
}

func normalizeWhitespace(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

var credentialsInURL = regexp.MustCompile(`(?i)(wss?|https?)://[^/@\s]+@`)

func redactSecrets(message string) string {
	if message == "" {
		return message
	}

	redacted := credentialsInURL.ReplaceAllStringFunc(message, func(m string) string {
		idx := strings.Index(m, "://")
		if idx == -1 {
			return m
		}
		scheme := m[:idx+3]
		return fmt.Sprintf("%s***@", scheme)
	})

	return redacted
}

type cdpClient struct {
	connection *websocket.Conn

	nextID int64

	mu      sync.Mutex
	pending map[int64]chan cdpIncoming

	events chan cdpIncoming
}

type cdpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type cdpIncoming struct {
	ID        *int64          `json:"id,omitempty"`
	Method    string          `json:"method,omitempty"`
	Params    json.RawMessage `json:"params,omitempty"`
	Result    json.RawMessage `json:"result,omitempty"`
	Error     *cdpError       `json:"error,omitempty"`
	SessionID string          `json:"sessionId,omitempty"`
}

type cdpOutgoing struct {
	ID        int64 `json:"id"`
	Method    string `json:"method"`
	Params    any    `json:"params,omitempty"`
	SessionID string `json:"sessionId,omitempty"`
}

func dialCDP(ctx context.Context, parsed *url.URL) (*cdpClient, error) {
	headers := http.Header{}

	username := ""
	password := ""
	if parsed.User != nil {
		username = parsed.User.Username()
		password, _ = parsed.User.Password()
	}

	if username != "" {
		creds := username
		if password != "" {
			creds = creds + ":" + password
		}
		headers.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(creds)))
	}

	dialURL := *parsed
	dialURL.User = nil

	dialer := websocket.Dialer{
		HandshakeTimeout: 20 * time.Second,
	}

	wsConn, resp, err := dialer.DialContext(ctx, dialURL.String(), headers)
	if err != nil {
		status := ""
		if resp != nil {
			status = resp.Status
		}
		return nil, fmt.Errorf("cdp websocket dial failed (status=%s): %w", status, err)
	}

	client := &cdpClient{
		connection: wsConn,
		pending:    map[int64]chan cdpIncoming{},
		events:     make(chan cdpIncoming, 128),
	}

	go client.readLoop()
	return client, nil
}

func (c *cdpClient) Close() error {
	return c.connection.Close()
}

func (c *cdpClient) readLoop() {
	for {
		_, data, err := c.connection.ReadMessage()
		if err != nil {
			return
		}

		var msg cdpIncoming
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		if msg.ID != nil {
			c.mu.Lock()
			ch := c.pending[*msg.ID]
			if ch != nil {
				delete(c.pending, *msg.ID)
			}
			c.mu.Unlock()

			if ch != nil {
				ch <- msg
				close(ch)
			}
			continue
		}

		select {
		case c.events <- msg:
		default:
		}
	}
}

func (c *cdpClient) call(ctx context.Context, sessionID string, method string, params any, out any) error {
	id := atomic.AddInt64(&c.nextID, 1)
	responseCh := make(chan cdpIncoming, 1)

	c.mu.Lock()
	c.pending[id] = responseCh
	c.mu.Unlock()

	payload := cdpOutgoing{
		ID:     id,
		Method: method,
		Params: params,
	}
	if sessionID != "" {
		payload.SessionID = sessionID
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := c.connection.WriteMessage(websocket.TextMessage, body); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case resp := <-responseCh:
		if resp.Error != nil {
			return fmt.Errorf("%s: %s", method, resp.Error.Message)
		}
		if out != nil && len(resp.Result) > 0 {
			return json.Unmarshal(resp.Result, out)
		}
		return nil
	}
}

func (c *cdpClient) waitForEvent(ctx context.Context, sessionID string, method string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev := <-c.events:
			if ev.Method != method {
				continue
			}
			if sessionID != "" && ev.SessionID != sessionID {
				continue
			}
			return nil
		}
	}
}

func (c *cdpClient) createTarget(ctx context.Context, targetURL string) (string, error) {
	var resp struct {
		TargetID string `json:"targetId"`
	}
	if err := c.call(ctx, "", "Target.createTarget", map[string]any{"url": targetURL}, &resp); err != nil {
		return "", err
	}
	if resp.TargetID == "" {
		return "", errors.New("missing targetId from Target.createTarget")
	}
	return resp.TargetID, nil
}

func (c *cdpClient) attachToTarget(ctx context.Context, targetID string, flatten bool) (string, error) {
	var resp struct {
		SessionID string `json:"sessionId"`
	}
	if err := c.call(ctx, "", "Target.attachToTarget", map[string]any{"targetId": targetID, "flatten": flatten}, &resp); err != nil {
		return "", err
	}
	if resp.SessionID == "" {
		return "", errors.New("missing sessionId from Target.attachToTarget")
	}
	return resp.SessionID, nil
}
