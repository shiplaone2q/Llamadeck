import type { CreateInstanceOptions } from '@/schemas/instanceOptions'

export type { CreateInstanceOptions } from '@/schemas/instanceOptions'

export const BackendType = {
  LLAMA_CPP: 'llama_cpp',
  MLX_LM: 'mlx_lm',
  VLLM: 'vllm',
  CUSTOM: 'custom',
  // MLX_VLM: 'mlx_vlm',  // Future expansion
} as const

export type BackendTypeValue = typeof BackendType[keyof typeof BackendType]
