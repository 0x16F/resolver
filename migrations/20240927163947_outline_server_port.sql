-- +goose Up

ALTER TABLE vp_servers ADD COLUMN user_port INT NOT NULL;

-- +goose Down

ALTER TABLE vp_servers DROP COLUMN user_port;