package monitor

import (
	"context"
	"github.com/chlp/ui/internal/api/grpc/proto"
	"github.com/chlp/ui/internal/device"
	"google.golang.org/grpc"
	"time"
)

const grpcClientTimeout = 2 * time.Second

func getGrpcInfo(address string) (*device.Info, error) {
	ctx, cancel := context.WithTimeout(context.Background(), grpcClientTimeout)
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

	return &device.Info{
		ID:              resp.Id,
		HardwareVersion: resp.HardwareVersion,
		SoftwareVersion: resp.SoftwareVersion,
		FirmwareVersion: resp.FirmwareVersion,
		Status:          device.StatusType(resp.Status),
		Checksum:        resp.Checksum,
	}, nil
}
