package ufcstats

func ScrapeFightDetailsByURL(url string) (*FightDetailsScrape, error) {
	html, err := fetchHTML(url)
	if err != nil {
		return nil, err
	}

	return parseFightDetails(html, url)
}
