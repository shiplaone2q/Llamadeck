# llamadiff - llama-server Flag Validator

This tool validates that all flags supported by `llama-server` are properly handled by the llamadeck parser.

## Usage

Run the validator:

```bash
go run cmd/llamadiff/main.go
```

The tool will output a report showing:
- **Working**: Flags that are correctly parsed and mapped to struct fields
- **Missing**: Flags from `llama-server --help` that are not yet supported
- **Ignored**: Utility flags (--help, --version, etc.) that don't need struct fields
- **Conflicts**: Fields that have both positive and negative flag variants mapping to the same struct field

## How It Works

1. Runs `llama-server --help` to extract all available flags
2. Parses the help output using regex to identify flag names
