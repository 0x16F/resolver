package definitions

import (
	"github.com/0x16f/vpn-resolver/internal/config"
	"github.com/0x16f/vpn-resolver/internal/controller/httpsrv"
	"github.com/0x16f/vpn-resolver/internal/controller/httpsrv/handlers/outline"
	"github.com/0x16f/vpn-resolver/internal/controller/httpsrv/handlers/servers"
	"github.com/0x16f/vpn-resolver/internal/controller/httpsrv/handlers/users"
	"github.com/0x16f/vpn-resolver/internal/usecase/configparser"
	"github.com/0x16f/vpn-resolver/internal/usecase/errors"
	"github.com/0x16f/vpn-resolver/internal/usecase/serversservice"
	"github.com/0x16f/vpn-resolver/internal/usecase/usersservice"
	"github.com/sarulabs/di"
)

const (
	HttpSrvDef = "http-srv"
)

func getHTTPSrvDef() di.Def {
	return di.Def{
		Name:  HttpSrvDef,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg, _ := ctn.Get(ConfigDef).(config.Config)

			serversSrv, _ := ctn.Get(serversServiceDef).(*serversservice.Service)
			usersSrv, _ := ctn.Get(usersServiceDef).(*usersservice.Service)
			errSrv, _ := ctn.Get(errorsService).(errors.Service)
			parserSrv, _ := ctn.Get(parserServiceDef).(*configparser.Service)

			srv := httpsrv.New()

			usersHandler := users.New(cfg.App.URI, usersSrv, errSrv, parserSrv)
			serversHandler := servers.New(serversSrv, errSrv)
			outlineHandler := outline.New(serversSrv, usersSrv, parserSrv, errSrv)

			v1 := srv.Group("/api/v1")
			{
				users := v1.Group("/users")
				{
					users.Get("", usersHandler.GetUsers)
					users.Get("/:id", usersHandler.GetUser)
					users.Post("", usersHandler.CreateUser)
					users.Delete("/:id", usersHandler.DeleteUser)
				}

				servers := v1.Group("/servers")
				{
					servers.Get("", serversHandler.GetServers)
					servers.Get("/:id", serversHandler.GetServer)
					servers.Post("", serversHandler.CreateServer)
					servers.Delete("/:id", serversHandler.DeleteServer)
				}
			}

			configs := srv.Group("")
			{
				configs.Get("/:id", outlineHandler.GetConfig)
			}

			return srv, nil
		},
	}
}
