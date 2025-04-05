package grpc

import (
	"github.com/chlp/ui/internal/api/grpc/proto"
	"github.com/chlp/ui/internal/config"
	"github.com/chlp/ui/internal/model"
	"github.com/chlp/ui/pkg/logger"
	"google.golang.org/grpc"
	"net"
)

type server struct {
	DeviceConfig *model.DeviceInfo
	proto.UnimplementedDeviceServiceServer
}

func StartGrpcServer(cfg *config.Config) {
	if cfg.GrpcPort == "" {
		return
	}

	lis, err := net.Listen("tcp", cfg.GrpcPort)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	proto.RegisterDeviceServiceServer(s, &server{DeviceConfig: cfg.Device})

	logger.Printf("StartGrpcServer: starting server on %s", cfg.GrpcPort)
	if err := s.Serve(lis); err != nil {
		logger.Fatalf("StartGrpcServer: failed to serve: %v", err)
	}
}
