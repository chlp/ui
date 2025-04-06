package config

import (
	"errors"
	"fmt"
	"github.com/chlp/ui/internal/device"
	"github.com/chlp/ui/pkg/filestore"
	"github.com/chlp/ui/pkg/logger"
	"os"
	"time"
)

type Config struct {
	LogFile           string       `json:"log_file"`          // file where logs will be written
	Debug             bool         `json:"debug"`             // write debug logs
	DevicesListFile   string       `json:"devices_list_file"` // file with list of devices
	DevicesStatusFile string       `json:"monitor_file"`      // file for persisting all polled data from devices
	GrpcPort          string       `json:"grpc_port"`         // port to start gRPC server
	RestPort          string       `json:"rest_port"`         // port to start REST server
	Device            *device.Info `json:"device"`            // configuration if using the app as a monitored device
	ChecksumCmd       string       `json:"checksum_cmd"`      // command to calculate checksum
	ChecksumEmulate   *bool        `json:"checksum_emulate"`  // emulate checksum if checksum cmd return empty string
}

const (
	defaultDevicesListFile = "devices_list.json"
	defaultMonitorFile     = "monitor.json"
	defaultLogFile         = "app.log"
	defaultGrpcPort        = ":50051"
	defaultRestPort        = ":8080"
	defaultChecksumEmulate = true
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
	defaultChecksumEmulateVar := defaultChecksumEmulate
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		deviceConfig := &device.Info{
			ID:              generateID(),
			HardwareVersion: "0.0.1",
			SoftwareVersion: "1.0.0",
			FirmwareVersion: "0.0.5",
			Status:          device.StatusOk,
			Checksum:        "",
		}
		cfg = &Config{
			DevicesListFile:   defaultDevicesListFile,
			DevicesStatusFile: defaultMonitorFile,
			LogFile:           defaultLogFile,
			GrpcPort:          defaultGrpcPort,
			RestPort:          defaultRestPort,
			Device:            deviceConfig,
			ChecksumCmd:       "",
			ChecksumEmulate:   &defaultChecksumEmulateVar,
		}
		return cfg, filestore.SaveJSON(configFile, cfg)
	} else {
		if err = filestore.LoadJSON(configFile, &cfg); err != nil {
			return nil, err
		}
		if cfg.ChecksumEmulate == nil {
			cfg.ChecksumEmulate = &defaultChecksumEmulateVar
		}
		return cfg, nil
	}
}

func generateID() string {
	return fmt.Sprintf("dev-%d", time.Now().UnixNano())
}
