package scraping

import (
	"context"
	"errors"
	"os"
	"strings"
	"ufc-backend/internal/scraping/tapology"
	"ufc-backend/internal/scraping/ufcstats"
	"ufc-backend/internal/shared/logger"

	"github.com/jackc/pgx/v5"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ScrapeAndSaveEvents(ctx context.Context) ([]ufcstats.Event, error) {
	logger.Debugf("scrape_events_start")
	events, err := ufcstats.ScrapeEvents()
	if err != nil {
		logger.Errorf("scrape_events_scrape_failed err=%s", err.Error())
		return nil, err
	}

	err = s.repo.UpsertEvents(ctx, events)
	if err != nil {
		logger.Errorf("scrape_events_upsert_failed err=%s", err.Error())
		return nil, err
	}

	logger.Debugf("scrape_events_done count=%d", len(events))
	return events, nil
}

func (s *Service) ScrapeAndSaveEventByID(ctx context.Context, id string) (*ufcstats.Event, error) {
	logger.Debugf("scrape_event_start id=%s", strings.TrimSpace(id))
	event, err := ufcstats.ScrapeEventByID(id)
	if err != nil {
		logger.Errorf("scrape_event_scrape_failed id=%s err=%s", strings.TrimSpace(id), err.Error())
		return nil, err
	}

	err = s.repo.UpsertEvent(ctx, event)
	if err != nil {
		logger.Errorf("scrape_event_upsert_failed id=%s err=%s", strings.TrimSpace(id), err.Error())
		return nil, err
	}

	_, err = s.ScrapeAndSaveEventFights(ctx, "", event.URL)
	if err != nil {
		logger.Errorf("scrape_event_fights_after_event_failed id=%s url=%s err=%s", strings.TrimSpace(id), strings.TrimSpace(event.URL), err.Error())
		return nil, err
	}

	logger.Debugf("scrape_event_done id=%s url=%s", strings.TrimSpace(id), strings.TrimSpace(event.URL))
	return event, nil
}

func (s *Service) GetEvent(ctx context.Context, id string) (*ufcstats.Event, error) {
	return s.repo.GetEventByID(ctx, id)
}

func (s *Service) ScrapeAndSaveTapologyUFCEvents(ctx context.Context) ([]tapology.Event, error) {
	wsURL := os.Getenv("SCRAPING_BROWSER_WS_URL")
	promotionURL := "https://www.tapology.com/fightcenter/promotions/1-ultimate-fighting-championship-ufc"

	logger.Debugf("tapology_events_start promotion_url=%s", promotionURL)
	events, err := tapology.ScrapePromotionEvents(
		ctx,
		wsURL,
		promotionURL,
	)
	if err != nil {
		logger.Errorf("tapology_events_scrape_failed err=%s", err.Error())
		return nil, err
	}

	err = s.repo.UpsertTapologyEvents(ctx, events)
	if err != nil {
		logger.Errorf("tapology_events_upsert_failed err=%s", err.Error())
		return nil, err
	}

	logger.Debugf("tapology_events_done count=%d", len(events))
	return events, nil
}

func (s *Service) ScrapeAndSaveEventFights(
	ctx context.Context,
	id string,
	url string,
) ([]ufcstats.Fight, error) {
	var event *ufcstats.Event
	var err error

	if id != "" {
		logger.Debugf("scrape_event_fights_load_event_by_id id=%s", strings.TrimSpace(id))
		event, err = s.repo.GetEventByID(ctx, id)
		if err != nil {
			logger.Errorf("scrape_event_fights_get_event_by_id_failed id=%s err=%s", strings.TrimSpace(id), err.Error())
			return nil, err
		}
	} else {
		logger.Debugf("scrape_event_fights_load_event_by_url url=%s", strings.TrimSpace(url))
		event, err = s.repo.GetEventByURL(ctx, url)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				logger.Debugf("scrape_event_fights_event_not_found_scrape_event url=%s", strings.TrimSpace(url))
				scrapedEvent, scrapeErr := ufcstats.ScrapeEventByURL(url)
				if scrapeErr != nil {
					logger.Errorf("scrape_event_fights_scrape_event_failed url=%s err=%s", strings.TrimSpace(url), scrapeErr.Error())
					return nil, scrapeErr
				}

				if err := s.repo.UpsertEvent(ctx, scrapedEvent); err != nil {
					logger.Errorf("scrape_event_fights_upsert_event_failed url=%s err=%s", strings.TrimSpace(url), err.Error())
					return nil, err
				}

				event, err = s.repo.GetEventByURL(ctx, url)
				if err != nil {
					logger.Errorf("scrape_event_fights_get_event_by_url_failed url=%s err=%s", strings.TrimSpace(url), err.Error())
					return nil, err
				}
			} else {
				logger.Errorf("scrape_event_fights_get_event_by_url_failed url=%s err=%s", strings.TrimSpace(url), err.Error())
				return nil, err
			}
		}
	}

	eventID := ""
	if event.ID != nil {
		eventID = *event.ID
	}

	logger.Debugf("scrape_event_fights_scrape_fight_urls event_url=%s", strings.TrimSpace(event.URL))
	fightURLs, err := ufcstats.ScrapeFightURLsForEvent(event.URL)
	if err != nil {
		logger.Errorf("scrape_event_fights_scrape_fight_urls_failed event_url=%s err=%s", strings.TrimSpace(event.URL), err.Error())
		return nil, err
	}

	var fights []ufcstats.Fight

	for _, fightURL := range fightURLs {
		fight, redRef, blueRef, err := ufcstats.ScrapeFightByURL(fightURL)
		if err != nil {
			logger.Errorf("scrape_event_fights_scrape_fight_failed fight_url=%s err=%s", strings.TrimSpace(fightURL), err.Error())
			return nil, err
		}

		if redRef != nil && strings.TrimSpace(redRef.URL) != "" {
			redFighter, err := ufcstats.ScrapeFighterByURL(redRef.URL)
			if err != nil {
				logger.Errorf("scrape_event_fights_scrape_fighter_failed fighter_url=%s err=%s", strings.TrimSpace(redRef.URL), err.Error())
				return nil, err
			}
			redID, err := s.repo.UpsertFighter(ctx, redFighter)
			if err != nil {
				logger.Errorf("scrape_event_fights_upsert_fighter_failed fighter_url=%s err=%s", strings.TrimSpace(redRef.URL), err.Error())
				return nil, err
			}
			fight.RedFighterID = &redID
		}

		if blueRef != nil && strings.TrimSpace(blueRef.URL) != "" {
			blueFighter, err := ufcstats.ScrapeFighterByURL(blueRef.URL)
			if err != nil {
				logger.Errorf("scrape_event_fights_scrape_fighter_failed fighter_url=%s err=%s", strings.TrimSpace(blueRef.URL), err.Error())
				return nil, err
			}
			blueID, err := s.repo.UpsertFighter(ctx, blueFighter)
			if err != nil {
				logger.Errorf("scrape_event_fights_upsert_fighter_failed fighter_url=%s err=%s", strings.TrimSpace(blueRef.URL), err.Error())
				return nil, err
			}
			fight.BlueFighterID = &blueID
		}

		fight.EventID = &eventID

		_, err = s.repo.UpsertFight(ctx, fight)
		if err != nil {
			logger.Errorf("scrape_event_fights_upsert_fight_failed fight_url=%s err=%s", strings.TrimSpace(fight.URL), err.Error())
			return nil, err
		}

		fights = append(fights, *fight)
	}

	if strings.TrimSpace(eventID) != "" {
		if err := s.repo.MarkEventSynced(ctx, eventID); err != nil {
			logger.Errorf("scrape_event_fights_mark_synced_failed event_id=%s err=%s", strings.TrimSpace(eventID), err.Error())
			return nil, err
		}
	}

	logger.Debugf("scrape_event_fights_done event_id=%s count=%d", strings.TrimSpace(eventID), len(fights))
	return fights, nil
}

func (s *Service) ScrapeAndSaveFighter(
	ctx context.Context,
	id string,
	url string,
) (*ufcstats.Fighter, error) {
	if strings.TrimSpace(url) == "" && strings.TrimSpace(id) != "" {
		url = "http://ufcstats.com/fighter-details/" + strings.TrimSpace(id)
	}

	logger.Debugf("scrape_fighter_start url=%s", strings.TrimSpace(url))
	fighter, err := ufcstats.ScrapeFighterByURL(url)
	if err != nil {
		logger.Errorf("scrape_fighter_scrape_failed url=%s err=%s", strings.TrimSpace(url), err.Error())
		return nil, err
	}

	fighterID, err := s.repo.UpsertFighter(ctx, fighter)
	if err != nil {
		logger.Errorf("scrape_fighter_upsert_failed url=%s err=%s", strings.TrimSpace(url), err.Error())
		return nil, err
	}

	fighter.ID = &fighterID
	logger.Debugf("scrape_fighter_done id=%s url=%s", strings.TrimSpace(fighterID), strings.TrimSpace(url))
	return fighter, nil
}
