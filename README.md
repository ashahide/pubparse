# pubparse

**pubparse** is a command-line tool for parsing PubMed and PMC XML files, converting them into JSON format, and validating them against structured JSON Schema definitions. It supports parallel processing and generates detailed reports for each session.

---

## Features

- Supports both **PubMed** and **PMC** XML formats.
- Converts XML to compact **JSON**.
- **Schema validation** for structural correctness.
- **Concurrent processing** with a configurable number of workers.
- **Progress bar** and real-time console feedback.
- Generates a `report.tsv` summarizing input-output mappings.

---

## Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/yourname/pubparse.git
cd pubparse
go build -o pubparse ./cmd/pubparse
```

Or run directly:

```bash
go run ./cmd/pubparse [subcommand] -i input_dir -o output_dir [--workers N]
```

---

## Usage

```bash
pubparse [pubmed|pmc] -i <input_path> -o <output_path> [--workers N]
```

### Required Arguments

- `-i`: Path to an XML file or directory of XML files.
- `-o`: Path to the output directory for generated JSON files.

### Optional Flags

- `--workers`: Number of concurrent workers to use (default: 8). Caps at available CPU cores.

---

## Example

### Parse PubMed XML

```bash
pubparse pubmed -i ./test/pubmed_xml -o ./output/pubmed_json --workers 4
```

### Parse PMC Full Text XML

```bash
pubparse pmc -i ./test/pmc_xml -o ./output/pmc_json
```

---

## Output

Each execution generates:

- **JSON files** for each input XML file in the output directory.
- A `report.tsv` file inside the output directory that logs:
  - Start time
  - Input/output directories
  - Number of files processed
  - Worker count
  - Mapping of each input file to its output

---

## JSON Validation

Each JSON file is validated against:
- `pubmed_json_schema.json` for PubMed files
- `pmc_json_schema.json` for PMC files

You can find and edit these schema files under `internal/jsonTools/`.

---

## 🛠️ Development

To run tests:

```bash
go test ./internal/...
```

To rebuild:

```bash
go build -o pubparse ./cmd/pubparse
```

---

## Requirements

- Go 1.20+
- Unix-based OS (Linux/macOS) or Windows (PowerShell compatible)

---

## License

MIT License. See `LICENSE` file for details.

---

## Author

- Andrew Shahidehpour  
- GitHub: [@ashahide](https://github.com/ashahide)