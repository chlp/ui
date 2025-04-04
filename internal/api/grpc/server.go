package grpc

import (
	"context"
	"github.com/chlp/ui/internal/api/grpc/proto"
	"github.com/chlp/ui/internal/model"
)

type Server struct {
	DeviceConfig *model.DeviceInfo
	proto.UnimplementedDeviceServiceServer
}

func (s *Server) GetInfo(ctx context.Context, in *proto.Empty) (*proto.DeviceInfo, error) {
	return &proto.DeviceInfo{
		Id:              s.DeviceConfig.ID,
		Name:            s.DeviceConfig.Name,
		HardwareVersion: s.DeviceConfig.HardwareVersion,
		SoftwareVersion: s.DeviceConfig.SoftwareVersion,
		FirmwareVersion: s.DeviceConfig.FirmwareVersion,
	}, nil
}
