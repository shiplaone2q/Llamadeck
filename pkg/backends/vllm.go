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

	// Performance and serving configuration
	BlockSize                int     `json:"block_size,omitempty"`
	EnablePrefixCaching      bool    `json:"enable_prefix_caching,omitempty"`
	DisableSlidingWindow     bool    `json:"disable_sliding_window,omitempty"`
	UseV2BlockManager        bool    `json:"use_v2_block_manager,omitempty"`
	NumLookaheadSlots        int     `json:"num_lookahead_slots,omitempty"`
	SwapSpace                int     `json:"swap_space,omitempty"`
	CPUOffloadGB             int     `json:"cpu_offload_gb,omitempty"`
	GPUMemoryUtilization     float64 `json:"gpu_memory_utilization,omitempty"`
	NumGPUBlocksOverride     int     `json:"num_gpu_blocks_override,omitempty"`
	MaxNumBatchedTokens      int     `json:"max_num_batched_tokens,omitempty"`
	MaxNumSeqs               int     `json:"max_num_seqs,omitempty"`
	MaxLogprobs              int     `json:"max_logprobs,omitempty"`
	DisableLogStats          bool    `json:"disable_log_stats,omitempty"`
	Quantization             string  `json:"quantization,omitempty"`
	RopeScaling              string  `json:"rope_scaling,omitempty"`
	RopeTheta                float64 `json:"rope_theta,omitempty"`
	EnforceEager             bool    `json:"enforce_eager,omitempty"`
	MaxContextLenToCapture   int     `json:"max_context_len_to_capture,omitempty"`
	MaxSeqLenToCapture       int     `json:"max_seq_len_to_capture,omitempty"`
	DisableCustomAllReduce   bool    `json:"disable_custom_all_reduce,omitempty"`
	TokenizerPoolSize        int     `json:"tokenizer_pool_size,omitempty"`
	TokenizerPoolType        string  `json:"tokenizer_pool_type,omitempty"`
	TokenizerPoolExtraConfig string  `json:"tokenizer_pool_extra_config,omitempty"`
	EnableLoraBias           bool    `json:"enable_lora_bias,omitempty"`
	LoraExtraVocabSize       int     `json:"lora_extra_vocab_size,omitempty"`
	LoraRank                 int     `json:"lora_rank,omitempty"`
	PromptLookbackDistance   int     `json:"prompt_lookback_distance,omitempty"`
	PreemptionMode           string  `json:"preemption_mode,omitempty"`

	// Distributed and parallel processing
	TensorParallelSize            int     `json:"tensor_parallel_size,omitempty"`
	PipelineParallelSize          int     `json:"pipeline_parallel_size,omitempty"`
	MaxParallelLoadingWorkers     int     `json:"max_parallel_loading_workers,omitempty"`
	DisableAsyncOutputProc        bool    `json:"disable_async_output_proc,omitempty"`
	WorkerClass                   string  `json:"worker_class,omitempty"`
	EnabledLoraModules            string  `json:"enabled_lora_modules,omitempty"`
	MaxLoraRank                   int     `json:"max_lora_rank,omitempty"`
	FullyShardedLoras             bool    `json:"fully_sharded_loras,omitempty"`
	LoraModules                   string  `json:"lora_modules,omitempty"`
	PromptAdapters                string  `json:"prompt_adapters,omitempty"`
	MaxPromptAdapterToken         int     `json:"max_prompt_adapter_token,omitempty"`
	Device                        string  `json:"device,omitempty"`
	SchedulerDelay                float64 `json:"scheduler_delay,omitempty"`
	EnableChunkedPrefill          bool    `json:"enable_chunked_prefill,omitempty"`
	SpeculativeModel              string  `json:"speculative_model,omitempty"`
	SpeculativeModelQuantization  string  `json:"speculative_model_quantization,omitempty"`
	SpeculativeRevision           string  `json:"speculative_revision,omitempty"`
	SpeculativeMaxModelLen        int     `json:"speculative_max_model_len,omitempty"`
	SpeculativeDisableByBatchSize int     `json:"speculative_disable_by_batch_size,omitempty"`
