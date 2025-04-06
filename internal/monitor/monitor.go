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

	devicesStatus   map[string]model.DeviceStatus
	devicesStatusMu sync.RWMutex
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
		devicesList:        make([]string, 0),
		devicesListStore:   devicesListStore,
		devicesListMu:      sync.RWMutex{},
		devicesStatusStore: devicesStatusStore,
		devicesStatus:      make(map[string]model.DeviceStatus),
		devicesStatusMu:    sync.RWMutex{},
	}

	if err := m.syncDevicesListWithStore(); err != nil {
		return nil, err
	}
	go m.watchDevicesListStoreChanges()

	go m.pollAllDevicesStatus()

	return m, nil
}

func (m *Monitor) GetDevicesStatus() map[string]model.DeviceStatus {
	m.devicesStatusMu.RLock()
	defer m.devicesStatusMu.RUnlock()

	devicesStatus := make(map[string]model.DeviceStatus, len(m.devicesStatus))
	for k, v := range m.devicesStatus {
		devicesStatus[k] = v
	}
	return devicesStatus
}
