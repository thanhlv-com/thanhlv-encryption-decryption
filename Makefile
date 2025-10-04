.PHONY: build clean test install deps cross-compile

APP_NAME=thanhlv-ed
VERSION?=1.0.0

# Build for current platform
build:
	go build -ldflags="-s -w" -o bin/$(APP_NAME) main.go

# Clean build artifacts
clean:
	rm -rf build/ bin/

# Run tests
test:
	go test -v ./...

# Install dependencies
deps:
	go mod tidy
	go mod download

# Cross-platform compilation
build-all:
	./build.sh $(VERSION)

# Install binary to system
install: build
	cp bin/$(APP_NAME) /usr/local/bin/

# Development build (with debug info)
dev:
	go build -o bin/$(APP_NAME) main.go

# Create binary directory
bin:
	mkdir -p bin

# Default target
all: deps build