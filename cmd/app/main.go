package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/0x16f/vpn-resolver/internal/config"
	"github.com/0x16f/vpn-resolver/internal/controller/httpsrv"
	"github.com/0x16f/vpn-resolver/internal/definitions"
	"github.com/sirupsen/logrus"
)

func main() {
	stopCtx, cancel := context.WithCancel(context.Background())

	di, err := definitions.New(stopCtx)
	if err != nil {
		logrus.Fatalf("failed to create DI container: %v", err)
	}

	cfg, _ := di.Get(definitions.ConfigDef).(config.Config)

	server, _ := di.Get(definitions.HttpSrvDef).(*httpsrv.Server)

	go func() {
		if err := server.Start(cfg.App.Port); err != nil {
			logrus.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Graceful shutdown
	_ = <-stop

	cancel()
}
