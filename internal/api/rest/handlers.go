package rest

import (
	"encoding/json"
	"github.com/chlp/ui/pkg/logger"
	"net/http"
)

func (s *server) getInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := json.Marshal(s.device)
	if err != nil {
		logger.Printf("Rest::getInfoHandler: failed to marshal: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setContentTypeJson(w)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (s *server) getDevicesStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := json.Marshal(s.monitor.GetDevicesStatus())
	if err != nil {
		logger.Printf("Rest::getDevicesStatusHandler: failed to marshal: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setContentTypeJson(w)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (s *server) getDevicesListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := json.Marshal(s.monitor.GetDevicesList())
	if err != nil {
		logger.Printf("Rest::getDevicesListHandler: failed to marshal: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setContentTypeJson(w)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (s *server) deviceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.addDevice(w, r)
		return
	case http.MethodDelete:
		s.removeDevice(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *server) addDevice(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Address string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Printf("Rest::addDevice: %s", payload.Address)

	deleted, err := s.monitor.AddDevice(payload.Address)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	if !deleted {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) removeDevice(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Address string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Printf("Rest::removeDevice: %s", payload.Address)

	deleted, err := s.monitor.RemoveDevice(payload.Address)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	if !deleted {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func setContentTypeJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
