package config

import (
	"errors"
	"fmt"
	"github.com/chlp/ui/internal/model"
	"github.com/chlp/ui/pkg/file_store"
	"github.com/chlp/ui/pkg/logger"
	"os"
	"time"
)

type Config struct {
	DevicesListFile   string            `json:"devices_list_file"` // file with list of devices
	DevicesStatusFile string            `json:"monitor_file"`      // file for persisting all polled data from devices
	LogFile           string            `json:"log_file"`          // file where logs will be written
	GrpcPort          string            `json:"grpc_port"`         // port to start gRPC server
	RestPort          string            `json:"rest_port"`         // port to start REST server
	Device            *model.DeviceInfo `json:"device"`            // configuration if using the app as a monitored device
}

const (
	defaultDevicesListFile = "devices_list.json"
	defaultMonitorFile     = "monitor.json"
	defaultLogFile         = "app.log"
	defaultGrpcPort        = ":50051"
	defaultRestPort        = ":8080"
)

func MustLoadOrCreateConfig(configFile string) *Config {
	deviceConfig, err := LoadOrCreateConfig(configFile)
	if err != nil {
		logger.Fatalf("MustLoadOrCreateConfig: failed to load/create config: %v", err)
		return nil
	}
	return deviceConfig
}

func LoadOrCreateConfig(configFile string) (*Config, error) {
	var cfg *Config
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		deviceConfig := &model.DeviceInfo{
			ID:              generateID(),
			HardwareVersion: "0.0.1",
			SoftwareVersion: "1.0.0",
			FirmwareVersion: "0.0.5",
			Status:          model.DeviceStatusOk,
			Checksum:        "",
		}
		cfg = &Config{
			DevicesListFile:   defaultDevicesListFile,
			DevicesStatusFile: defaultMonitorFile,
			LogFile:           defaultLogFile,
			GrpcPort:          defaultGrpcPort,
			RestPort:          defaultRestPort,
			Device:            deviceConfig,
		}
		return cfg, file_store.SaveJSON(configFile, cfg)
	} else {
		if err = file_store.LoadJSON(configFile, &cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	}
}

func generateID() string {
	return fmt.Sprintf("dev-%d", time.Now().UnixNano())
}
