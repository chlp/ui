package main

import (
	"github.com/chlp/ui/internal/api/grpc"
	"github.com/chlp/ui/internal/api/rest"
	"github.com/chlp/ui/internal/config"
	"github.com/chlp/ui/internal/monitor"
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

	devicesMonitor := monitor.MustNewMonitor(
		filestore.NewFileStore(cfg.DevicesListFile),
		filestore.NewFileStore(cfg.DevicesStatusFile),
	)
	go rest.StartRestServer(cfg.RestPort, device, devicesMonitor)
	go grpc.StartGrpcServer(cfg.GrpcPort, device)
}
