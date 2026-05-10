package ufc

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"ufc-backend/internal/shared/logger"
)

func ScrapeAthleteDetailsByName(ctx context.Context, name string) (*AthleteDetails, error) {
	slug, err := AthleteSlugFromName(name)
	if err != nil {
		return nil, err
	}

	return ScrapeAthleteDetailsBySlug(ctx, slug)
}

func ScrapeAthleteDetailsBySlug(ctx context.Context, slug string) (*AthleteDetails, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil, ErrInvalidAthleteSlug
	}

	athleteURL := "https://www.ufc.com.br/athlete/" + slug
	html, status, err := fetchHTML(ctx, athleteURL)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, ErrAthleteNotFound
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, status)
	}

	details, imgURL, err := parseAthleteDetails(html, athleteURL, slug)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(imgURL) != "" {
		photo, err := fetchAndConvertToWebPBase64(ctx, imgURL)
		if err != nil {
			logger.Errorf("ufc_photo_convert_failed athlete_url=%s img_url=%s err=%s", athleteURL, imgURL, err.Error())
		} else if strings.TrimSpace(photo) != "" {
			details.PhotoWebPBase64 = &photo
		}
	}

	return details, nil
}

func fetchHTML(ctx context.Context, url string) (string, int, error) {
	client := &http.Client{Timeout: 25 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	start := time.Now()
	logger.Debugf("ufc_request method=GET url=%s", url)
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, err
	}

	logger.Debugf("ufc_response method=GET url=%s status=%d bytes=%d duration_ms=%d", url, resp.StatusCode, len(body), time.Since(start).Milliseconds())
	return string(body), resp.StatusCode, nil
}

func fetchAndConvertToWebPBase64(ctx context.Context, url string) (string, error) {
	client := &http.Client{Timeout: 25 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "image/webp,image/*;q=0.8,*/*;q=0.5")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	contentType := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Type")))
	if !strings.HasPrefix(contentType, "image/webp") {
		return "", fmt.Errorf("%w: expected webp got %s", ErrUnexpectedStatusCode, contentType)
	}

	encoded := base64.StdEncoding.EncodeToString(b)
	return encoded, nil
}
