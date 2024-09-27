package definitions

import (
	"github.com/0x16f/vpn-resolver/internal/config"
	"github.com/0x16f/vpn-resolver/internal/infrastructure/repo/servers"
	"github.com/0x16f/vpn-resolver/internal/infrastructure/repo/users"
	"github.com/0x16f/vpn-resolver/internal/service/outline"
	"github.com/0x16f/vpn-resolver/internal/usecase/configparser"
	"github.com/0x16f/vpn-resolver/internal/usecase/errors"
	"github.com/0x16f/vpn-resolver/internal/usecase/serversservice"
	"github.com/0x16f/vpn-resolver/internal/usecase/usersservice"
	"github.com/sarulabs/di"
)

const (
	errorsService = "errors-service"

	serversServiceDef = "servers-service"
	usersServiceDef   = "users-service"

	outlineServiceDef = "outline-service"
	parserServiceDef  = "parser-service"
)

func getErrorsServiceDef() di.Def {
	return di.Def{
		Name:  errorsService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg, _ := ctn.Get(ConfigDef).(config.Config)

			return errors.New(cfg.App.ErrorsPath), nil
		},
	}
}

func getServersServiceDef() di.Def {
	return di.Def{
		Name:  serversServiceDef,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			repo, _ := ctn.Get(serversRepoDef).(*servers.Repo)
			errSrv, _ := ctn.Get(errorsService).(errors.Service)

			return serversservice.New(repo, errSrv), nil
		},
	}
}

func getUsersServiceDef() di.Def {
	return di.Def{
		Name:  usersServiceDef,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			repo, _ := ctn.Get(usersRepoDef).(*users.Repo)
			errSrv, _ := ctn.Get(errorsService).(errors.Service)
			serversSrv, _ := ctn.Get(serversServiceDef).(*serversservice.Service)
			outlineSrv, _ := ctn.Get(outlineServiceDef).(*outline.Service)

			return usersservice.New(repo, errSrv, serversSrv, outlineSrv), nil
		},
	}
}

func getOutlineServiceDef() di.Def {
	return di.Def{
		Name:  outlineServiceDef,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return outline.New(), nil
		},
	}
}

func getParserServiceDef() di.Def {
	return di.Def{
		Name:  parserServiceDef,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg, _ := ctn.Get(ConfigDef).(config.Config)

			return configparser.New(cfg.App.Secret)
		},
	}
}
