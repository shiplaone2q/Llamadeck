package validation_test

import (
	"llamadeck/pkg/backends"
	"llamadeck/pkg/validation"
	"strings"
	"testing"
)

func TestValidateInstanceName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"simple name", "myinstance", false},
		{"with numbers", "instance123", false},
		{"with hyphens", "my-instance", false},
		{"with underscores", "my_instance", false},
		{"with dots", "my.instance", false},
		{"mixed valid chars", "test-instance_123", false},
		{"single char", "a", false},
		{"max length", strings.Repeat("a", 50), false},
