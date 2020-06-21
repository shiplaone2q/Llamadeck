import type React from 'react'
import { useState } from 'react'
import type { CreateInstanceOptions } from '@/types/instance'
import { BackendType } from '@/types/instance'
import { getAdvancedBackendFields } from '@/lib/zodFormUtils'
import BackendFormField from '@/components/BackendFormField'
import { ChevronDown, ChevronRight } from 'lucide-react'

interface AdvancedTabProps {
  formData: CreateInstanceOptions
  onBackendFieldChange: (key: string, value: unknown) => void
}

const AdvancedTab: React.FC<AdvancedTabProps> = ({
  formData,
  onBackendFieldChange
}) => {
  const [showAdvanced, setShowAdvanced] = useState(false)
  const advancedBackendFields = getAdvancedBackendFields(formData.backend_type)

  const extraArgs = (formData.backend_options as Record<string, unknown>)?.extra_args as Record<string, string> | undefined
  const advancedOptionsCount = advancedBackendFields.filter(fieldKey => {
    const value = (formData.backend_options as Record<string, unknown>)?.[fieldKey]
    return value !== undefined && value !== null && value !== ''
  }).length
