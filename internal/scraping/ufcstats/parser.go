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
		name := strings.TrimSpace(link.Text())
		url, _ := link.Attr("href")

		// No ufcstats.com/statistics/events/completed:
		// td[0] tem o nome e o link, e também a data em um span
		// td[1] tem a localização

		date := strings.TrimSpace(
			selection.Find("td").Eq(0).Find("span").Text(),
		)

		location := strings.TrimSpace(
			selection.Find("td").Eq(1).Text(),
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

func parseEventDetail(html string, url string) (*Event, error) {
	document, err := goquery.NewDocumentFromReader(
		strings.NewReader(html),
	)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(
		document.Find(".b-content__title-highlight").Text(),
	)

	dateStr := strings.TrimSpace(
		document.Find(".b-list__box-list-item").Eq(0).Text(),
	)
	dateStr = strings.TrimSpace(strings.Replace(dateStr, "Date:", "", 1))

	locationStr := strings.TrimSpace(
		document.Find(".b-list__box-list-item").Eq(1).Text(),
	)
	locationStr = strings.TrimSpace(strings.Replace(locationStr, "Location:", "", 1))

	return &Event{
		Name:     name,
		URL:      url,
		Date:     dateStr,
		Location: locationStr,
	}, nil
}
