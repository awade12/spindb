.PHONY: build build-all test lint clean install dev-setup

BINARY_NAME=spindb
BUILD_DIR=build
LDFLAGS=-ldflags "-s -w"

build:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

build-all:
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