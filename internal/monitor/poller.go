package monitor

import (
	"github.com/chlp/ui/internal/model"
	"github.com/chlp/ui/pkg/logger"
	"sync"
	"time"
)

const (
	devicesStatusPollInterval = 5 * time.Second
	maxParallelPolling        = 128
)

func (m *Monitor) pollAllDevicesStatus() {
	ticker := time.NewTicker(devicesStatusPollInterval)
	defer ticker.Stop()
	for range ticker.C {
		devicesList := m.GetDevicesList()

		wg := sync.WaitGroup{}
		goRoutinesLimiterChan := make(chan struct{}, maxParallelPolling)

		for _, address := range devicesList {
			wg.Add(1)
			goRoutinesLimiterChan <- struct{}{}
			go func() {
				defer wg.Done()
				defer func() { <-goRoutinesLimiterChan }()

				if err := m.pollDeviceStatus(address); err != nil {
					logger.Printf("Monitor::pollAllDevicesStatus: pollDeviceStatus err (%s): %v", address, err)
				}
			}()
		}

		wg.Wait()
	}
}

func (m *Monitor) pollDeviceStatus(address string) error {
	info, err := getRestInfo(address)
	if err != nil {
		logger.Printf("Monitor::pollDevice: REST failed, trying gRPC (%s): %v", address, err)
		info, err = getGrpcInfo(address)
		if err != nil {
			logger.Printf("Monitor::pollDevice: gRPC failed (%s): %v", address, err)
			return err
		}
	}

	if info == nil {
		logger.Printf("Monitor::pollDevice: empty info (%s)", address)
		return nil
	}

	m.devicesStatusMu.Lock()
	m.devicesStatus[address] = model.DeviceStatus{
		DeviceInfo: *info,
		UpdatedAt:  time.Now(),
	}
	m.devicesStatusMu.Unlock()

	err = m.devicesStatusStore.SaveJSON(&m.devicesStatus)
	if err != nil {
		logger.Printf("Monitor::pollDevice: devicesStatusStore.SaveJSON (%s): %v", address, err)
		return err
	}

	return nil
}
