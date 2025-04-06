package monitor

import (
	"errors"
	"github.com/chlp/ui/pkg/logger"
	"os"
	"time"
)

const devicesListFileWatchInterval = 5 * time.Second

func (m *Monitor) GetDevicesList() []string {
	m.devicesListMu.RLock()
	m.devicesListMu.RUnlock()

	devicesList := make([]string, len(m.devicesList))
	copy(devicesList, m.devicesList)
	return devicesList
}

func (m *Monitor) AddDevice(address string) (bool, error) {
	m.devicesListMu.Lock()
	defer m.devicesListMu.Unlock()

	for _, device := range m.devicesList {
		if device == address {
			return false, nil
		}
	}

	m.devicesList = append(m.devicesList, address)

	if err := m.devicesListStore.SaveJSON(&m.devicesList); err != nil {
		logger.Printf("Monitor::AddDevice: failed to save devices list: %v", err)
		return false, err
	}

	return true, nil
}

func (m *Monitor) RemoveDevice(address string) (bool, error) {
	m.devicesListMu.Lock()
	listChanged := false
	for i, device := range m.devicesList {
		if device == address {
			m.devicesList = append(m.devicesList[:i], m.devicesList[i+1:]...)
			listChanged = true
			break
		}
	}
	m.devicesListMu.Unlock()

	if !listChanged {
		return false, nil
	}

	if err := m.devicesListStore.SaveJSON(&m.devicesList); err != nil {
		logger.Printf("Monitor::RemoveDevice: failed to save devices list: %v", err)
		return false, err
	}

	return true, nil
}

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
