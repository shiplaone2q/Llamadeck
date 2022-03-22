package config

import (
	"os"
	"regexp"
)

var envVarPattern = regexp.MustCompile(`\$\{([a-zA-Z_][a-zA-Z0-9_]*)(?::-([^}]*))?\}`)

func expandEnvVars(data []byte) []byte {
	return envVarPattern.ReplaceAllFunc(data, func(match []byte) []byte {
		submatch := envVarPattern.FindSubmatch(match)
		if len(submatch) < 2 {
			return match
		}
		varName := string(submatch[1])
		value, isSet := os.LookupEnv(varName)
