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

		{"empty name", "", true},
		{"with spaces", "my instance", true},
		{"with special chars", "my@instance", true},
		{"too long", strings.Repeat("a", 51), true},
		{"shell metachar semicolon", "test;ls", true},
		{"shell metachar pipe", "test|ls", true},
		{"shell metachar ampersand", "test&ls", true},
		{"shell metachar dollar", "test$var", true},
		{"shell metachar backtick", "test`cmd`", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, err := validation.ValidateInstanceName(tt.input)
			if (err != nil) != tt.wantErr {
