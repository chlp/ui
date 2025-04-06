package monitor

import (
	"github.com/chlp/ui/internal/device"
	"github.com/chlp/ui/pkg/application"
	"sync"
	"testing"
	"time"
)

type MockStore struct {
	sync.Mutex
	Data    interface{}
	LoadErr error
	SaveErr error
}

func (m *MockStore) LoadJSON(v interface{}) error {
	m.Lock()
	defer m.Unlock()
	if m.LoadErr != nil {
		return m.LoadErr
	}

	switch data := m.Data.(type) {
	case map[string]string:
		p, ok := v.(*map[string]string)
		if ok {
			*p = data
		}
	case map[string]interface{}:
		p, ok := v.(*map[string]interface{})
		if ok {
			*p = data
		}
	case []string:
		p, ok := v.(*[]string)
		if ok {
			*p = data
		}
	default:
		if p, ok := v.(*map[string]device.Status); ok {
			*p = m.Data.(map[string]device.Status)
		}
	}

	return nil
}

func (m *MockStore) SaveJSON(v interface{}) error {
	m.Lock()
	defer m.Unlock()
	if m.SaveErr != nil {
		return m.SaveErr
	}
	m.Data = v
	return nil
}

func TestNewMonitor_Success(t *testing.T) {
	app, _ := application.NewApp("app", "", false)
	devicesListStore := &MockStore{Data: []string{"device1", "device2"}}
	devicesStatusStore := &MockStore{Data: map[string]device.Status{}}

	m, err := NewMonitor(app, devicesListStore, devicesStatusStore)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if m == nil {
		t.Fatal("expected monitor instance, got nil")
	}
}

func TestMustNewMonitor_Success(t *testing.T) {
	app, _ := application.NewApp("app", "", false)
	devicesListStore := &MockStore{Data: []string{"device1", "device2"}}
	devicesStatusStore := &MockStore{Data: map[string]device.Status{}}

	m := MustNewMonitor(app, devicesListStore, devicesStatusStore)
	if m == nil {
		t.Fatal("expected monitor instance, got nil")
	}
}

func TestNewMonitor_NilStores(t *testing.T) {
	app, _ := application.NewApp("app", "", false)

	m, err := NewMonitor(app, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if m != nil {
		t.Fatal("expected monitor to be nil if stores are nil")
	}
}

func TestGetDevicesStatus(t *testing.T) {
	app, _ := application.NewApp("app", "", false)
	devicesListStore := &MockStore{Data: []string{"device1"}}
	devicesStatusStore := &MockStore{Data: map[string]device.Status{
		"device1": {
			Info:      device.Info{ID: "device1"},
			UpdatedAt: time.Now(),
		},
	}}

	m, _ := NewMonitor(app, devicesListStore, devicesStatusStore)

	status := m.GetDevicesStatus()
	if len(status) != 1 {
		t.Fatalf("expected 1 device status, got %d", len(status))
	}
}
