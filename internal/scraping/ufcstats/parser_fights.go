package ufcstats

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseEventFightURLs(
	html string,
) ([]string, error) {
	document, err := goquery.NewDocumentFromReader(
		strings.NewReader(html),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParseFailed, err.Error())
	}

	seen := map[string]bool{}
	var urls []string

	document.Find("a").Each(func(_ int, selection *goquery.Selection) {
		href, ok := selection.Attr("href")
		if !ok {
			return
		}
		if !strings.Contains(href, "/fight-details/") {
			return
		}
		href = strings.TrimSpace(href)
		if href == "" {
			return
		}
		if seen[href] {
			return
		}
		seen[href] = true
		urls = append(urls, href)
	})

	if len(urls) == 0 {
		return nil, fmt.Errorf("%w: no fight urls found", ErrMissingRequiredField)
	}

	return urls, nil
}

func parseFightDetail(
	html string,
	url string,
) (*Fight, *Fighter, *Fighter, error) {
	document, err := goquery.NewDocumentFromReader(
		strings.NewReader(html),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: %s", ErrParseFailed, err.Error())
	}

	fight := &Fight{
		URL: url,
	}

	title := normalizeWhitespace(document.Find(".b-fight-details__fight-title").First().Text())
	if title != "" {
		title = strings.TrimSpace(strings.TrimSuffix(title, "Bout"))
		fight.WeightClass = title
	}

	var red *Fighter
	var blue *Fighter

	document.Find(".b-fight-details__persons .b-fight-details__person").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Find(".b-fight-details__person-name a").First().Text())
		href, _ := s.Find(".b-fight-details__person-name a").First().Attr("href")
		status := strings.TrimSpace(s.Find(".b-fight-details__person-status").First().Text())

		f := &Fighter{
			Name:   name,
			URL:    strings.TrimSpace(href),
			Record: "",
		}

		if i == 0 {
			red = f
			if status == "W" {
				fight.Winner = "red"
			}
		}
		if i == 1 {
			blue = f
			if status == "W" {
				fight.Winner = "blue"
			}
		}
	})

	document.Find(".b-fight-details__text-item, .b-fight-details__text-item_first").Each(func(_ int, s *goquery.Selection) {
		t := normalizeWhitespace(s.Text())
		switch {
		case strings.HasPrefix(t, "Weight:"):
			if strings.TrimSpace(fight.WeightClass) == "" {
				fight.WeightClass = strings.TrimSpace(strings.TrimPrefix(t, "Weight:"))
			}
		case strings.HasPrefix(t, "Method:"):
			fight.Method = strings.TrimSpace(strings.TrimPrefix(t, "Method:"))
		case strings.HasPrefix(t, "Round:"):
			roundStr := strings.TrimSpace(strings.TrimPrefix(t, "Round:"))
			if n, err := strconv.Atoi(roundStr); err == nil {
				fight.Round = n
			}
		case strings.HasPrefix(t, "Time:"):
			fight.Time = strings.TrimSpace(strings.TrimPrefix(t, "Time:"))
		}
	})

	if red == nil || blue == nil {
		return nil, nil, nil, fmt.Errorf("%w: fighters not found", ErrMissingRequiredField)
	}

	return fight, red, blue, nil
}

func normalizeWhitespace(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}
