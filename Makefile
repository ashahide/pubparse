.PHONY: help build test clean

help:
	@echo ""
	@echo "Available targets:"
	@echo "  make build     Build the Go project"
	@echo "  make test      Run all tests"
	@echo "  make coverage  Run all tests and include coverage report"
	@echo "  make clean     Remove build artifacts"
	@echo "  make help      Show this help message"
	@echo ""
	
# Variables
GO=go

## Build the Go project
build:
	go build -o bin/pubparse ./cmd/pubparse


test:
	go test ./internal/... -cover

coverage:
	go test ./internal/... -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

## Remove build artifacts
clean:
	rm -rf bin
	rm -f coverage.out