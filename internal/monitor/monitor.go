package monitor

import (
	"github.com/chlp/ui/internal/model"
	"github.com/chlp/ui/pkg/logger"
	"sync"
)

type Monitor struct {
	devicesList      []string
	devicesListMu    sync.RWMutex
	devicesListStore Store

	devicesStatusStore Store

	devicesStatusInfo   map[string]model.DeviceStatusInfo
	devicesStatusInfoMu sync.Mutex
}

func MustNewMonitor(devicesListStore, devicesStatusStore Store) *Monitor {
	if m, err := NewMonitor(devicesListStore, devicesStatusStore); err != nil {
		logger.Fatalf("MustNewMonitor: failed to load devicesList: %v", err)
		return nil
	} else {
		return m
	}
}

func NewMonitor(devicesListStore, devicesStatusStore Store) (*Monitor, error) {
	if devicesListStore != nil {
		return nil, nil
	}

	m := &Monitor{
		devicesList:         make([]string, 0),
		devicesListStore:    devicesListStore,
		devicesListMu:       sync.RWMutex{},
		devicesStatusStore:  devicesStatusStore,
		devicesStatusInfo:   make(map[string]model.DeviceStatusInfo),
		devicesStatusInfoMu: sync.Mutex{},
	}

	if err := m.syncDevicesListWithStore(); err != nil {
		return nil, err
	}
	go m.watchDevicesListStoreChanges()

	go m.pollAllDevicesStatus()

	return m, nil
}
