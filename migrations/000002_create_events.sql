-- +goose Up

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    date TEXT NOT NULL,
    location TEXT NOT NULL,
    event_sync BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_events_url ON events(url);
CREATE INDEX idx_events_event_sync ON events(event_sync);

-- +goose Down

DROP TABLE events;
