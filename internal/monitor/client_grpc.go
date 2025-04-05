package monitor

import (
	"context"
	"github.com/chlp/ui/internal/api/grpc/proto"
	"github.com/chlp/ui/internal/model"
	"google.golang.org/grpc"
	"time"
)

func getGrpcInfo(address string) (*model.DeviceInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock()) // todo: deprecated
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := proto.NewDeviceServiceClient(conn)
	resp, err := client.GetInfo(ctx, &proto.Empty{})
	if err != nil {
		return nil, err
	}

	return &model.DeviceInfo{
		ID:              resp.Id,
		HardwareVersion: resp.HardwareVersion,
		SoftwareVersion: resp.SoftwareVersion,
		FirmwareVersion: resp.FirmwareVersion,
	}, nil
}
