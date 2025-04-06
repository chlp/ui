package main

import (
	"github.com/chlp/ui/internal/api/grpc"
	"github.com/chlp/ui/internal/api/rest"
	"github.com/chlp/ui/internal/config"
	"github.com/chlp/ui/internal/device"
	"github.com/chlp/ui/internal/monitor"
	"github.com/chlp/ui/pkg/application"
	"github.com/chlp/ui/pkg/filestore"
	"github.com/chlp/ui/pkg/logger"
	"os"
)

const (
	defaultAppName    = "ui-monitor"
	defaultConfigFile = "config.json"
)

func main() {
	configPath := defaultConfigFile
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	cfg := config.MustLoadOrCreateConfig(configPath)

	appName := defaultAppName
	if cfg.Device.ID != "" {
		appName = cfg.Device.ID
	}

	app, appDone := application.NewApp(appName, cfg.LogFile, cfg.Debug)

	localDevice := device.GetLocalDevice(cfg.Device, cfg.ChecksumCmd, *cfg.ChecksumEmulate)

	devicesMonitor := monitor.MustNewMonitor(
		app,
		filestore.NewFileStore(cfg.DevicesListFile),
		filestore.NewFileStore(cfg.DevicesStatusFile),
	)
	go rest.StartServer(app, cfg.RestPort, localDevice, devicesMonitor)
	go grpc.StartServer(app, cfg.GrpcPort, localDevice)

	<-appDone

	logger.Printf("Application gracefully shut down")
}
