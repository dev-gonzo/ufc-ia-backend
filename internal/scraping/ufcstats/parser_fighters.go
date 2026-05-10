package ufcstats

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseFighterDetail(
	html string,
	url string,
) (*Fighter, error) {
	document, err := goquery.NewDocumentFromReader(
		strings.NewReader(html),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParseFailed, err.Error())
	}

	name := strings.TrimSpace(document.Find(".b-content__title-highlight").First().Text())
	record := normalizeWhitespace(document.Find(".b-content__title-record").First().Text())
	record = strings.TrimSpace(strings.TrimPrefix(record, "Record:"))
	nicknameText := strings.TrimSpace(document.Find(".b-content__Nickname").First().Text())

	fighter := &Fighter{
		Name:   name,
		URL:    url,
		Record: record,
	}

	if nicknameText != "" {
		fighter.Nickname = &nicknameText
	}

	document.Find(".b-list__box-list-item").Each(func(_ int, s *goquery.Selection) {
		t := normalizeWhitespace(s.Text())
		switch {
		case strings.HasPrefix(t, "Nickname:"):
			v := strings.TrimSpace(strings.TrimPrefix(t, "Nickname:"))
			if v != "" {
				fighter.Nickname = &v
			}
		case strings.HasPrefix(t, "Height:"):
			v := strings.TrimSpace(strings.TrimPrefix(t, "Height:"))
			if v != "" {
				fighter.Height = &v
			}
		case strings.HasPrefix(t, "Weight:"):
			v := strings.TrimSpace(strings.TrimPrefix(t, "Weight:"))
			if v != "" {
				fighter.Weight = &v
			}
		case strings.HasPrefix(t, "Reach:"):
			v := strings.TrimSpace(strings.TrimPrefix(t, "Reach:"))
			if v != "" {
				fighter.Reach = &v
			}
		case strings.HasPrefix(t, "STANCE:"):
			v := strings.TrimSpace(strings.TrimPrefix(t, "STANCE:"))
			if v != "" {
				fighter.Stance = &v
			}
		case strings.HasPrefix(t, "DOB:"):
			v := strings.TrimSpace(strings.TrimPrefix(t, "DOB:"))
			if v != "" {
				fighter.DOB = &v
			}
		}
	})

	if strings.TrimSpace(fighter.Name) == "" {
		return nil, fmt.Errorf("%w: fighter name not found", ErrMissingRequiredField)
	}

	return fighter, nil
}
