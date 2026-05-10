-- +goose Up

CREATE TABLE referees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE judges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE fight_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fight_id UUID NOT NULL UNIQUE REFERENCES fights(id) ON DELETE CASCADE,
    is_title_bout BOOLEAN,
    rounds INTEGER,
    referee_id UUID REFERENCES referees(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fight_details_fight_id ON fight_details(fight_id);

CREATE TABLE fight_judges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fight_id UUID NOT NULL REFERENCES fights(id) ON DELETE CASCADE,
    judge_id UUID NOT NULL REFERENCES judges(id) ON DELETE CASCADE,
    red_score INTEGER,
    blue_score INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (fight_id, judge_id)
);

CREATE INDEX idx_fight_judges_fight_id ON fight_judges(fight_id);

CREATE TABLE fight_round_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fight_id UUID NOT NULL REFERENCES fights(id) ON DELETE CASCADE,
    round INTEGER NOT NULL,
    corner TEXT NOT NULL,
    kd INTEGER,
    sig_landed INTEGER,
    sig_attempted INTEGER,
    sig_pct INTEGER,
    total_landed INTEGER,
    total_attempted INTEGER,
    td_landed INTEGER,
    td_attempted INTEGER,
    td_pct INTEGER,
    sub_att INTEGER,
    rev INTEGER,
    ctrl TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (fight_id, round, corner)
);

CREATE INDEX idx_fight_round_stats_fight_id ON fight_round_stats(fight_id);

-- +goose Down

DROP TABLE fight_round_stats;
DROP TABLE fight_judges;
DROP TABLE fight_details;
DROP TABLE judges;
DROP TABLE referees;
