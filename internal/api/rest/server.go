package rest

import (
	"github.com/chlp/ui/internal/model"
	"github.com/chlp/ui/pkg/logger"
	"net/http"
)

type server struct {
	device  *model.DeviceInfo
	monitor Monitor
}

func StartRestServer(port string, device *model.DeviceInfo, monitor Monitor) {
	if port == "" {
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

	logger.Printf("StartRestServer: starting server on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		logger.Fatalf("StartRestServer: failed to serve: %v", err)
	}
}
