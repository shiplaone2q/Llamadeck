package backends

import (
	"encoding/json"
	"fmt"
	"llamadeck/pkg/validation"
)

type MlxServerOptions struct {
	// Basic connection options
	Model string `json:"model,omitempty"`
	Host  string `json:"host,omitempty"`
	Port  int    `json:"port,omitempty"`

	// Model and adapter options
	AdapterPath     string `json:"adapter_path,omitempty"`
	DraftModel      string `json:"draft_model,omitempty"`
	NumDraftTokens  int    `json:"num_draft_tokens,omitempty"`
	TrustRemoteCode bool   `json:"trust_remote_code,omitempty"`

	// Logging and templates
	LogLevel               string `json:"log_level,omitempty"`
	ChatTemplate           string `json:"chat_template,omitempty"`
	UseDefaultChatTemplate bool   `json:"use_default_chat_template,omitempty"`
	ChatTemplateArgs       string `json:"chat_template_args,omitempty"` // JSON string

	// Sampling defaults
	Temp      float64 `json:"temp,omitempty"`
	TopP      float64 `json:"top_p,omitempty"`
	TopK      int     `json:"top_k,omitempty"`
	MinP      float64 `json:"min_p,omitempty"`
	MaxTokens int     `json:"max_tokens,omitempty"`

	// ExtraArgs are additional command line arguments.
	// Example: {"verbose": "", "log-file": "/logs/mlx.log"}
	ExtraArgs map[string]string `json:"extra_args,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling to collect unknown fields into ExtraArgs
func (o *MlxServerOptions) UnmarshalJSON(data []byte) error {
	// First unmarshal into a map to capture all fields
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Create a temporary struct for standard unmarshaling
	type tempOptions MlxServerOptions
	temp := tempOptions{}

	// Standard unmarshal first
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Copy to our struct
	*o = MlxServerOptions(temp)

	// Get all known canonical field names from struct tags
	knownFields := getKnownFieldNames(o)

	// Collect unknown fields into ExtraArgs
	if o.ExtraArgs == nil {
