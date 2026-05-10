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

func (r *Repository) MarkEventSynced(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE events
		SET event_sync = true, updated_at = NOW()
		WHERE id = $1
	`, id)
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

func (r *Repository) UpsertFighter(ctx context.Context, fighter *ufcstats.Fighter) (string, error) {
	var id string

	err := r.db.QueryRow(ctx, `
		INSERT INTO fighters (name, url, record, nickname, height, weight, reach, stance, dob, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		ON CONFLICT (url) DO UPDATE SET
			name = EXCLUDED.name,
			record = EXCLUDED.record,
			nickname = EXCLUDED.nickname,
			height = EXCLUDED.height,
			weight = EXCLUDED.weight,
			reach = EXCLUDED.reach,
			stance = EXCLUDED.stance,
			dob = EXCLUDED.dob,
			updated_at = NOW()
		RETURNING id
	`,
		fighter.Name,
		fighter.URL,
		fighter.Record,
		fighter.Nickname,
		fighter.Height,
		fighter.Weight,
		fighter.Reach,
		fighter.Stance,
		fighter.DOB,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *Repository) UpsertFight(ctx context.Context, fight *ufcstats.Fight) (string, error) {
	var id string
	eventID := ""
	if fight.EventID != nil {
		eventID = *fight.EventID
	}

	err := r.db.QueryRow(ctx, `
		INSERT INTO fights (
			event_id,
			url,
			weight_class,
			method,
			round,
			time,
			winner,
			red_fighter_id,
			blue_fighter_id,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		ON CONFLICT (url) DO UPDATE SET
			event_id = EXCLUDED.event_id,
			weight_class = EXCLUDED.weight_class,
			method = EXCLUDED.method,
			round = EXCLUDED.round,
			time = EXCLUDED.time,
			winner = EXCLUDED.winner,
			red_fighter_id = EXCLUDED.red_fighter_id,
			blue_fighter_id = EXCLUDED.blue_fighter_id,
			updated_at = NOW()
		RETURNING id
	`,
		eventID,
		fight.URL,
		fight.WeightClass,
		fight.Method,
		fight.Round,
		fight.Time,
		fight.Winner,
		fight.RedFighterID,
		fight.BlueFighterID,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}
