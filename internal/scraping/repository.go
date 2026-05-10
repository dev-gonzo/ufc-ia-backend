package scraping

import (
	"context"
	"ufc-backend/internal/scraping/tapology"
	"ufc-backend/internal/scraping/ufcstats"

	"time"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) UpsertEvents(ctx context.Context, events []ufcstats.Event) error {
	for _, event := range events {
		if err := r.UpsertEvent(ctx, &event); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) UpsertEvent(ctx context.Context, event *ufcstats.Event) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO events (name, url, date, location, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (url) DO UPDATE SET
			name = EXCLUDED.name,
			date = EXCLUDED.date,
			location = EXCLUDED.location,
			updated_at = NOW()
	`, event.Name, event.URL, event.Date, event.Location)
	return err
}

func (r *Repository) UpsertTapologyEvents(ctx context.Context, events []tapology.Event) error {
	for _, event := range events {
		if err := r.UpsertTapologyEvent(ctx, &event); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) UpsertTapologyEvent(ctx context.Context, event *tapology.Event) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO events_tapology (name, url, details, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (url) DO UPDATE SET
			name = EXCLUDED.name,
			details = EXCLUDED.details,
			updated_at = NOW()
	`, event.Name, event.URL, event.Details)
	return err
}

func (r *Repository) GetEventByURL(ctx context.Context, url string) (*ufcstats.Event, error) {
	var (
		e         ufcstats.Event
		id        string
		eventSync bool
		createdAt time.Time
		updatedAt time.Time
	)

	err := r.db.QueryRow(ctx, `
		SELECT id, name, url, date, location, event_sync, created_at, updated_at
		FROM events
		WHERE url = $1
	`, url).Scan(&id, &e.Name, &e.URL, &e.Date, &e.Location, &eventSync, &createdAt, &updatedAt)

	if err != nil {
		return nil, err
	}

	e.ID = &id
	e.EventSync = &eventSync
	e.CreatedAt = &createdAt
	e.UpdatedAt = &updatedAt

	return &e, nil
}

func (r *Repository) GetEventByID(ctx context.Context, id string) (*ufcstats.Event, error) {
	var (
		e         ufcstats.Event
		eventID   string
		eventSync bool
		createdAt time.Time
		updatedAt time.Time
	)

	err := r.db.QueryRow(ctx, `
		SELECT id, name, url, date, location, event_sync, created_at, updated_at
		FROM events
		WHERE id = $1
	`, id).Scan(&eventID, &e.Name, &e.URL, &e.Date, &e.Location, &eventSync, &createdAt, &updatedAt)

	if err != nil {
		return nil, err
	}

	e.ID = &eventID
	e.EventSync = &eventSync
	e.CreatedAt = &createdAt
	e.UpdatedAt = &updatedAt

	return &e, nil
}
