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

func ScrapeEventByID(id string) (*Event, error) {
	url := "http://ufcstats.com/event-details/" + id

	html, err := fetchHTML(url)
	if err != nil {
		return nil, err
	}

	return parseEventDetail(html, url)
}

func ScrapeEventByURL(url string) (*Event, error) {
	html, err := fetchHTML(url)
	if err != nil {
		return nil, err
	}

	return parseEventDetail(html, url)
}
