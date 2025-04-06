package device

import "time"

type StatusType string

const (
	StatusUnknown     StatusType = "unknown"
	StatusOk          StatusType = "ok"
	StatusWarning     StatusType = "warning"
	StatusFatal       StatusType = "fatal"
	StatusUnavailable StatusType = "unavailable"
)

type Info struct {
	ID              string     `json:"id"`
	HardwareVersion string     `json:"hardware_version"`
	SoftwareVersion string     `json:"software_version"`
	FirmwareVersion string     `json:"firmware_version"`
	Status          StatusType `json:"status"`
	Checksum        string     `json:"checksum"`
}

type Status struct {
	Info
	UpdatedAt time.Time `json:"updated_at"`
}
