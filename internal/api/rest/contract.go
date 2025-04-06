package rest

import "github.com/chlp/ui/internal/device"

type Monitor interface {
	GetDevicesStatus() map[string]device.Status
	GetDevicesList() []string
	AddDevice(address string) (bool, error)
	RemoveDevice(address string) (bool, error)
}
