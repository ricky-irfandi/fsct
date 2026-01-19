.PHONY: all build test lint clean install

BINARY_NAME := fsct
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)

all: build

build:
	go build -ldflags "$(LDFLAGS)" -o ./bin/$(BINARY_NAME) ./cmd/fsct/

test:
	go test ./... -coverprofile=coverage.out

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

clean:
	rm -f ./bin/$(BINARY_NAME)
	rm -f coverage.out coverage.html

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/fsct/

build-all:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./bin/$(BINARY_NAME)-linux-amd64 ./cmd/fsct/
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o ./bin/$(BINARY_NAME)-linux-arm64 ./cmd/fsct/
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./bin/$(BINARY_NAME)-darwin-amd64 ./cmd/fsct/
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o ./bin/$(BINARY_NAME)-darwin-arm64 ./cmd/fsct/
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/fsct/
