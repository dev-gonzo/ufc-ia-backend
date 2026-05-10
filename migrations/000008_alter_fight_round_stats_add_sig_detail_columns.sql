-- +goose Up

ALTER TABLE fight_round_stats
    ADD COLUMN head_landed INTEGER,
    ADD COLUMN head_attempted INTEGER,
    ADD COLUMN body_landed INTEGER,
    ADD COLUMN body_attempted INTEGER,
    ADD COLUMN leg_landed INTEGER,
    ADD COLUMN leg_attempted INTEGER,
    ADD COLUMN distance_landed INTEGER,
    ADD COLUMN distance_attempted INTEGER,
    ADD COLUMN clinch_landed INTEGER,
    ADD COLUMN clinch_attempted INTEGER,
    ADD COLUMN ground_landed INTEGER,
    ADD COLUMN ground_attempted INTEGER;

-- +goose Down

ALTER TABLE fight_round_stats
    DROP COLUMN ground_attempted,
    DROP COLUMN ground_landed,
    DROP COLUMN clinch_attempted,
    DROP COLUMN clinch_landed,
    DROP COLUMN distance_attempted,
    DROP COLUMN distance_landed,
    DROP COLUMN leg_attempted,
    DROP COLUMN leg_landed,
    DROP COLUMN body_attempted,
    DROP COLUMN body_landed,
    DROP COLUMN head_attempted,
    DROP COLUMN head_landed;
