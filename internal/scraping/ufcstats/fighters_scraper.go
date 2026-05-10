package ufcstats

import "strings"

func ScrapeFighterByURL(url string) (*Fighter, error) {
	html, err := fetchHTML(url)
	if err != nil {
		return nil, err
	}

	fighter, err := parseFighterDetail(html, url)
	if err != nil {
		return nil, err
	}

	fighter.Name = strings.TrimSpace(fighter.Name)
	fighter.Record = strings.TrimSpace(fighter.Record)
	return fighter, nil
}
