FROM golang:1.22

RUN apt-get update && apt-get install -y protobuf-compiler

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN protoc --go_out=. --go-grpc_out=. internal/api/grpc/proto/device.proto

RUN go build -o app ./cmd

CMD ["./app"]