import { z } from 'zod'

// Define the LlamaCpp backend options schema
export const LlamaCppBackendOptionsSchema = z.object({
  // Common params (ordered as in llama-cpp.md)
  verbose_prompt: z.boolean().optional(),
  threads: z.number().optional(), // -t, --threads N
  threads_batch: z.number().optional(), // -tb, --threads-batch N
  cpu_mask: z.string().optional(), // -C, --cpu-mask M
  cpu_range: z.string().optional(), // -Cr, --cpu-range lo-hi
  cpu_strict: z.number().optional(), // --cpu-strict <0|1>
  prio: z.number().optional(), // --prio N
  poll: z.number().optional(), // --poll <0...100>
  cpu_mask_batch: z.string().optional(), // -Cb, --cpu-mask-batch M
  cpu_range_batch: z.string().optional(), // -Crb, --cpu-range-batch lo-hi
  cpu_strict_batch: z.number().optional(),
  prio_batch: z.number().optional(),
  poll_batch: z.number().optional(),
  ctx_size: z.number().optional(), // -c, --ctx-size N
  predict: z.number().optional(), // -n, --predict, --n-predict N
  batch_size: z.number().optional(), // -b, --batch-size N
  ubatch_size: z.number().optional(), // -ub, --ubatch-size N
  keep: z.number().optional(), // --keep N
  swa_full: z.boolean().optional(), // --swa-full
  flash_attn: z.string().optional(), // -fa, --flash-attn [on|off|auto]
  perf: z.boolean().optional(), // --perf
  no_perf: z.boolean().optional(), // --no-perf
  escape: z.boolean().optional(), // -e, --escape
  no_escape: z.boolean().optional(), // --no-escape
  rope_scaling: z.string().optional(), // --rope-scaling {none,linear,yarn}
  rope_scale: z.number().optional(), // --rope-scale N
  rope_freq_base: z.number().optional(),
  rope_freq_scale: z.number().optional(),
  yarn_orig_ctx: z.number().optional(),
  yarn_ext_factor: z.number().optional(),
  yarn_attn_factor: z.number().optional(),
  yarn_beta_slow: z.number().optional(),
  yarn_beta_fast: z.number().optional(),
  kv_offload: z.boolean().optional(), // -kvo, --kv-offload
  no_kv_offload: z.boolean().optional(), // -nkvo, --no-kv-offload
  repack: z.boolean().optional(), // --repack
  no_repack: z.boolean().optional(), // -nr, --no-repack
  no_host: z.boolean().optional(), // --no-host
  cache_type_k: z.string().optional(), // -ctk, --cache-type-k TYPE
  cache_type_v: z.string().optional(), // -ctv, --cache-type-v TYPE
  defrag_thold: z.number().optional(), // -dt, --defrag-thold N
  mlock: z.boolean().optional(), // --mlock
  mmap: z.boolean().optional(), // --mmap
  no_mmap: z.boolean().optional(), // --no-mmap
  direct_io: z.boolean().optional(), // -dio, --direct-io
  no_direct_io: z.boolean().optional(), // -ndio, --no-direct-io
  numa: z.string().optional(), // --numa TYPE
  device: z.string().optional(), // -dev, --device <dev1,dev2,..>
  override_tensor: z.array(z.string()).optional(), // -ot, --override-tensor
  cpu_moe: z.boolean().optional(), // -cmoe, --cpu-moe
  n_cpu_moe: z.number().optional(), // -ncmoe, --n-cpu-moe N
  gpu_layers: z.number().optional(), // -ngl, --gpu-layers, --n-gpu-layers N
  split_mode: z.string().optional(), // -sm, --split-mode {none,layer,row}
  tensor_split: z.string().optional(), // -ts, --tensor-split N0,N1,N2,...
  main_gpu: z.number().optional(), // -mg, --main-gpu INDEX
  fit: z.string().optional(), // -fit, --fit [on|off]
  fit_target: z.string().optional(), // -fitt, --fit-target MiB0,MiB1,MiB2,...
  fit_ctx: z.number().optional(), // -fitc, --fit-ctx N
  check_tensors: z.boolean().optional(), // --check-tensors
  override_kv: z.array(z.string()).optional(), // --override-kv KEY=TYPE:VALUE,...
  op_offload: z.boolean().optional(), // --op-offload
  no_op_offload: z.boolean().optional(), // --no-op-offload
  lora: z.array(z.string()).optional(), // --lora FNAME
  lora_scaled: z.array(z.string()).optional(), // --lora-scaled FNAME:SCALE,...
  control_vector: z.array(z.string()).optional(), // --control-vector FNAME
  control_vector_scaled: z.array(z.string()).optional(), // --control-vector-scaled FNAME:SCALE,...
  control_vector_layer_range: z.string().optional(), // --control-vector-layer-range START END
  model: z.string().optional(), // -m, --model FNAME
  model_url: z.string().optional(), // -mu, --model-url MODEL_URL
  docker_repo: z.string().optional(), // -dr, --docker-repo [<repo>/]<model>[:quant]
  hf_repo: z.string().optional(), // -hf, -hfr, --hf-repo <user>/<model>[:quant]
  hf_repo_draft: z.string().optional(), // -hfd, -hfrd, --hf-repo-draft <user>/<model>[:quant]
  hf_file: z.string().optional(), // -hff, --hf-file FILE
  hf_repo_v: z.string().optional(), // -hfv, -hfrv, --hf-repo-v <user>/<model>[:quant]
  hf_file_v: z.string().optional(), // -hffv, --hf-file-v FILE
  hf_token: z.string().optional(), // -hft, --hf-token TOKEN
  log_disable: z.boolean().optional(), // --log-disable
  log_file: z.string().optional(), // --log-file FNAME
  log_colors: z.string().optional(), // --log-colors [on|off|auto]
  verbose: z.boolean().optional(), // -v, --verbose, --log-verbose
  offline: z.boolean().optional(), // --offline
  verbosity: z.number().optional(), // -lv, --verbosity, --log-verbosity N
  log_prefix: z.boolean().optional(), // --log-prefix
  log_timestamps: z.boolean().optional(), // --log-timestamps
  cache_type_k_draft: z.string().optional(), // -ctkd, --cache-type-k-draft TYPE
  cache_type_v_draft: z.string().optional(), // -ctvd, --cache-type-v-draft TYPE

  // Sampling params (ordered as in llama-cpp.md)
  samplers: z.string().optional(), // --samplers SAMPLERS
  seed: z.number().optional(), // -s, --seed SEED
