FROM golang:1.24.2

WORKDIR /app

COPY . .

RUN go build -o app ./cmd/main.go

CMD ["./app"]