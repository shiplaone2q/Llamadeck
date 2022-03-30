import { z } from 'zod'

// Define the vLLM backend options schema
export const VllmBackendOptionsSchema = z.object({
  // Basic connection options (auto-assigned by llamadeck)
  host: z.string().optional(),
  port: z.number().optional(),

  // Model and engine configuration
  model: z.string().optional(),
  tokenizer: z.string().optional(),
  skip_tokenizer_init: z.boolean().optional(),
  revision: z.string().optional(),
  code_revision: z.string().optional(),
  tokenizer_revision: z.string().optional(),
  tokenizer_mode: z.string().optional(),
  trust_remote_code: z.boolean().optional(),
  download_dir: z.string().optional(),
  load_format: z.string().optional(),
  config_format: z.string().optional(),
  dtype: z.string().optional(),
  kv_cache_dtype: z.string().optional(),
  quantization_param_path: z.string().optional(),
  seed: z.number().optional(),
  max_model_len: z.number().optional(),
  guided_decoding_backend: z.string().optional(),
  distributed_executor_backend: z.string().optional(),
  worker_use_ray: z.boolean().optional(),
  ray_workers_use_nsight: z.boolean().optional(),

  // Performance and serving configuration
  block_size: z.number().optional(),
  enable_prefix_caching: z.boolean().optional(),
  disable_sliding_window: z.boolean().optional(),
  use_v2_block_manager: z.boolean().optional(),
  num_lookahead_slots: z.number().optional(),
  swap_space: z.number().optional(),
  cpu_offload_gb: z.number().optional(),
  gpu_memory_utilization: z.number().optional(),
  num_gpu_blocks_override: z.number().optional(),
  max_num_batched_tokens: z.number().optional(),
  max_num_seqs: z.number().optional(),
  max_logprobs: z.number().optional(),
  disable_log_stats: z.boolean().optional(),
  quantization: z.string().optional(),
