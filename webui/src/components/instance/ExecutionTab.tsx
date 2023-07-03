import type React from 'react'
import { BackendType, type CreateInstanceOptions } from '@/types/instance'
import CheckboxInput from '@/components/form/CheckboxInput'
import TextInput from '@/components/form/TextInput'
import EnvVarsInput from '@/components/form/EnvVarsInput'
import SelectInput from '@/components/form/SelectInput'
import { useBackendSettings } from '@/hooks/useConfig'

interface ExecutionTabProps {
  formData: CreateInstanceOptions
  onChange: (key: keyof CreateInstanceOptions, value: unknown) => void
}

const ExecutionTab: React.FC<ExecutionTabProps> = ({
  formData,
  onChange
}) => {
  const backendSettings = useBackendSettings(formData.backend_type)

  const getCommandPlaceholder = () => {
    if (backendSettings?.command) {
      return backendSettings.command
    }

    switch (formData.backend_type) {
      case BackendType.LLAMA_CPP:
        return "llama-server"
      case BackendType.VLLM:
        return "vllm"
      case BackendType.MLX_LM:
        return "mlx_lm.server"
      default:
        return ""
    }
  }

  return (
    <div className="space-y-6 py-4">
      <SelectInput
        id="backend_type"
