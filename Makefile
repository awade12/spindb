.PHONY: build build-all test lint clean install dev-setup

BINARY_NAME=spindb
BUILD_DIR=build

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

LDFLAGS=-ldflags "-s -w \
	-X 'github.com/awade12/spindb/cmd.Version=$(VERSION)' \
	-X 'github.com/awade12/spindb/cmd.GitCommit=$(GIT_COMMIT)' \
	-X 'github.com/awade12/spindb/cmd.BuildDate=$(BUILD_DATE)'"

build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

build-all:
	@echo "Building $(BINARY_NAME) $(VERSION) for all platforms..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -rf $(BUILD_DIR)

install: build
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

dev-setup:
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest 