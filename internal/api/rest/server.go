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

	http.HandleFunc("/v1/info", s.getInfo)

	if s.monitor != nil {
		http.HandleFunc("/v1/devices_status", s.getDevicesStatus)
		http.HandleFunc("/v1/devices_list", s.getDevicesList)
		http.HandleFunc("/v1/add_device", s.addDevice)
		http.HandleFunc("/v1/remove_device", s.removeDevice)
	}

	httpServer := &http.Server{Addr: port, Handler: nil}
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
