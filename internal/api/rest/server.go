package rest

import (
	"github.com/chlp/ui/internal/model"
	"github.com/chlp/ui/pkg/application"
	"github.com/chlp/ui/pkg/logger"
	"net/http"
)

type server struct {
	device  *model.DeviceInfo
	monitor Monitor
}

func StartServer(app *application.App, port string, device *model.DeviceInfo, monitor Monitor) {
	if port == "" {
		logger.Printf("Rest::StartServer: starting without rest server")
		return
	}

	s := &server{
		device:  device,
		monitor: monitor,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/info", s.getInfoHandler)

	if s.monitor != nil {
		mux.HandleFunc("/v1/devices_status", s.getDevicesStatusHandler)
		mux.HandleFunc("/v1/devices", s.getDevicesListHandler)
		mux.HandleFunc("/v1/device", s.deviceHandler)
	}

	serveSwaggerFiles(mux)

	httpServer := &http.Server{Addr: port, Handler: mux}
	app.Wg.Add(1)
	go func() {
		<-app.Ctx.Done()
		_ = httpServer.Close()
		app.Wg.Done()
	}()

	logger.Printf("Rest::StartServer: starting server on %s", port)
	if err := httpServer.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			logger.Fatalf("Rest::StartServer: failed to serve: %v", err)
		}
	}
}
