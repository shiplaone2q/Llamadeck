import type React from 'react'
import type { CreateInstanceOptions } from '@/types/instance'
import { getBasicBackendFields, getAdvancedBackendFields } from '@/lib/zodFormUtils'
import BackendFormField from '@/components/BackendFormField'

type BackendFieldValue = string | number | boolean | string[] | Record<string, string> | undefined

interface BackendConfigurationProps {
  formData: CreateInstanceOptions
  onBackendFieldChange: (key: string, value: BackendFieldValue) => void
  showAdvanced?: boolean
}

const BackendConfiguration: React.FC<BackendConfigurationProps> = ({
  formData,
  onBackendFieldChange,
  showAdvanced = false
}) => {
  const basicBackendFields = getBasicBackendFields(formData.backend_type)
  const advancedBackendFields = getAdvancedBackendFields(formData.backend_type)

  return (
    <div className="space-y-4">
      <h3 className="text-lg font-medium">Backend Configuration</h3>
