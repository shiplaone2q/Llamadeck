package backends

import (
	"encoding/json"
	"fmt"
	"llamadeck/pkg/validation"
	"reflect"
)

// llamaMultiValuedFlags defines flags that should be repeated for each value rather than comma-separated
// Keys use snake_case as the parser converts kebab-case flags to snake_case before lookup
var llamaMultiValuedFlags = map[string]struct{}{
	"override_tensor":       {},
	"override_tensor_draft": {},
	"override_kv":           {},
	"lora":                  {},
	"lora_scaled":           {},
	"control_vector":        {},
	"control_vector_scaled": {},
	"dry_sequence_breaker":  {},
	"logit_bias":            {},
}

type LlamaServerOptions struct {
	// Common params
	VerbosePrompt           bool     `json:"verbose_prompt,omitempty"`
	Threads                 int      `json:"threads,omitempty"`         // -t, --threads N
	ThreadsBatch            int      `json:"threads_batch,omitempty"`   // -tb, --threads-batch N
	CPUMask                 string   `json:"cpu_mask,omitempty"`        // -C, --cpu-mask M
	CPURange                string   `json:"cpu_range,omitempty"`       // -Cr, --cpu-range lo-hi
	CPUStrict               int      `json:"cpu_strict,omitempty"`      // --cpu-strict <0|1>
	Prio                    int      `json:"prio,omitempty"`            // --prio N
	Poll                    int      `json:"poll,omitempty"`            // --poll <0...100>
	CPUMaskBatch            string   `json:"cpu_mask_batch,omitempty"`  // -Cb, --cpu-mask-batch M
	CPURangeBatch           string   `json:"cpu_range_batch,omitempty"` // -Crb, --cpu-range-batch lo-hi
	CPUStrictBatch          int      `json:"cpu_strict_batch,omitempty"`
	PrioBatch               int      `json:"prio_batch,omitempty"`
	PollBatch               int      `json:"poll_batch,omitempty"`
	CtxSize                 int      `json:"ctx_size,omitempty"`     // -c, --ctx-size N
	Predict                 int      `json:"predict,omitempty"`      // -n, --predict, --n-predict N
	BatchSize               int      `json:"batch_size,omitempty"`   // -b, --batch-size N
	UBatchSize              int      `json:"ubatch_size,omitempty"`  // -ub, --ubatch-size N
	Keep                    int      `json:"keep,omitempty"`         // --keep N
	SWAFull                 bool     `json:"swa_full,omitempty"`     // --swa-full
	FlashAttn               string   `json:"flash_attn,omitempty"`   // -fa, --flash-attn [on|off|auto]
	Perf                    bool     `json:"perf,omitempty"`         // --perf
	NoPerf                  bool     `json:"no_perf,omitempty"`      // --no-perf
	Escape                  bool     `json:"escape,omitempty"`       // -e, --escape
	NoEscape                bool     `json:"no_escape,omitempty"`    // --no-escape
	RopeScaling             string   `json:"rope_scaling,omitempty"` // --rope-scaling {none,linear,yarn}
	RopeScale               float64  `json:"rope_scale,omitempty"`   // --rope-scale N
	RopeFreqBase            float64  `json:"rope_freq_base,omitempty"`
	RopeFreqScale           float64  `json:"rope_freq_scale,omitempty"`
	YarnOrigCtx             int      `json:"yarn_orig_ctx,omitempty"`
	YarnExtFactor           float64  `json:"yarn_ext_factor,omitempty"`
	YarnAttnFactor          float64  `json:"yarn_attn_factor,omitempty"`
	YarnBetaSlow            float64  `json:"yarn_beta_slow,omitempty"`
	YarnBetaFast            float64  `json:"yarn_beta_fast,omitempty"`
	KVOffload               bool     `json:"kv_offload,omitempty"`                 // -kvo, --kv-offload
	NoKVOffload             bool     `json:"no_kv_offload,omitempty"`              // -nkvo, --no-kv-offload
	Repack                  bool     `json:"repack,omitempty"`                     // --repack
	NoRepack                bool     `json:"no_repack,omitempty"`                  // -nr, --no-repack
	NoHost                  bool     `json:"no_host,omitempty"`                    // --no-host
	CacheTypeK              string   `json:"cache_type_k,omitempty"`               // -ctk, --cache-type-k TYPE
	CacheTypeV              string   `json:"cache_type_v,omitempty"`               // -ctv, --cache-type-v TYPE
	DefragThold             float64  `json:"defrag_thold,omitempty"`               // -dt, --defrag-thold N
	Mlock                   bool     `json:"mlock,omitempty"`                      // --mlock
	Mmap                    bool     `json:"mmap,omitempty"`                       // --mmap
	NoMmap                  bool     `json:"no_mmap,omitempty"`                    // --no-mmap
	DirectIO                bool     `json:"direct_io,omitempty"`                  // -dio, --direct-io
	NoDirectIO              bool     `json:"no_direct_io,omitempty"`               // -ndio, --no-direct-io
	Numa                    string   `json:"numa,omitempty"`                       // --numa TYPE
	Device                  string   `json:"device,omitempty"`                     // -dev, --device <dev1,dev2,..>
	OverrideTensor          []string `json:"override_tensor,omitempty"`            // -ot, --override-tensor
	CPUMoe                  bool     `json:"cpu_moe,omitempty"`                    // -cmoe, --cpu-moe
	NCPUMoe                 int      `json:"n_cpu_moe,omitempty"`                  // -ncmoe, --n-cpu-moe N
	GPULayers               int      `json:"gpu_layers,omitempty"`                 // -ngl, --gpu-layers, --n-gpu-layers N
	SplitMode               string   `json:"split_mode,omitempty"`                 // -sm, --split-mode {none,layer,row}
	TensorSplit             string   `json:"tensor_split,omitempty"`               // -ts, --tensor-split N0,N1,N2,...
	MainGPU                 int      `json:"main_gpu,omitempty"`                   // -mg, --main-gpu INDEX
	Fit                     string   `json:"fit,omitempty"`                        // -fit, --fit [on|off]
	FitTarget               string   `json:"fit_target,omitempty"`                 // -fitt, --fit-target MiB0,MiB1,MiB2,...
	FitCtx                  int      `json:"fit_ctx,omitempty"`                    // -fitc, --fit-ctx N
	CheckTensors            bool     `json:"check_tensors,omitempty"`              // --check-tensors
	OverrideKV              []string `json:"override_kv,omitempty"`                // --override-kv KEY=TYPE:VALUE,...
	OpOffload               bool     `json:"op_offload,omitempty"`                 // --op-offload
	NoOpOffload             bool     `json:"no_op_offload,omitempty"`              // --no-op-offload
	Lora                    []string `json:"lora,omitempty"`                       // --lora FNAME
	LoraScaled              []string `json:"lora_scaled,omitempty"`                // --lora-scaled FNAME:SCALE,...
	ControlVector           []string `json:"control_vector,omitempty"`             // --control-vector FNAME
	ControlVectorScaled     []string `json:"control_vector_scaled,omitempty"`      // --control-vector-scaled FNAME:SCALE,...
	ControlVectorLayerRange string   `json:"control_vector_layer_range,omitempty"` // --control-vector-layer-range START END
	Model                   string   `json:"model,omitempty"`                      // -m, --model FNAME
	ModelURL                string   `json:"model_url,omitempty"`                  // -mu, --model-url MODEL_URL
	DockerRepo              string   `json:"docker_repo,omitempty"`                // -dr, --docker-repo [<repo>/]<model>[:quant]
	HFRepo                  string   `json:"hf_repo,omitempty"`                    // -hf, -hfr, --hf-repo <user>/<model>[:quant]
	HFRepoDraft             string   `json:"hf_repo_draft,omitempty"`              // -hfd, -hfrd, --hf-repo-draft <user>/<model>[:quant]
	HFFile                  string   `json:"hf_file,omitempty"`                    // -hff, --hf-file FILE
	HFRepoV                 string   `json:"hf_repo_v,omitempty"`                  // -hfv, -hfrv, --hf-repo-v <user>/<model>[:quant]
	HFFileV                 string   `json:"hf_file_v,omitempty"`                  // -hffv, --hf-file-v FILE
	HFToken                 string   `json:"hf_token,omitempty"`                   // -hft, --hf-token TOKEN
	LogDisable              bool     `json:"log_disable,omitempty"`                // --log-disable
	LogFile                 string   `json:"log_file,omitempty"`                   // --log-file FNAME
	LogColors               string   `json:"log_colors,omitempty"`                 // --log-colors [on|off|auto]
	Verbose                 bool     `json:"verbose,omitempty"`                    // -v, --verbose, --log-verbose
	Offline                 bool     `json:"offline,omitempty"`                    // --offline
	Verbosity               int      `json:"verbosity,omitempty"`                  // -lv, --verbosity, --log-verbosity N
	LogPrefix               bool     `json:"log_prefix,omitempty"`                 // --log-prefix
	LogTimestamps           bool     `json:"log_timestamps,omitempty"`             // --log-timestamps
	CacheTypeKDraft         string   `json:"cache_type_k_draft,omitempty"`         // -ctkd, --cache-type-k-draft TYPE
	CacheTypeVDraft         string   `json:"cache_type_v_draft,omitempty"`         // -ctvd, --cache-type-v-draft TYPE

	// Sampling params
	Samplers           string   `json:"samplers,omitempty"`          // --samplers SAMPLERS
	Seed               int      `json:"seed,omitempty"`              // -s, --seed SEED
	SamplingSeq        string   `json:"sampling_seq,omitempty"`      // --sampler-seq, --sampling-seq SEQUENCE
	IgnoreEOS          bool     `json:"ignore_eos,omitempty"`        // --ignore-eos
	Temperature        float64  `json:"temp,omitempty"`              // --temp N
	TopK               int      `json:"top_k,omitempty"`             // --top-k N
	TopP               float64  `json:"top_p,omitempty"`             // --top-p N
	MinP               float64  `json:"min_p,omitempty"`             // --min-p N
	TopNSigma          float64  `json:"top_nsigma,omitempty"`        // --top-nsigma N
	XTCProbability     float64  `json:"xtc_probability,omitempty"`   // --xtc-probability N
	XTCThreshold       float64  `json:"xtc_threshold,omitempty"`     // --xtc-threshold N
	Typical            float64  `json:"typical,omitempty"`           // --typical N
	RepeatLastN        int      `json:"repeat_last_n,omitempty"`     // --repeat-last-n N
	RepeatPenalty      float64  `json:"repeat_penalty,omitempty"`    // --repeat-penalty N
	PresencePenalty    float64  `json:"presence_penalty,omitempty"`  // --presence-penalty N
	FrequencyPenalty   float64  `json:"frequency_penalty,omitempty"` // --frequency-penalty N
	DryMultiplier      float64  `json:"dry_multiplier,omitempty"`    // --dry-multiplier N
	DryBase            float64  `json:"dry_base,omitempty"`          // --dry-base N
	DryAllowedLength   int      `json:"dry_allowed_length,omitempty"`
	DryPenaltyLastN    int      `json:"dry_penalty_last_n,omitempty"`
	DrySequenceBreaker []string `json:"dry_sequence_breaker,omitempty"` // --dry-sequence-breaker STRING
	AdaptiveTarget     float64  `json:"adaptive_target,omitempty"`      // --adaptive-target N
	AdaptiveDecay      float64  `json:"adaptive_decay,omitempty"`       // --adaptive-decay N
	DynatempRange      float64  `json:"dynatemp_range,omitempty"`       // --dynatemp-range N
	DynatempExp        float64  `json:"dynatemp_exp,omitempty"`         // --dynatemp-exp N
	Mirostat           int      `json:"mirostat,omitempty"`             // --mirostat N
	MirostatLR         float64  `json:"mirostat_lr,omitempty"`          // --mirostat-lr N
	MirostatEnt        float64  `json:"mirostat_ent,omitempty"`         // --mirostat-ent N
	LogitBias          []string `json:"logit_bias,omitempty"`           // -l, --logit-bias TOKEN_ID(+/-)BIAS
	Grammar            string   `json:"grammar,omitempty"`              // --grammar GRAMMAR
	GrammarFile        string   `json:"grammar_file,omitempty"`         // --grammar-file FNAME
	JSONSchema         string   `json:"json_schema,omitempty"`          // -j, --json-schema SCHEMA
	JSONSchemaFile     string   `json:"json_schema_file,omitempty"`     // -jf, --json-schema-file FILE
	BackendSampling    bool     `json:"backend_sampling,omitempty"`     // -bs, --backend-sampling

	// Server-specific params
	CtxCheckpoints       int      `json:"ctx_checkpoints,omitempty"`        // --ctx-checkpoints, --swa-checkpoints N
	CacheRAM             int      `json:"cache_ram,omitempty"`              // -cram, --cache-ram N
	LookupCacheStatic    bool     `json:"lookup_cache_static,omitempty"`    // -lcs, --lookup-cache-static
	LookupCacheDynamic   bool     `json:"lookup_cache_dynamic,omitempty"`   // -lcd, --lookup-cache-dynamic
