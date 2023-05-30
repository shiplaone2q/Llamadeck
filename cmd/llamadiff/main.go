package main

import (
	"fmt"
	"llamadeck/pkg/backends"
	"os/exec"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

type FlagInfo struct {
	Flag        string
	Description string
	EnvVar      string
}

// Utility flags that are not expected to be in the struct
var utilityFlags = map[string]bool{
	"-h":                true,
	"--help":            true,
	"--usage":           true,
	"--version":         true,
	"--license":         true,
	"-cl":               true,
	"--cache-list":      true,
	"--completion-bash": true,
	"--list-devices":    true,
}

func main() {
	// Run llama-server --help to get all flags
	fmt.Println("Running llama-server --help to extract flags...")
	helpOutput, err := runLlamaServerHelp()
	if err != nil {
		fmt.Printf("Error running llama-server --help: %v\n", err)
		return
	}

	// Parse help output to extract all flags
	flags := parseHelpFlags(helpOutput)
	fmt.Printf("Found %d flags in llama-server --help\n\n", len(flags))

	// Test each flag with the actual parser
	var missing []FlagInfo
	var working []string
	var ignored []string
	fieldToFlags := make(map[string][]string) // Track which flags set which fields

	for _, flag := range flags {
		// Skip utility flags
		if utilityFlags[flag.Flag] {
			ignored = append(ignored, flag.Flag)
			continue
		}

		fieldName, ok := testFlagAndGetField(flag.Flag)
		if ok {
			working = append(working, flag.Flag)
			if fieldName != "" {
				fieldToFlags[fieldName] = append(fieldToFlags[fieldName], flag.Flag)
			}
		} else {
			missing = append(missing, flag)
		}
	}

	// Check for conflicts (positive/negative pairs setting the same field)
	conflicts := findConflicts(fieldToFlags)

	// Report results
	fmt.Println("=== FLAGS MISSING FROM llamadeck ===")
	if len(missing) == 0 {
