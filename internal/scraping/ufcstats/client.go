package ufcstats

import (
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	"ufc-backend/internal/shared/logger"
)

func fetchHTML(
	url string,
) (string, error) {

	if strings.TrimSpace(url) == "" {
		return "", fmt.Errorf("%w: empty", ErrInvalidURL)
	}

	parsed, err := neturl.ParseRequestURI(url)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("%w: %s", ErrInvalidURL, url)
	}

	start := time.Now()
	logger.Debugf("ufcstats_request method=GET url=%s", url)

	response, err := http.Get(url)

	if err != nil {
		logger.Errorf("ufcstats_request_failed method=GET url=%s err=%s", url, err.Error())
		return "", fmt.Errorf("%w: %s", ErrRequestFailed, err.Error())
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		logger.Errorf("ufcstats_unexpected_status method=GET url=%s status=%d", url, response.StatusCode)
		return "", fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, response.StatusCode)
	}

	body, err := io.ReadAll(
		response.Body,
	)

	if err != nil {
		logger.Errorf("ufcstats_read_body_failed method=GET url=%s err=%s", url, err.Error())
		return "", fmt.Errorf("%w: %s", ErrReadResponseBody, err.Error())
	}

	logger.Debugf("ufcstats_response method=GET url=%s status=%d bytes=%d duration_ms=%d", url, response.StatusCode, len(body), time.Since(start).Milliseconds())
	return string(body), nil
}
