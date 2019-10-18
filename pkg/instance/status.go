package instance

import (
	"encoding/json"
	"log"
	"sync"
)

// Status is the enum for status values (exported).
type Status int

const (
	Stopped Status = iota
	Running
	Failed
	Restarting
	ShuttingDown
)

var nameToStatus = map[string]Status{
	"stopped":       Stopped,
	"running":       Running,
	"failed":        Failed,
	"restarting":    Restarting,
	"shutting_down": ShuttingDown,
}

var statusToName = map[Status]string{
	Stopped:      "stopped",
	Running:      "running",
	Failed:       "failed",
	Restarting:   "restarting",
	ShuttingDown: "shutting_down",
}

// Status enum JSON marshaling methods
func (s Status) MarshalJSON() ([]byte, error) {
	name, ok := statusToName[s]
	if !ok {
		name = "stopped" // Default to "stopped" for unknown status
	}
