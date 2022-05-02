package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type FieldInfo struct {
	Name    string
	Type    string
	JSONTag string
	Comment string
	AltKeys []string
}

func main() {
	// Parse the llama.go file
	fset := token.NewFileSet()
	sourceFile := filepath.Join("pkg", "backends", "llama.go")

	file, err := parser.ParseFile(fset, sourceFile, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
		os.Exit(1)
	}

	// Find LlamaServerOptions struct
	var fields []FieldInfo
	ast.Inspect(file, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok || typeSpec.Name.Name != "LlamaServerOptions" {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return false
		}

		// Extract fields with section headers
		for _, field := range structType.Fields.List {
			// Check for section header in Doc comments
			if field.Doc != nil {
				for _, comment := range field.Doc.List {
					text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
					if text != "" && (strings.Contains(text, "params") || strings.Contains(text, "Params") || strings.Contains(text, "Parameters")) {
						fields = append(fields, FieldInfo{
							Comment: text,
						})
					}
				}
			}

			if len(field.Names) == 0 {
				continue // Skip embedded fields
			}

			fieldName := field.Names[0].Name

			// Skip unexported fields
			if !ast.IsExported(fieldName) {
				continue
			}

			// Extract JSON tag
			var jsonTag string
			if field.Tag != nil {
				tag := strings.Trim(field.Tag.Value, "`")
				if strings.Contains(tag, "json:") {
					parts := strings.Split(tag, "\"")
					if len(parts) >= 2 {
						jsonTag = strings.Split(parts[1], ",")[0]
					}
				}
			}

			// Skip fields without JSON tags or with "-"
			if jsonTag == "" || jsonTag == "-" {
				continue
			}

			// Extract type
			goType := typeToString(field.Type)

			// Extract inline comment
			var comment string
			if field.Comment != nil && len(field.Comment.List) > 0 {
				comment = strings.TrimSpace(strings.TrimPrefix(field.Comment.List[0].Text, "//"))
			}

			// Extract alternative keys from comment
