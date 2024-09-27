package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     uint16 `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-default:"postgres"`
	Database string `env:"DB_DATABASE" env-default:"postgres"`
}

func (c Config) DSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", c.User, c.Password, c.Host, c.Port, c.Database)
}

func NewConnection(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, cfg.DSN())
}
