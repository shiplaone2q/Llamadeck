package testutil

import "slices"

// Helper functions for pointer fields
func BoolPtr(b bool) *bool {
	return &b
}

func IntPtr(i int) *int {
	return &i
