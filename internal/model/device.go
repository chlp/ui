package model

import "time"

type DeviceStatusType string

const (
	DeviceStatusUnknown DeviceStatusType = "unknown"
	DeviceStatusOk      DeviceStatusType = "ok"
	DeviceStatusWarning DeviceStatusType = "warning"
	DeviceStatusFatal   DeviceStatusType = "fatal"
)

type DeviceInfo struct {
	ID              string           `json:"id"`
	HardwareVersion string           `json:"hardware_version"`
	SoftwareVersion string           `json:"software_version"`
	FirmwareVersion string           `json:"firmware_version"`
	Status          DeviceStatusType `json:"status"`
	Checksum        string           `json:"checksum"`
}

type DeviceStatus struct {
	DeviceInfo
	UpdatedAt time.Time `json:"last_success"`
}
