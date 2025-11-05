import type React from 'react'
import { useState } from 'react'
import { BackendType, type CreateInstanceOptions } from '@/types/instance'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Terminal, ChevronDown, ChevronRight } from 'lucide-react'
import { getBasicBackendFields, getAdvancedBackendFields } from '@/lib/zodFormUtils'
import BackendFormField from '@/components/BackendFormField'
import SelectInput from '@/components/form/SelectInput'
import ExecutionContextSection from '@/components/instance/ExecutionContextSection'

interface BackendConfigurationCardProps {
  formData: CreateInstanceOptions
  onBackendFieldChange: (key: string, value: unknown) => void
  onChange: (key: keyof CreateInstanceOptions, value: unknown) => void
  onParseCommand: () => void
}

const BackendConfigurationCard: React.FC<BackendConfigurationCardProps> = ({
  formData,
  onBackendFieldChange,
  onChange,
  onParseCommand
}) => {
  const [showAdvanced, setShowAdvanced] = useState(false)
  const basicBackendFields = getBasicBackendFields(formData.backend_type)
  const advancedBackendFields = getAdvancedBackendFields(formData.backend_type)

  return (
    <Card>
      <CardHeader>
        <CardTitle>Backend Configuration</CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Backend Type Selection */}
        <SelectInput
          id="backend_type"
          label="Backend Type"
          value={formData.backend_type || BackendType.LLAMA_CPP}
          onChange={(value) => onChange('backend_type', value)}
          options={[
            { value: BackendType.LLAMA_CPP, label: 'Llama Server' },
            { value: BackendType.MLX_LM, label: 'MLX LM' },
            { value: BackendType.VLLM, label: 'vLLM' }
          ]}
          description="Select the backend server type"
        />

        {/* Parse Command Section */}
        <div className="flex flex-col gap-2">
          <Button
            variant="outline"
            onClick={onParseCommand}
            className="flex items-center gap-2 w-fit"
          >
            <Terminal className="h-4 w-4" />
