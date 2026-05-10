-- +goose Up

CREATE TABLE fighter_ufc_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fighter_id UUID NOT NULL UNIQUE REFERENCES fighters(id) ON DELETE CASCADE,
    athlete_url TEXT NOT NULL,
    athlete_slug TEXT NOT NULL,
    photo_webp_base64 TEXT,
    is_active BOOLEAN,
    hometown TEXT,
    fighting_style TEXT,
    height TEXT,
    weight TEXT,
    ufc_debut TEXT,
    reach TEXT,
    leg_reach TEXT,
    fighter_facts TEXT,
    ufc_history TEXT,
    qa TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fighter_ufc_details_fighter_id ON fighter_ufc_details(fighter_id);

-- +goose Down

DROP TABLE fighter_ufc_details;
