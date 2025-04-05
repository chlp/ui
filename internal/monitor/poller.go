package monitor

import (
	"github.com/chlp/ui/internal/model"
	"github.com/chlp/ui/pkg/logger"
	"sync"
	"time"
)

const (
	devicesStatusInfoPollInterval = 5 * time.Second
	maxParallelPolling            = 128
)

func (m *Monitor) pollAllDevicesStatus() {
	ticker := time.NewTicker(devicesStatusInfoPollInterval)
	defer ticker.Stop()
	for range ticker.C {
		m.devicesListMu.RLock()
		devicesList := make([]string, len(m.devicesList))
		copy(devicesList, m.devicesList)
		m.devicesListMu.RUnlock()

		wg := sync.WaitGroup{}
		semaphore := make(chan struct{}, maxParallelPolling)

		for _, address := range devicesList {
			wg.Add(1)
			semaphore <- struct{}{}
			go func() {
				defer wg.Done()
				defer func() { <-semaphore }()

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

	m.devicesStatusInfoMu.Lock()
	m.devicesStatusInfo[address] = model.DeviceStatusInfo{
		DeviceInfo: *info,
		UpdatedAt:  time.Now(),
	}
	m.devicesStatusInfoMu.Unlock()

	err = m.devicesStatusStore.SaveJSON(&m.devicesStatusInfo)
	if err != nil {
		logger.Printf("Monitor::pollDevice: devicesStatusStore.SaveJSON (%s): %v", address, err)
		return err
	}

	return nil
}
