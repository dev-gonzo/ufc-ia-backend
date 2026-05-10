package scraping

import (
	"context"
	"strings"
	"time"
	"ufc-backend/internal/scraping/tapology"
	"ufc-backend/internal/scraping/ufc"
	"ufc-backend/internal/scraping/ufcstats"

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

func (r *Repository) UpsertFighterUFCDetails(ctx context.Context, fighterID string, d *ufc.AthleteDetails) error {
	var photo any
	if d.PhotoWebPBase64 != nil && strings.TrimSpace(*d.PhotoWebPBase64) != "" {
		photo = *d.PhotoWebPBase64
	}

	var active any
	if d.IsActive != nil {
		active = *d.IsActive
	}

	var hometown any
	if d.Hometown != nil && strings.TrimSpace(*d.Hometown) != "" {
		hometown = *d.Hometown
	}

	var style any
	if d.FightingStyle != nil && strings.TrimSpace(*d.FightingStyle) != "" {
		style = *d.FightingStyle
	}

	var height any
	if d.Height != nil && strings.TrimSpace(*d.Height) != "" {
		height = *d.Height
	}

	var weight any
	if d.Weight != nil && strings.TrimSpace(*d.Weight) != "" {
		weight = *d.Weight
	}

	var debut any
	if d.UFCDebut != nil && strings.TrimSpace(*d.UFCDebut) != "" {
		debut = *d.UFCDebut
	}

	var reach any
	if d.Reach != nil && strings.TrimSpace(*d.Reach) != "" {
		reach = *d.Reach
	}

	var legReach any
	if d.LegReach != nil && strings.TrimSpace(*d.LegReach) != "" {
		legReach = *d.LegReach
	}

	var facts any
	if d.FighterFacts != nil && strings.TrimSpace(*d.FighterFacts) != "" {
		facts = *d.FighterFacts
	}

	var history any
	if d.UFCHistory != nil && strings.TrimSpace(*d.UFCHistory) != "" {
		history = *d.UFCHistory
	}

	var qa any
	if d.QA != nil && strings.TrimSpace(*d.QA) != "" {
		qa = *d.QA
	}

	_, err := r.db.Exec(ctx, `
		INSERT INTO fighter_ufc_details (
			fighter_id,
			athlete_url,
			athlete_slug,
			photo_webp_base64,
			is_active,
			hometown,
			fighting_style,
			height,
			weight,
			ufc_debut,
			reach,
			leg_reach,
			fighter_facts,
			ufc_history,
			qa,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,NOW())
		ON CONFLICT (fighter_id) DO UPDATE SET
			athlete_url = EXCLUDED.athlete_url,
			athlete_slug = EXCLUDED.athlete_slug,
			photo_webp_base64 = EXCLUDED.photo_webp_base64,
			is_active = EXCLUDED.is_active,
			hometown = EXCLUDED.hometown,
			fighting_style = EXCLUDED.fighting_style,
			height = EXCLUDED.height,
			weight = EXCLUDED.weight,
			ufc_debut = EXCLUDED.ufc_debut,
			reach = EXCLUDED.reach,
			leg_reach = EXCLUDED.leg_reach,
			fighter_facts = EXCLUDED.fighter_facts,
			ufc_history = EXCLUDED.ufc_history,
			qa = EXCLUDED.qa,
			updated_at = NOW()
	`,
		fighterID,
		d.AthleteURL,
		d.AthleteSlug,
		photo,
		active,
		hometown,
		style,
		height,
		weight,
		debut,
		reach,
		legReach,
		facts,
		history,
		qa,
	)
	return err
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
