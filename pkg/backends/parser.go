package backends

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// parseCommand parses a command string into a target struct
func parseCommand(command string, executableNames []string, subcommandNames []string, multiValuedFlags map[string]struct{}, target any) error {
	return parseCommandWithAliases(command, executableNames, subcommandNames, multiValuedFlags, nil, target)
}

// parseCommandWithAliases parses a command string into a target struct with support for flag aliases
func parseCommandWithAliases(command string, executableNames []string, subcommandNames []string, multiValuedFlags map[string]struct{}, flagAliases map[string]string, target any) error {
	// Normalize multiline commands
	command = normalizeCommand(command)
	if command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	// Extract arguments and positional model
	args, modelFromPositional, err := extractArgs(command, executableNames, subcommandNames)
	if err != nil {
		return err
	}

	// Parse flags into map
	options, err := parseFlags(args, multiValuedFlags, flagAliases)
	if err != nil {
		return err
	}

	// If we found a positional model and no --model flag was provided, set the model
	if modelFromPositional != "" {
		if _, hasModelFlag := options["model"]; !hasModelFlag {
			options["model"] = modelFromPositional
		}
	}

	// Convert to target struct via JSON
	jsonData, err := json.Marshal(options)
	if err != nil {
		return fmt.Errorf("failed to marshal options: %w", err)
	}

	if err := json.Unmarshal(jsonData, target); err != nil {
		return fmt.Errorf("failed to unmarshal to target: %w", err)
	}

	return nil
}

// normalizeCommand handles multiline commands with backslashes
func normalizeCommand(command string) string {
	re := regexp.MustCompile(`\\\s*\n\s*`)
	normalized := re.ReplaceAllString(command, " ")
	re = regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(normalized, " "))
}

// extractArgs extracts arguments from command, removing executable and subcommands
// Returns: args, modelFromPositional, error
func extractArgs(command string, executableNames []string, subcommandNames []string) ([]string, string, error) {
	// Check for unterminated quotes
	if strings.Count(command, `"`)%2 != 0 || strings.Count(command, `'`)%2 != 0 {
		return nil, "", fmt.Errorf("unterminated quoted string")
	}

	tokens := strings.Fields(command)
	if len(tokens) == 0 {
		return nil, "", fmt.Errorf("no tokens found")
	}

	// Skip executable
	start := 0
	firstToken := tokens[0]

	// Check for executable name (with or without path)
	if strings.Contains(firstToken, string(filepath.Separator)) {
		baseName := filepath.Base(firstToken)
		for _, execName := range executableNames {
			if strings.HasSuffix(strings.ToLower(baseName), strings.ToLower(execName)) {
				start = 1
				break
			}
		}
	} else {
