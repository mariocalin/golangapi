MAIN_PACKAGE_PATH := ./cmd/server
MAIN_SERVER_FILE := api.go
BINARY_NAME := server

run: build run-bin

build:
	mkdir -p bin
	go build -o bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}/${MAIN_SERVER_FILE}

test:
	go test -v ./...

test/cover:
	mkdir -p tmp
	go test -v -coverprofile=tmp/coverage.out ./...
	go tool cover -func=tmp/coverage.out

run-bin:
	bin/${BINARY_NAME}