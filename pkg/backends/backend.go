package backends

import (
	"encoding/json"
	"fmt"
	"llamadeck/pkg/config"
	"llamadeck/pkg/validation"
	"maps"
	"strings"
)

type BackendType string

const (
	BackendTypeLlamaCpp BackendType = "llama_cpp"
	BackendTypeMlxLm    BackendType = "mlx_lm"
	BackendTypeVllm     BackendType = "vllm"
	BackendTypeCustom   BackendType = "custom"
	BackendTypeUnknown  BackendType = "unknown"
)

type backend interface {
	BuildCommandArgs() []string
	BuildDockerArgs() []string
	GetModel() string
	GetPort() int
	SetPort(int)
	GetHost() string
	Validate() error
	ParseCommand(string) (any, error)
}

var backendConstructors = map[BackendType]func() backend{
	BackendTypeLlamaCpp: func() backend { return &LlamaServerOptions{} },
	BackendTypeMlxLm:    func() backend { return &MlxServerOptions{} },
	BackendTypeVllm:     func() backend { return &VllmServerOptions{} },
	BackendTypeCustom:   func() backend { return &CustomServerOptions{} },
}

type Options struct {
	BackendType    BackendType    `json:"backend_type"`
	BackendOptions map[string]any `json:"backend_options,omitempty"`

	// Backend-specific options
	LlamaServerOptions  *LlamaServerOptions  `json:"-"`
	MlxServerOptions    *MlxServerOptions    `json:"-"`
	VllmServerOptions   *VllmServerOptions   `json:"-"`
	CustomServerOptions *CustomServerOptions `json:"-"`
}

func (o *Options) UnmarshalJSON(data []byte) error {
	type Alias Options
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Create backend from constructor map
	constructor, exists := backendConstructors[o.BackendType]
	if !exists {
		return fmt.Errorf("unsupported backend type: %s", o.BackendType)
	}

	backend := constructor()

	// If backend_options is provided, unmarshal into the backend
	if o.BackendOptions != nil {
		optionsData, err := json.Marshal(o.BackendOptions)
		if err != nil {
			return fmt.Errorf("failed to marshal backend options: %w", err)
		}

		if err := json.Unmarshal(optionsData, backend); err != nil {
			return fmt.Errorf("failed to unmarshal backend options: %w", err)
		}
	}
	// If backend_options is nil or empty, backend remains as empty struct (for router mode)

	// Store in the appropriate typed field
	o.setBackendOptions(backend)

	return nil
}

func (o *Options) MarshalJSON() ([]byte, error) {
	// Get backend and marshal it
	var backendOptions map[string]any
	backend := o.getBackend()
	if backend != nil {
