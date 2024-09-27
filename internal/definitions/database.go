package definitions

import (
	"context"

	"github.com/0x16f/vpn-resolver/internal/config"
	"github.com/0x16f/vpn-resolver/internal/infrastructure/repo/database"
	"github.com/0x16f/vpn-resolver/internal/infrastructure/repo/migrations"
	"github.com/sarulabs/di"
)

const (
	databaseDef = "database"
)

func getDatabaseDef(ctx context.Context) di.Def {
	return di.Def{
		Name:  databaseDef,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg, _ := ctn.Get(ConfigDef).(config.Config)

			err := migrations.Run(cfg.Database.DSN(), cfg.App.MigrationsPath)
			if err != nil {
				return nil, err
			}

			return database.NewConnection(ctx, cfg.Database)
		},
	}
}
