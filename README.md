
# pubparse

**pubparse** is a command-line tool for parsing PubMed and PMC XML files, converting them into compact JSON, and validating them against structured JSON Schema definitions. It supports parallel processing, real-time progress tracking, and generates reports for reproducibility and auditing.

---

## Features

- Supports **PubMed** and **PMC** XML formats
- Converts to compact **JSON**
- **Schema validation** using JSON Schema
- **Parallel processing** with `--workers`
- Interactive **progress bar**
- Generates session-level `report.tsv` with file mappings

---

## Installation

Clone and build:

```bash
git clone https://github.com/yourname/pubparse.git
cd pubparse
make build
```

Or run directly with Go:

```bash
go run ./cmd/pubparse [pubmed|pmc] -i <input_path> -o <output_path> --workers 4
```

---

## Makefile Targets

```bash
make help           # Show usage and available targets
make build          # Build binary at ./bin/pubparse
make test           # Run tests with coverage info
make coverage       # Generate HTML test coverage report
make data           # Generate test XML data for PubMed/PMC
make example        # Run example pipeline (add FETCH=1 to refresh data)
make clean          # Remove bin and coverage artifacts
```

---

## Usage

```bash
pubparse [pubmed|pmc] -i <input_path> -o <output_path> [--workers N]
```

### Required Flags

- `-i`: Path to a single XML file or directory of files
- `-o`: Output directory for JSON files

### Optional Flags

- `--workers`: Number of concurrent workers (default: 8, capped at CPU cores)

---

## Example

```bash
make example FETCH=1
```

This will:
1. Download test data (`fetch_test_data.sh`)
2. Run both PubMed and PMC processing
3. Save output JSONs and `report.tsv` under `test/data/test_pubmed/json/` and `test/data/test_pmc/json/`

---

## Output

Each run produces:

- JSON files for each XML input
- A `report.tsv` containing:
  - Timestamp
  - Input/output paths
  - File count
  - Worker count
  - Per-file conversion status

---

## JSON Schema Validation

All JSON outputs are validated against schemas:

- `pubmed_json_schema.json` for PubMed
- `pmc_json_schema.json` for PMC

Schemas live in:

```bash
internal/jsonTools/
```

---

## Testing

Run tests and generate coverage:

```bash
make test
make coverage
```

---

## Requirements

- Go 1.20+
- Bash or POSIX-compatible shell
- Internet connection (for test data generation)

---

## License

MIT License. See `LICENSE` for details.

---

## Author

**Andrew Shahidehpour**  
GitHub: [@ashahide](https://github.com/ashahide)