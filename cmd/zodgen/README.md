# zodgen - Zod Schema Generator

This tool generates TypeScript Zod schemas for the frontend from the Go struct definition.

## Usage

Generate the Zod schemas to stdout:

```bash
go run cmd/zodgen/main.go
```

Then manually copy the generated schemas into `webui/src/schemas/backends/llamacpp.ts`, replacing the schema definitions.

## How It Works

1. Parses `pkg/backends/llama.go` using Go's AST parser
