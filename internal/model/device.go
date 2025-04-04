package model

type DeviceInfo struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	HardwareVersion string `json:"hardware_version"`
	SoftwareVersion string `json:"software_version"`
	FirmwareVersion string `json:"firmware_version"`
}
