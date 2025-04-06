package grpc

import (
	"github.com/chlp/ui/internal/api/grpc/proto"
	"github.com/chlp/ui/internal/model"
	"github.com/chlp/ui/pkg/application"
	"github.com/chlp/ui/pkg/logger"
	"google.golang.org/grpc"
	"net"
)

type server struct {
	device *model.DeviceInfo
	proto.UnimplementedDeviceServiceServer
}

func StartGrpcServer(app *application.App, port string, device *model.DeviceInfo) {
	if port == "" {
		return
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	proto.RegisterDeviceServiceServer(s, &server{device: device})

	app.Wg.Add(1)
	go func() {
		<-app.Ctx.Done()
		s.GracefulStop()
		app.Wg.Done()
	}()

	logger.Printf("StartGrpcServer: starting server on %s", port)
	if err = s.Serve(lis); err != nil {
		logger.Fatalf("StartGrpcServer: failed to serve: %v", err)
	}
}
