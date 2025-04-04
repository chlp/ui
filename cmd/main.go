package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	apiGrpc "github.com/chlp/ui/internal/api/grpc"
	"github.com/chlp/ui/internal/api/grpc/proto"
	"github.com/chlp/ui/internal/api/rest"
	"github.com/chlp/ui/internal/data"
	"github.com/chlp/ui/internal/model"
	"github.com/chlp/ui/pkg/logger"

	"google.golang.org/grpc"
)

const (
	configFile         = "config.json"
	stateFile          = "state.json"
	devicesFile        = "devices.json"
	logFile            = "app.log"
	pollInterval       = 3 * time.Second
	deviceScanInterval = 5 * time.Second
	grpcPort           = ":50051"
	restPort           = ":8080"
)

type State struct {
	LastSuccess map[string]time.Time `json:"last_success"`
	mu          sync.Mutex           `json:"-"`
}

var (
	deviceConfig model.DeviceInfo
	state        State
	devices      []string
	devicesMu    sync.Mutex
)

func loadOrCreateConfig() error {
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		deviceConfig = model.DeviceInfo{
			ID:              generateID(),
			Name:            "DefaultDevice",
			HardwareVersion: "1.0",
			SoftwareVersion: "1.0",
			FirmwareVersion: "1.0",
		}
		return data.SaveJSON(configFile, &deviceConfig)
	} else {
		return data.LoadJSON(configFile, &deviceConfig)
	}
}

func loadOrCreateState() error {
	if _, err := os.Stat(stateFile); errors.Is(err, os.ErrNotExist) {
		state = State{LastSuccess: make(map[string]time.Time)}
		return data.SaveJSON(stateFile, &state)
	} else {
		return data.LoadJSON(stateFile, &state)
	}
}

func generateID() string {
	return fmt.Sprintf("dev-%d", time.Now().UnixNano())
}

func getRestInfo(address string) (*model.DeviceInfo, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/info", address))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}
	var info model.DeviceInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

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
		Name:            resp.Name,
		HardwareVersion: resp.HardwareVersion,
		SoftwareVersion: resp.SoftwareVersion,
		FirmwareVersion: resp.FirmwareVersion,
	}, nil
}

func pollDevice(address string) {
	info, err := getRestInfo(address)
	if err != nil {
		logger.Printf("%s: REST failed, trying gRPC: %v", address, err)
		info, err = getGrpcInfo(address)
		if err != nil {
			logger.Printf("%s: gRPC failed: %v", address, err)
			return
		}
	}

	logger.Printf("%s: got info: %+v", address, info)
	state.mu.Lock()
	state.LastSuccess[address] = time.Now()
	state.mu.Unlock()
	data.SaveJSON(stateFile, &state)
}

func loadDevices() error {
	devicesMu.Lock()
	defer devicesMu.Unlock()
	if _, err := os.Stat(devicesFile); errors.Is(err, os.ErrNotExist) {
		devices = []string{}
		return data.SaveJSON(devicesFile, &devices)
	} else {
		return data.LoadJSON(devicesFile, &devices)
	}
}

func watchDevices() {
	ticker := time.NewTicker(deviceScanInterval)
	defer ticker.Stop()
	for range ticker.C {
		_ = loadDevices()
	}
}

func startGrpcServer() { // move to api/grpc
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterDeviceServiceServer(s, &apiGrpc.Server{DeviceConfig: &deviceConfig})
	logger.Printf("Starting gRPC server on %s", grpcPort)
	if err := s.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	logger.InitLogger(logFile)

	if err := loadOrCreateConfig(); err != nil {
		logger.Fatalf("failed to load/create config: %v", err)
	}
	if err := loadOrCreateState(); err != nil {
		logger.Fatalf("failed to load/create state: %v", err)
	}
	if err := loadDevices(); err != nil {
		logger.Fatalf("failed to load/create devices list: %v", err)
	}

	go watchDevices()
	go rest.StartRestServer(&deviceConfig, &devicesMu, &devices, devicesFile, restPort)
	go startGrpcServer()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for range ticker.C {
		devicesMu.Lock()
		currentDevices := append([]string{}, devices...)
		devicesMu.Unlock()
		for _, addr := range currentDevices {
			go pollDevice(addr)
		}
	}
}
