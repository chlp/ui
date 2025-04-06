package rest

import "github.com/chlp/ui/internal/model"

type Monitor interface {
	GetDevicesStatus() map[string]model.DeviceStatus
	GetDevicesList() []string
	AddDevice(address string) (bool, error)
	RemoveDevice(address string) (bool, error)
}
