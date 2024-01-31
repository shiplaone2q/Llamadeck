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
			if field.Bool() {
				args = append(args, "--"+flagName)
			}
		case reflect.Int:
			if field.Int() != 0 {
				args = append(args, "--"+flagName, strconv.FormatInt(field.Int(), 10))
			}
		case reflect.Float64:
			if field.Float() != 0 {
				args = append(args, "--"+flagName, strconv.FormatFloat(field.Float(), 'f', -1, 64))
			}
		case reflect.String:
			if field.String() != "" {
				args = append(args, "--"+flagName, field.String())
			}
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String && field.Len() > 0 {
				// Use jsonFieldName (snake_case) for multipleFlags lookup
				if _, isMultiValue := multipleFlags[jsonFieldName]; isMultiValue {
					// Multiple flags: --flag value1 --flag value2
					for j := 0; j < field.Len(); j++ {
