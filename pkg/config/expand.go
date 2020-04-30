package config

import (
	"os"
	"regexp"
)

var envVarPattern = regexp.MustCompile(`\$\{([a-zA-Z_][a-zA-Z0-9_]*)(?::-([^}]*))?\}`)

func expandEnvVars(data []byte) []byte {
