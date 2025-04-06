package main

import (
	"github.com/chlp/ui/internal/api/grpc"
	"github.com/chlp/ui/internal/api/rest"
	"github.com/chlp/ui/internal/config"
	"github.com/chlp/ui/internal/monitor"
	"github.com/chlp/ui/pkg/application"
	"github.com/chlp/ui/pkg/filestore"
	"github.com/chlp/ui/pkg/logger"
)

const (
	configFile = "config.json"
)

func main() {
	cfg := config.MustLoadOrCreateConfig(configFile)
	device := cfg.Device
	// todo: update device checksum

	logger.InitLogger(cfg.LogFile)

	app, appDone := application.NewApp()

	devicesMonitor := monitor.MustNewMonitor(
		app,
		filestore.NewFileStore(cfg.DevicesListFile),
		filestore.NewFileStore(cfg.DevicesStatusFile),
	)
	go rest.StartRestServer(app, cfg.RestPort, device, devicesMonitor)
	go grpc.StartGrpcServer(app, cfg.GrpcPort, device)

	<-appDone

	logger.Printf("Application gracefully shut down")
}
