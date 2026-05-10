package ufcstats

import "strings"

func ScrapeFightURLsForEvent(eventURL string) ([]string, error) {
	html, err := fetchHTML(eventURL)
	if err != nil {
		return nil, err
	}

	return parseEventFightURLs(html)
}

func ScrapeFightByURL(url string) (*Fight, *Fighter, *Fighter, error) {
	html, err := fetchHTML(url)
	if err != nil {
		return nil, nil, nil, err
	}

	fight, red, blue, err := parseFightDetail(html, url)
	if err != nil {
		return nil, nil, nil, err
	}

	if red != nil {
		red.URL = strings.TrimSpace(red.URL)
	}
	if blue != nil {
		blue.URL = strings.TrimSpace(blue.URL)
	}

	return fight, red, blue, nil
}
