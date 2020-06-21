package backends

import (
	"encoding/json"
	"fmt"
	"llamadeck/pkg/validation"
)

// vllmMultiValuedFlags defines flags that should be repeated for each value rather than comma-separated
// Based on vLLM's CLI argument definitions with action='append' or List types
// Keys use snake_case as the parser converts kebab-case flags to snake_case before lookup
var vllmMultiValuedFlags = map[string]struct{}{
	"api_key":         {}, // --api-key (action='append')
	"allowed_origins": {}, // --allowed-origins (List type)
	"allowed_methods": {}, // --allowed-methods (List type)
	"allowed_headers": {}, // --allowed-headers (List type)
	"middleware":      {}, // --middleware (action='append')
	"lora_modules":    {}, // --lora-modules (custom LoRAParserAction, accepts multiple)
	"prompt_adapters": {}, // --prompt-adapters (similar to lora-modules, accepts multiple)
}

type VllmServerOptions struct {
	// Basic connection options (auto-assigned by llamadeck)
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`

	// Model and engine configuration
	Model                      string `json:"model,omitempty"`
	Tokenizer                  string `json:"tokenizer,omitempty"`
	SkipTokenizerInit          bool   `json:"skip_tokenizer_init,omitempty"`
	Revision                   string `json:"revision,omitempty"`
	CodeRevision               string `json:"code_revision,omitempty"`
	TokenizerRevision          string `json:"tokenizer_revision,omitempty"`
	TokenizerMode              string `json:"tokenizer_mode,omitempty"`
	TrustRemoteCode            bool   `json:"trust_remote_code,omitempty"`
	DownloadDir                string `json:"download_dir,omitempty"`
	LoadFormat                 string `json:"load_format,omitempty"`
	ConfigFormat               string `json:"config_format,omitempty"`
	Dtype                      string `json:"dtype,omitempty"`
	KVCacheDtype               string `json:"kv_cache_dtype,omitempty"`
	QuantizationParamPath      string `json:"quantization_param_path,omitempty"`
	Seed                       int    `json:"seed,omitempty"`
	MaxModelLen                int    `json:"max_model_len,omitempty"`
	GuidedDecodingBackend      string `json:"guided_decoding_backend,omitempty"`
	DistributedExecutorBackend string `json:"distributed_executor_backend,omitempty"`
	WorkerUseRay               bool   `json:"worker_use_ray,omitempty"`
	RayWorkersUseNSight        bool   `json:"ray_workers_use_nsight,omitempty"`

