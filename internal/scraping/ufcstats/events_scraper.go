package ufcstats

func ScrapeEvents() ([]Event, error) {

	html, err := fetchHTML(
		"http://ufcstats.com/statistics/events/completed?page=all",
	)

	if err != nil {
		return nil, err
	}

	return parseEvents(html)
}
