.PHONY: help build test coverage clean data example

# ------------------------ Help ------------------------

help:
	@echo ""
	@echo "Usage: make [target] [FETCH=1]"
	@echo ""
	@echo "Available targets:"
	@echo "  build       Build the Go project (outputs to ./bin/pubparse)"
	@echo "  test        Run all unit tests with coverage info"
	@echo "  coverage    Generate HTML test coverage report"
	@echo "  data        Generate test PubMed and PMC XML files (via generate_test_data.sh)"
	@echo "  example     Run example pipeline; use FETCH=1 to refresh test data"
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

# ------------------------ Example ------------------------

## Run example pipeline with optional FETCH=1 to refresh test data
example:
	@echo ">>> Running example with test data..."
ifeq ($(FETCH),1)
	@echo ">>> FETCH=1: Regenerating test data..."
	bash fetch_test_data.sh
else
	@echo ">>> Skipping test data generation (set FETCH=1 to enable)"
endif

	@echo ">>> Running PMC example..."
	$(GO) run ./cmd/pubparse pmc -i test/data/test_pmc/xml/ -o test/data/test_pmc/json/

	@echo ">>> Running PubMed example..."
	$(GO) run ./cmd/pubparse pubmed -i test/data/test_pubmed/xml/ -o test/data/test_pubmed/json/

	@echo ">>> Example completed. Output written to:"
	@echo "    test/data/test_pubmed/json/"
	@echo "    test/data/test_pmc/json/"

# ------------------------ Test Data ------------------------

## Generate test data from PubMed and PMC
data:
	@echo ">>> Generating test data from PubMed and PMC..."
	bash scripts/generate_test_data.sh

# ------------------------ Clean ------------------------

## Clean build and test artifacts
clean:
	@echo ">>> Cleaning up..."
	rm -rf bin
	rm -f coverage.out
