-- +goose Up

CREATE TABLE fight_bonuses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fight_id UUID NOT NULL REFERENCES fights(id) ON DELETE CASCADE,
    bonus_type TEXT NOT NULL,
    recipient_corner TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (fight_id, bonus_type, recipient_corner)
);

CREATE INDEX idx_fight_bonuses_fight_id ON fight_bonuses(fight_id);

CREATE TABLE fight_sig_strike_breakdown (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fight_id UUID NOT NULL REFERENCES fights(id) ON DELETE CASCADE,
    scope TEXT NOT NULL,
    label TEXT NOT NULL,
    red_pct INTEGER,
    blue_pct INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (fight_id, scope, label)
);

CREATE INDEX idx_fight_sig_strike_breakdown_fight_id ON fight_sig_strike_breakdown(fight_id);

-- +goose Down

DROP TABLE fight_sig_strike_breakdown;
DROP TABLE fight_bonuses;
