import type React from 'react'
import { useState } from 'react'
import { BackendType, type CreateInstanceOptions } from '@/types/instance'
import { Button } from '@/components/ui/button'
import { Terminal, ChevronDown, ChevronRight } from 'lucide-react'
import { getBasicBackendFields } from '@/lib/zodFormUtils'
import BackendFormField from '@/components/BackendFormField'
import SelectInput from '@/components/form/SelectInput'
import CheckboxInput from '@/components/form/CheckboxInput'
import TextInput from '@/components/form/TextInput'
import EnvVarsInput from '@/components/form/EnvVarsInput'
import { useBackendSettings, useCustomBackends, useCustomBackendSettings } from '@/hooks/useConfig'
import PresetDialog from './PresetDialog'

const CUSTOM_OPTION_PREFIX = 'custom:'

interface BackendTabProps {
  formData: CreateInstanceOptions
  onBackendFieldChange: (key: string, value: unknown) => void
  onChange: (key: keyof CreateInstanceOptions, value: unknown) => void
  onParseCommand: () => void
}

const BackendTab: React.FC<BackendTabProps> = ({
  formData,
  onBackendFieldChange,
  onChange,
  onParseCommand
}) => {
  const [showExecutionContext, setShowExecutionContext] = useState(false)
  const [showPresetDialog, setShowPresetDialog] = useState(false)
  const backendSettings = useBackendSettings(formData.backend_type)
  const basicBackendFields = getBasicBackendFields(formData.backend_type)

  const isCustom = formData.backend_type === BackendType.CUSTOM
  const customName = (formData.backend_options as Record<string, unknown> | undefined)?.name as string | undefined
  const customBackends = useCustomBackends()
  const customSettings = useCustomBackendSettings(customName)

  const handleBackendTypeChange = (value: string | undefined) => {
    if (value?.startsWith(CUSTOM_OPTION_PREFIX)) {
      const name = value.slice(CUSTOM_OPTION_PREFIX.length)
      onChange('backend_type', BackendType.CUSTOM)
      onBackendFieldChange('name', name)
      onChange('docker_enabled', customBackends[name]?.docker?.enabled ?? false)
      return
    }
    onChange('backend_type', value)
  }

  const getCommandPlaceholder = () => {
    if (isCustom) {
      return customSettings?.command || 'configured in backends.custom.<name>'
    }

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
