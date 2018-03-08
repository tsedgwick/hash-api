
TARGET=server

all: buildServer test

build:
	go build -o bin/$(TARGET) cmd/$(TARGET)/main.go

run:
	go run cmd/$(TARGET)/main.go $(MAKECMDGOALS)


test:
	go test -v ./...
