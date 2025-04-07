# Device Availability Monitoring Application

## Feature List

1. Retrieves a list of device addresses from the `devices_list.json` file.
2. Watches for changes in the `devices_list.json` file in real-time.
3. Allows managing the device list through the REST API.
4. Periodically pings devices via REST and gRPC APIs to update their statuses:
    1. Marks a device as `StatusUnavailable` after one minute of inactivity.
    2.  Marks a device as `StatusUnknown` if no data has ever been received from it.
5. Saves monitoring state to a file. Not a database, ensuring lightweight deployment without external dependencies.
6. Can emulate and act as a monitored device itself.
7. Allows enabling monitoring, REST API, and gRPC API, as well as configuring their ports via `config.json`.
8. Supports enabling debug mode through `config.json` for troubleshooting.
9. Allows specifying a path to a checksum application via `config.json` and can emulate checksum behavior if the application is unavailable.
10. Includes a Swagger specification for the REST API (device info and monitoring methods).
11. All resources (including Swagger files) are embedded into the compiled binary.
12. A `Makefile` is provided for building the application on popular systems (`make build-all`) and for running it via Docker (`make docker`).
13. Generated protobuf files and the `vendor` directory are committed to the repository for easier setup and reduced external dependencies.
14. Example unit tests available at `internal/monitor/monitor_test.go`.
15. Example integration tests available at `scenarios/start_scenarios.sh`:  
    Starts monitoring along with 5 simulated devices in different states.  
    Some devices are unreachable, one device appears and disappears from the device list, and one device stops responding after 30 seconds.

---

## Deployment Options

1. Download a ready-to-use binary for your platform.
2. Build the application manually: see the `Makefile` or run `make build-all`.
3. Run via Docker: execute `make docker`.

---

## Device List

Edit the `devices_list.json` file.  
The list is updated automatically **without restarting** the application.

Example:

```json
[
  "127.0.0.1:50052",
  "127.0.0.1:8083",
  "127.0.0.1:8084",
  "127.0.0.1:8085",
  "127.0.0.1:8089"
]
```

---

## Application Configuration

Edit the `config.json` file.  
Changes are applied **after restarting** the application.

Example:

```json
{
  "log_file": "app.log",
  "debug": false,
  "devices_list_file": "devices_list.json",
  "monitor_file": "monitor.json",
  "grpc_port": ":50051",
  "rest_port": ":8080",
  "device": {
    "id": "dev-1743983582705188000",
    "hardware_version": "0.0.1",
    "software_version": "1.0.0",
    "firmware_version": "0.0.5",
    "status": "ok",
    "checksum": ""
  },
  "checksum_cmd": "",
  "checksum_emulate": true
}
```

---

## REST API and Swagger Documentation

Swagger specification file:  
`internal/api/rest/swagger/swagger.yaml`

After starting the application (with `rest_port=":8080"`), Swagger UI is available at:  
[http://localhost:8080/swagger](http://localhost:8080/swagger)
