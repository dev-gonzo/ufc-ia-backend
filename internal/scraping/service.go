package scraping

import (
	"context"
	"log"
	"os"
	"ufc-backend/internal/scraping/tapology"
	"ufc-backend/internal/scraping/ufcstats"
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
	events, err := ufcstats.ScrapeEvents()
	if err != nil {
		return nil, err
	}

	err = s.repo.UpsertEvents(ctx, events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Service) ScrapeAndSaveEventByID(ctx context.Context, id string) (*ufcstats.Event, error) {
	event, err := ufcstats.ScrapeEventByID(id)
	if err != nil {
		return nil, err
	}

	err = s.repo.UpsertEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *Service) GetEvent(ctx context.Context, id string) (*ufcstats.Event, error) {
	return s.repo.GetEventByID(ctx, id)
}

func (s *Service) ScrapeAndSaveTapologyUFCEvents(ctx context.Context) ([]tapology.Event, error) {
	wsURL := os.Getenv("SCRAPING_BROWSER_WS_URL")
	promotionURL := "https://www.tapology.com/fightcenter/promotions/1-ultimate-fighting-championship-ufc"

	events, err := tapology.ScrapePromotionEvents(
		ctx,
		wsURL,
		promotionURL,
	)
	if err != nil {
		return nil, err
	}

	err = s.repo.UpsertTapologyEvents(ctx, events)
	if err != nil {
		log.Printf("tapology: db upsert failed err=%s", err.Error())
		return nil, err
	}

	return events, nil
}
