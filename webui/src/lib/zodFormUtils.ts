import {
  type LlamaCppBackendOptions,
  type MlxBackendOptions,
  type VllmBackendOptions,
  type CustomBackendOptions,
  LlamaCppBackendOptionsSchema,
  MlxBackendOptionsSchema,
  VllmBackendOptionsSchema,
  CustomBackendOptionsSchema,
  getAllLlamaCppFieldKeys,
  getAllMlxFieldKeys,
  getAllVllmFieldKeys,
  getAllCustomFieldKeys,
  getLlamaCppFieldType,
  getMlxFieldType,
  getVllmFieldType,
  getCustomFieldType
} from '@/schemas/instanceOptions'

// LlamaCpp backend-specific basic fields
const basicLlamaCppFieldsConfig: Record<string, {
  label: string
  description?: string
  placeholder?: string
}> = {
  model: {
    label: 'Model Path',
    placeholder: '/path/to/model.gguf',
    description: 'Path to the model file'
  },
  hf_repo: {
    label: 'Hugging Face Repository',
    placeholder: 'microsoft/DialoGPT-medium',
    description: 'Hugging Face model repository'
  },
  gpu_layers: {
    label: 'GPU Layers',
    placeholder: '0',
    description: 'Number of layers to offload to GPU'
  },
  models_preset: {
    label: 'Models Preset Path',
    placeholder: '/path/to/preset.ini',
    description: 'Optional: Path to preset.ini for router mode. Leave empty for normal operation, or use Preset tab to auto-generate'
  }
}

// MLX backend-specific basic fields
const basicMlxFieldsConfig: Record<string, {
  label: string
