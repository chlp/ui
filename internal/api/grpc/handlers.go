package grpc

import (
	"context"
	"github.com/chlp/ui/internal/api/grpc/proto"
)

func (s *server) GetInfo(ctx context.Context, in *proto.Empty) (*proto.DeviceInfo, error) {
	return &proto.DeviceInfo{
		Id:              s.device.ID,
		HardwareVersion: s.device.HardwareVersion,
		SoftwareVersion: s.device.SoftwareVersion,
		FirmwareVersion: s.device.FirmwareVersion,
		Status:          string(s.device.Status),
		Checksum:        s.device.Checksum,
	}, nil
}
