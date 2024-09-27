package definitions

import (
	"github.com/0x16f/vpn-resolver/internal/infrastructure/repo/metrics"
	"github.com/0x16f/vpn-resolver/internal/infrastructure/repo/servers"
	"github.com/0x16f/vpn-resolver/internal/infrastructure/repo/users"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sarulabs/di"
)

const (
	serversRepoDef = "servers-repo"
	usersRepoDef   = "users-repo"
	metricsRepoDef = "metrics-repo"
)

func getServersRepoDef() di.Def {
	return di.Def{
		Name:  serversRepoDef,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			db, _ := ctn.Get(databaseDef).(*pgxpool.Pool)

			return servers.NewRepo(db), nil
		},
	}
}

func getUsersRepoDef() di.Def {
	return di.Def{
		Name:  usersRepoDef,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			db, _ := ctn.Get(databaseDef).(*pgxpool.Pool)

			return users.NewRepo(db), nil
		},
	}
}

func getMetricsRepoDef() di.Def {
	return di.Def{
		Name:  metricsRepoDef,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			db, _ := ctn.Get(databaseDef).(*pgxpool.Pool)

			return metrics.NewRepo(db), nil
		},
	}
}
