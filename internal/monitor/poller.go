package monitor

import (
	"errors"
	"github.com/chlp/ui/internal/device"
	"github.com/chlp/ui/pkg/application"
	"github.com/chlp/ui/pkg/logger"
	"os"
	"sync"
	"time"
)

const (
	devicesStatusPollInterval = 5 * time.Second
	maxParallelPolling        = 128
)

func (m *Monitor) loadPersistedDevicesStatus() error {
	m.devicesStatusMu.Lock()
	defer m.devicesStatusMu.Unlock()

	err := m.devicesStatusStore.LoadJSON(&m.devicesStatus)
	if errors.Is(err, os.ErrNotExist) {
		return m.devicesStatusStore.SaveJSON(&m.devicesStatus)
	}
	return err
}

func (m *Monitor) pollAllDevicesStatus(app *application.App) {
	app.Wg.Add(1)
	defer app.Wg.Done()

	m.pollAllDevicesStatusOnce()

	ticker := time.NewTicker(devicesStatusPollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-app.Ctx.Done():
			return
		case <-ticker.C:
			m.pollAllDevicesStatusOnce()
		}
	}
}

func (m *Monitor) pollAllDevicesStatusOnce() {
	devicesList := m.GetDevicesList()

	wg := sync.WaitGroup{}
	goRoutinesLimiterChan := make(chan struct{}, maxParallelPolling)

	for _, address := range devicesList {
		wg.Add(1)
		goRoutinesLimiterChan <- struct{}{}
		go func(address string) {
			defer wg.Done()
			defer func() { <-goRoutinesLimiterChan }()

			if err := m.pollDeviceStatus(address); err != nil {
				logger.Debugf("Monitor::pollAllDevicesStatus: pollDeviceStatus err (%s): %v", address, err)
			}
		}(address)
	}

	wg.Wait()
}

func (m *Monitor) pollDeviceStatus(address string) error {
	info, err := getRestInfo(address)
	if err != nil {
		logger.Debugf("Monitor::pollDeviceStatus: REST failed, trying gRPC (%s): %v", address, err)
		info, err = getGrpcInfo(address)
		if err != nil {
			logger.Debugf("Monitor::pollDeviceStatus: gRPC failed (%s): %v", address, err)
			return err
		}
	}

	if info == nil {
		logger.Debugf("Monitor::pollDeviceStatus: empty info (%s)", address)
		return nil
	}

	m.devicesStatusMu.Lock()
	defer m.devicesStatusMu.Unlock()

	m.devicesStatus[address] = device.Status{
		Info:      *info,
		UpdatedAt: time.Now(),
	}

	err = m.devicesStatusStore.SaveJSON(&m.devicesStatus)
	if err != nil {
		logger.Printf("Monitor::pollDeviceStatus: devicesStatusStore.SaveJSON (%s): %v", address, err)
		return err
	}

	return nil
}
