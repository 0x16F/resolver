package config

import (
	"github.com/0x16f/vpn-resolver/internal/infrastructure/repo/database"
	"github.com/ilyakaznacheev/cleanenv"
)

type App struct {
	Port           uint16 `env:"HTTP_PORT" env-default:"8080"`
	Name           string `env:"APP_NAME" env-default:"resolver"`
	URI            string `env:"APP_URI" required:"true"`
	MigrationsPath string `env:"MIGRATIONS_PATH" required:"true"`
	ErrorsPath     string `env:"ERRORS_PATH" required:"true"`
	Secret         string `env:"SECRET" required:"true"`
}

type Config struct {
	App      App
	Database database.Config
}

func New() (Config, error) {
	config := Config{}

	err := cleanenv.ReadEnv(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
