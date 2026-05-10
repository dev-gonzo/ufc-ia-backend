-- +goose Up

CREATE TABLE fighters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    record TEXT NOT NULL DEFAULT '',
    nickname TEXT,
    height TEXT,
    weight TEXT,
    reach TEXT,
    stance TEXT,
    dob TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fighters_url ON fighters(url);

CREATE TABLE fights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    url TEXT NOT NULL UNIQUE,
    weight_class TEXT NOT NULL DEFAULT '',
    method TEXT NOT NULL DEFAULT '',
    round INTEGER NOT NULL DEFAULT 0,
    time TEXT NOT NULL DEFAULT '',
    winner TEXT NOT NULL DEFAULT '',
    red_fighter_id UUID REFERENCES fighters(id) ON DELETE SET NULL,
    blue_fighter_id UUID REFERENCES fighters(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fights_event_id ON fights(event_id);
CREATE INDEX idx_fights_url ON fights(url);

-- +goose Down

DROP TABLE fights;
DROP TABLE fighters;
