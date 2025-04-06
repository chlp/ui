package monitor

import (
	"github.com/chlp/ui/internal/model"
	"github.com/chlp/ui/pkg/application"
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

func MustNewMonitor(app *application.App, devicesListStore, devicesStatusStore Store) *Monitor {
	if m, err := NewMonitor(app, devicesListStore, devicesStatusStore); err != nil {
		logger.Fatalf("MustNewMonitor: failed to load devicesList: %v", err)
		return nil
	} else {
		return m
	}
}

func NewMonitor(app *application.App, devicesListStore, devicesStatusStore Store) (*Monitor, error) {
	if devicesListStore == nil {
		logger.Printf("Monitor: starting without monitor (no devicesListStore)")
		return nil, nil
	}
	if devicesStatusStore == nil {
		logger.Printf("Monitor: starting without monitor (no devicesStatusStore)")
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
	go m.watchDevicesListStoreChanges(app)

	if err := m.loadPersistedDevicesStatus(); err != nil {
		return nil, err
	}
	go m.pollAllDevicesStatus()

	logger.Printf("Monitor: monitor started, devices in list: %d", len(m.devicesList))
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
