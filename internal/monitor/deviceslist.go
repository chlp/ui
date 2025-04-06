package monitor

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/chlp/ui/pkg/application"
	"github.com/chlp/ui/pkg/logger"
	"os"
	"strings"
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

	logger.Printf("Monitor::AddDevice: devices list changed: %d", len(m.devicesList))

	if err := m.devicesListStore.SaveJSON(&m.devicesList); err != nil {
		logger.Printf("Monitor::AddDevice: failed to save devices list: %v", err)
		return false, err
	}

	return true, nil
}

func (m *Monitor) RemoveDevice(address string) (bool, error) {
	m.devicesListMu.Lock()
	defer m.devicesListMu.Unlock()

	listChanged := false
	for i, device := range m.devicesList {
		if device == address {
			m.devicesList = append(m.devicesList[:i], m.devicesList[i+1:]...)
			listChanged = true
			break
		}
	}

	if !listChanged {
		return false, nil
	}

	logger.Printf("Monitor::RemoveDevice: devices list changed: %d", len(m.devicesList))

	if err := m.devicesListStore.SaveJSON(&m.devicesList); err != nil {
		logger.Printf("Monitor::RemoveDevice: failed to save devices list: %v", err)
		return false, err
	}

	return true, nil
}

func (m *Monitor) syncDevicesListWithStore() error {
	m.devicesListMu.Lock()
	defer m.devicesListMu.Unlock()

	hashBeforeChanges := m.devicesListHash()

	err := m.devicesListStore.LoadJSON(&m.devicesList)
	if err == nil {
		if hashBeforeChanges != m.devicesListHash() {
			logger.Printf("Monitor::syncDevicesListWithStore: devices list changed: %d", len(m.devicesList))
		}
		return nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return m.devicesListStore.SaveJSON(&m.devicesList)
	}

	return err
}

func (m *Monitor) watchDevicesListStoreChanges(app *application.App) {
	app.Wg.Add(1)
	defer app.Wg.Done()

	ticker := time.NewTicker(devicesListFileWatchInterval)
	defer ticker.Stop()
	for {
		select {
		case <-app.Ctx.Done():
			return
		case <-ticker.C:
			if err := m.syncDevicesListWithStore(); err != nil {
				logger.Printf("Monitor::watchDevicesListStoreChanges: syncDevicesListWithStore err: %v", err)
			}
		}
	}
}

func (m *Monitor) devicesListHash() string {
	joined := strings.Join(m.devicesList, ",")
	hash := sha256.Sum256([]byte(joined))
	return hex.EncodeToString(hash[:])
}
