# llamadiff - llama-server Flag Validator

This tool validates that all flags supported by `llama-server` are properly handled by the llamadeck parser.

## Usage

Run the validator:

```bash
go run cmd/llamadiff/main.go
```

The tool will output a report showing:
- **Working**: Flags that are correctly parsed and mapped to struct fields
