package ufcstats

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseEvents(
	html string,
) ([]Event, error) {

	document, err := goquery.NewDocumentFromReader(
		strings.NewReader(html),
	)

	if err != nil {
		return nil, err
	}

	var events []Event

	document.Find(
		"tr.b-statistics__table-row",
	).Each(func(i int, selection *goquery.Selection) {

		link := selection.Find("a").First()

		name := strings.TrimSpace(
			link.Text(),
		)

		url, _ := link.Attr("href")

		date := strings.TrimSpace(
			selection.Find("td").
				Eq(1).
				Text(),
		)

		location := strings.TrimSpace(
			selection.Find("td").
				Eq(2).
				Text(),
		)

		if name == "" {
			return
		}

		events = append(
			events,
			Event{
				Name:     name,
				URL:      url,
				Date:     date,
				Location: location,
			},
		)
	})

	return events, nil
}
