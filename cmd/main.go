package main

import (
	"github.com/chlp/ui/internal/monitor"
	"github.com/chlp/ui/pkg/file_store"
	"time"

	"github.com/chlp/ui/internal/api/grpc"
	"github.com/chlp/ui/internal/api/rest"
	"github.com/chlp/ui/internal/config"
	"github.com/chlp/ui/pkg/logger"
)

const (
	configFile = "config.json"
)

func main() {
	cfg := config.MustLoadOrCreateConfig(configFile)

	logger.InitLogger(cfg.LogFile)

	devicesMonitor := monitor.MustNewMonitor(
		file_store.NewFileStore(cfg.DevicesListFile),
		file_store.NewFileStore(cfg.DevicesStatusFile),
	)
	go rest.StartRestServer(cfg)
	go grpc.StartGrpcServer(cfg)
}
