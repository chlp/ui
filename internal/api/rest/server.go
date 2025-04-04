package rest

import (
	"encoding/json"
	"github.com/chlp/ui/internal/data"
	"github.com/chlp/ui/internal/model"
	"github.com/chlp/ui/pkg/logger"
	"net/http"
	"sync"
)

func StartRestServer(deviceConfig *model.DeviceInfo, devicesMu *sync.Mutex, devices *[]string, devicesFile string, restPort string) {
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(deviceConfig)
	})

	http.HandleFunc("/devices", func(w http.ResponseWriter, r *http.Request) {
		devicesMu.Lock()
		defer devicesMu.Unlock()
		_ = json.NewEncoder(w).Encode(devices)
	})

	http.HandleFunc("/add_device", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var payload struct {
			Address string `json:"address"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		devicesMu.Lock()
		defer devicesMu.Unlock()
		for _, d := range *devices {
			if d == payload.Address {
				w.WriteHeader(http.StatusConflict)
				return
			}
		}
		*devices = append(*devices, payload.Address)
		_ = saveDevices(devicesMu, devices, devicesFile)
		w.WriteHeader(http.StatusCreated)
	})

	http.HandleFunc("/remove_device", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var payload struct {
			Address string `json:"address"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		devicesMu.Lock()
		defer devicesMu.Unlock()
		var newDevices []string
		for _, d := range *devices {
			if d != payload.Address {
				newDevices = append(newDevices, d)
			}
		}
		devices = &newDevices
		_ = saveDevices(devicesMu, devices, devicesFile)
		w.WriteHeader(http.StatusOK)
	})

	logger.Printf("Starting REST server on %s", restPort)
	_ = http.ListenAndServe(restPort, nil)
}

func saveDevices(devicesMu *sync.Mutex, devices *[]string, devicesFile string) error {
	devicesMu.Lock()
	defer devicesMu.Unlock()
	return data.SaveJSON(devicesFile, &devices)
}
