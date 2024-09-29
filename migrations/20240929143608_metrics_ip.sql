-- +goose Up

ALTER TABLE vp_usage_metrics ADD COLUMN ip VARCHAR(255) NOT NULL DEFAULT '';

-- +goose Down

ALTER TABLE vp_usage_metrics DROP COLUMN ip;