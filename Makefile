.PHONY: help build test coverage clean gen-data

# ------------------------ Help ------------------------

help:
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  build       Build the Go project (outputs to ./bin/pubparse)"
	@echo "  test        Run all unit tests with coverage info"
	@echo "  coverage    Generate HTML test coverage report"
	@echo "  gen-data    Generate sample PubMed and PMC XML files for testing"
	@echo "  clean       Remove build artifacts and temporary files"
	@echo "  help        Show this help message"
	@echo ""

# ------------------------ Variables ------------------------

GO := go
BIN := bin/pubparse

# ------------------------ Build ------------------------

## Build the Go project
build:
	@echo ">>> Building pubparse..."
	$(GO) build -o $(BIN) ./cmd/pubparse

# ------------------------ Test ------------------------

## Run tests with basic coverage
test:
	@echo ">>> Running tests..."
	$(GO) test ./internal/... -cover

## Generate and open coverage report
coverage:
	@echo ">>> Generating coverage report..."
	$(GO) test ./internal/... -coverprofile=coverage.out
	@echo ">>> Opening HTML coverage viewer..."
	$(GO) tool cover -html=coverage.out

# ------------------------ Test Data ------------------------

## Generate test data from PubMed and PMC
gen-data:
	@echo ">>> Generating test data from PubMed and PMC..."
	bash scripts/generate_test_data.sh

# ------------------------ Clean ------------------------

## Clean build and test artifacts
clean:
	@echo ">>> Cleaning up..."
	rm -rf bin
	rm -f coverage.out
