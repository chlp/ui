package model

import "time"

type DeviceStatus string

const (
	DeviceStatusUnknown DeviceStatus = "unknown"
	DeviceStatusOk      DeviceStatus = "ok"
	DeviceStatusWarning DeviceStatus = "warning"
	DeviceStatusFatal   DeviceStatus = "fatal"
)

type DeviceInfo struct {
	ID              string       `json:"id"`
	HardwareVersion string       `json:"hardware_version"`
	SoftwareVersion string       `json:"software_version"`
	FirmwareVersion string       `json:"firmware_version"`
	Status          DeviceStatus `json:"status"`
	Checksum        string       `json:"checksum"`
}

type DeviceStatusInfo struct {
	DeviceInfo
	UpdatedAt time.Time `json:"last_success"`
}
