# Gendiff

A Go CLI utility for comparing configuration files and showing differences in various formats.

## Features

- **Multiple file formats**: Supports JSON and YAML files
- **Multiple output formats**: 
  - `stylish` (default) - Human-readable diff with +/- indicators
  - `plain` - Simple text descriptions of changes
  - `json` - Structured JSON output
- **Recursive comparison**: Handles nested objects and arrays
- **Cross-platform**: Works on Windows, macOS, and Linux

## Installation

### Prerequisites
- Go 1.21 or higher

### Build from source
```bash
git clone <repository-url>
cd go-test-project-2441
make build
```

## Usage

### Basic usage
```bash
./bin/gendiff file1.json file2.json
```

### Specify output format
```bash
./bin/gendiff -f plain file1.json file2.json
./bin/gendiff --format json file1.yml file2.yml
```

### Help
```bash
./bin/gendiff --help
```

## Examples

### Input files

**file1.json:**
```json
{
  "host": "hexlet.io",
  "timeout": 50,
  "proxy": "123.234.53.22",
  "follow": false
}
```

**file2.json:**
```json
{
  "timeout": 20,
  "verbose": true,
  "host": "hexlet.io"
}
```

### Output formats

#### Stylish (default)
```bash
./bin/gendiff file1.json file2.json
```
Output:
```
{
  - follow: false
    host: "hexlet.io"
  - proxy: "123.234.53.22"
  - timeout: 50
  + timeout: 20
  + verbose: true
}
```

#### Plain
```bash
./bin/gendiff -f plain file1.json file2.json
```
Output:
```
Property 'follow' was removed
Property 'proxy' was removed
Property 'timeout' was updated. From 50 to 20
Property 'verbose' was added with value: true
```

#### JSON
```bash
./bin/gendiff -f json file1.json file2.json
```
Output:
```json
{
  "type": "root",
  "children": [
    {
      "type": "removed",
      "key": "follow",
      "oldValue": false
    },
    {
      "type": "unchanged",
      "key": "host",
      "value": "hexlet.io"
    },
    {
      "type": "removed",
      "key": "proxy",
      "oldValue": "123.234.53.22"
    },
    {
      "type": "updated",
      "key": "timeout",
      "oldValue": 50,
      "newValue": 20
    },
    {
      "type": "added",
      "key": "verbose",
      "newValue": true
    }
  ]
}
```

## Development

### Project structure
```
.
├── cmd/gendiff/          # CLI application
├── internal/             # Internal packages
├── testdata/             # Test fixtures
├── gendiff.go            # Main library code
├── gendiff_test.go       # Tests
├── Makefile              # Build commands
└── go.mod                # Go module file
```

### Available make targets
```bash
make build              # Build the binary
make test               # Run tests
make test-with-coverage # Run tests with coverage report
make lint               # Run linter
make clean              # Clean build artifacts
```

### Running tests
```bash
go test -v              # Run all tests
go test -cover          # Run tests with coverage
go test -race           # Run tests with race detection
```

## API

### Library usage
```go
package main

import "code"

func main() {
    result, err := code.GenDiff("file1.json", "file2.json", "stylish")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result)
}
```

## Supported file formats

- **JSON**: Files with `.json` extension
- **YAML**: Files with `.yml` or `.yaml` extension

## Error handling

The utility provides clear error messages for common issues:
- File not found
- Unsupported file format
- Invalid JSON/YAML syntax
- Unsupported output format

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License.

