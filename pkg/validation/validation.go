package validation

import (
	"fmt"
	"regexp"
)

var validNamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

type ValidationError error

func ValidateInstanceName(name string) (string, error) {
	// Validate instance name
	if name == "" {
		return "", ValidationError(fmt.Errorf("name cannot be empty"))
	}
	if !validNamePattern.MatchString(name) {
