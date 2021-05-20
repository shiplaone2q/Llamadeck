package validation

import (
	"fmt"
	"regexp"
)

var validNamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

type ValidationError error

func ValidateInstanceName(name string) (string, error) {
