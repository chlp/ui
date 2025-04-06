package monitor

import (
	"fmt"
	"github.com/chlp/ui/internal/device"
	"github.com/chlp/ui/pkg/application"
	"github.com/chlp/ui/pkg/logger"
	"reflect"
	"sync"
	"time"
)

const durationToSetDeviceUnavailable = 15 * time.Second

type Monitor struct {
	devicesList      []string
	devicesListMu    sync.RWMutex
	devicesListStore Store

	devicesStatusStore Store

	devicesStatus   map[string]device.Status
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
	if app == nil {
		err := fmt.Errorf("could not start without app")
		logger.Printf("Monitor: err %v", err)
		return nil, err
	}
	if devicesListStore == nil || reflect.ValueOf(devicesListStore).IsNil() {
		logger.Printf("Monitor: starting without monitor (no devicesListStore)")
		return nil, nil
	}
	if devicesStatusStore == nil || reflect.ValueOf(devicesStatusStore).IsNil() {
		logger.Printf("Monitor: starting without monitor (no devicesStatusStore)")
		return nil, nil
	}

	m := &Monitor{
		devicesList:        make([]string, 0),
		devicesListStore:   devicesListStore,
		devicesListMu:      sync.RWMutex{},
		devicesStatusStore: devicesStatusStore,
		devicesStatus:      make(map[string]device.Status),
		devicesStatusMu:    sync.RWMutex{},
	}

	if err := m.syncDevicesListWithStore(); err != nil {
		return nil, err
	}
	go m.watchDevicesListStoreChanges(app)

	if err := m.loadPersistedDevicesStatus(); err != nil {
		return nil, err
	}
	go m.pollAllDevicesStatus(app)

	logger.Printf("Monitor: monitor started, devices in list: %d", len(m.devicesList))
	return m, nil
}

func (m *Monitor) GetDevicesStatus() map[string]device.Status {
	m.devicesStatusMu.RLock()
	defer m.devicesStatusMu.RUnlock()
	m.devicesListMu.RLock()
	defer m.devicesListMu.RUnlock()

	devicesStatus := make(map[string]device.Status, len(m.devicesList))
	for _, addr := range m.devicesList {
		if deviceStatus, ok := m.devicesStatus[addr]; ok {
			if deviceStatus.UpdatedAt.Before(time.Now().Add(-durationToSetDeviceUnavailable)) {
				deviceStatus.Status = device.StatusUnavailable
			}
			devicesStatus[addr] = deviceStatus
		} else {
			devicesStatus[addr] = device.Status{
				Info: device.Info{
					Status: device.StatusUnknown,
				},
				UpdatedAt: time.Time{},
			}
		}
	}
	return devicesStatus
}
