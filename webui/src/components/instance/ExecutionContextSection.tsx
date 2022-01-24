import type React from 'react'
import { BackendType, type CreateInstanceOptions } from '@/types/instance'
import CheckboxInput from '@/components/form/CheckboxInput'
import TextInput from '@/components/form/TextInput'
import EnvVarsInput from '@/components/form/EnvVarsInput'
import { useBackendSettings } from '@/hooks/useConfig'

interface ExecutionContextSectionProps {
  formData: CreateInstanceOptions
  onChange: (key: keyof CreateInstanceOptions, value: unknown) => void
}

const ExecutionContextSection: React.FC<ExecutionContextSectionProps> = ({
  formData,
  onChange
}) => {
  const backendSettings = useBackendSettings(formData.backend_type)

