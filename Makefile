.PHONY: all build build-all test test-coverage test-race test-single clean help install

all: build test

help:
	@echo "emojigate - Makefile commands:"
	@echo ""
	@echo "  make build           Build binary for current platform"
	@echo "  make build-all       Build binaries for all platforms"
	@echo "  make test            Run all tests"
	@echo "  make test-coverage   Run tests with coverage report"
	@echo "  make test-race       Run tests with race detector"
	@echo "  make test-single     Run specific test (TEST=TestName)"
	@echo "  make install         Install binary to GOPATH/bin"
	@echo "  make clean           Clean build artifacts"
	@echo "  make help            Show this help message"

build:
	@echo "Building emojigate..."
	@go build -ldflags="-s -w" -o emojigate.exe cmd/emojigate/main.go

build-all:
	@echo "Building for all platforms..."
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/emojigate-linux-amd64 cmd/emojigate/main.go
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/emojigate-linux-arm64 cmd/emojigate/main.go
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/emojigate-darwin-amd64 cmd/emojigate/main.go
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/emojigate-darwin-arm64 cmd/emojigate/main.go
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/emojigate-windows-amd64.exe cmd/emojigate/main.go
	@echo "Built all binaries in dist/"

install:
	@echo "Installing emojigate..."
	@go install -ldflags="-s -w" ./cmd/emojigate

test:
	@echo "Running tests..."
	@go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
	@echo "\nTo view HTML coverage report, run: go tool cover -html=coverage.out"

test-race:
	@echo "Running tests with race detector..."
	@go test -race ./...

test-single:
	@go test -v -run $(TEST) ./...

clean:
	@rm -f coverage.out
	@rm -f emojigate.exe
	@rm -rf dist/
	@echo "Cleaned up build and test artifacts"
