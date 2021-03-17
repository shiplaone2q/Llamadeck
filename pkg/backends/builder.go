package backends

import (
	"fmt"
	"llamadeck/pkg/config"
	"reflect"
	"strconv"
	"strings"
)

// BuildCommandArgs converts a struct to command line arguments
func BuildCommandArgs(options any, multipleFlags map[string]struct{}) []string {
	var args []string

	v := reflect.ValueOf(options).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanInterface() {
			continue
		}

		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Get flag name from JSON tag (snake_case)
		jsonFieldName := strings.Split(jsonTag, ",")[0]
		// Convert to kebab-case for CLI flags
		flagName := strings.ReplaceAll(jsonFieldName, "_", "-")

		switch field.Kind() {
		case reflect.Bool:
