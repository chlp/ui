package monitor

import (
	"errors"
	"github.com/chlp/ui/pkg/logger"
	"os"
	"time"
)

const devicesListFileWatchInterval = 5 * time.Second

func (m *Monitor) syncDevicesListWithStore() error {
	m.devicesListMu.Lock()
	defer m.devicesListMu.Unlock()

	err := m.devicesListStore.LoadJSON(&m.devicesList)
	if err == nil {
		return nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return m.devicesListStore.SaveJSON(&m.devicesList)
	}

	return err
}

func (m *Monitor) watchDevicesListStoreChanges() {
	ticker := time.NewTicker(devicesListFileWatchInterval)
	defer ticker.Stop()
	for range ticker.C {
		if err := m.syncDevicesListWithStore(); err != nil {
			logger.Printf("Monitor::watchDevicesfailed: syncDevicesListWithStore err: %v", err)
		}
	}
}
