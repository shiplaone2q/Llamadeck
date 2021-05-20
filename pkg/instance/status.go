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
	return json.Marshal(name)
}

// UnmarshalJSON implements json.Unmarshaler for Status enum
func (s *Status) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	status, ok := nameToStatus[str]
	if !ok {
		log.Printf("Unknown instance status: %s", str)
		status = Stopped // Default to Stopped on unknown status
	}

	*s = status
	return nil
}

// status represents the instance status with thread-safe access (unexported).
type status struct {
	mu sync.RWMutex
