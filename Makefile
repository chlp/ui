APP_NAME := ui_status_monitor
PROTO_FILE := internal/api/grpc/proto/device.proto

proto:
	protoc --go_out=. --go-grpc_out=. $(PROTO_FILE)

build-macos-arm:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $(APP_NAME)_macos_arm64 cmd/main.go

build-macos-intel:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(APP_NAME)_macos_amd64 cmd/main.go

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(APP_NAME)_linux_amd64 cmd/main.go

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(APP_NAME)_windows_amd64.exe cmd/main.go

build-all: build-macos-arm build-macos-intel build-linux build-windows

clean:
	rm -f $(APP_NAME) $(APP_NAME)_*

docker-build:
	docker build -t ui_monitor .

docker-run:
	docker run --rm -p 8080:8080 ui_monitor

docker: docker-build docker-run