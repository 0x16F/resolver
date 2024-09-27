-- +goose Up

CREATE TABLE vp_usage_metrics (
    id SERIAL PRIMARY KEY,
    server_id INT NOT NULL REFERENCES vp_servers(id),
    user_id UUID NOT NULL REFERENCES vp_users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose Down

DROP TABLE vp_usage_metrics;