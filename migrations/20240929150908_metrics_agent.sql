-- +goose Up

ALTER TABLE vp_usage_metrics ADD COLUMN agent VARCHAR(512) NOT NULL DEFAULT '';

-- +goose Down

ALTER TABLE vp_usage_metrics DROP COLUMN agent;