import type React from 'react'
import { useState, useEffect } from 'react'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { FileCode } from 'lucide-react'
import { getBackendFieldType, basicBackendFieldsConfig } from '@/lib/zodFormUtils'
import { parseArgs } from '@/lib/parseArgs'
import ExtraArgsInput from '@/components/form/ExtraArgsInput'
import type { CreateInstanceOptions } from '@/schemas/instanceOptions'

interface BackendFormFieldProps {
  fieldKey: string
  value: string | number | boolean | string[] | Record<string, string> | undefined
  onChange: (key: string, value: string | number | boolean | string[] | Record<string, string> | undefined) => void
  formData?: CreateInstanceOptions
  onOpenPresetDialog?: () => void
}

const BackendFormField: React.FC<BackendFormFieldProps> = ({ fieldKey, value, onChange, formData, onOpenPresetDialog }) => {
  // Local state for array fields to handle spaces correctly
  const [localArrayValue, setLocalArrayValue] = useState<string>('')
  const [isEditingArray, setIsEditingArray] = useState(false)

  // Sync local value when not editing
  useEffect(() => {
    if (!isEditingArray && Array.isArray(value)) {
      setLocalArrayValue(value.join(' '))
    }
  }, [value, isEditingArray])
  // Special handling for models_preset field
  if (fieldKey === 'models_preset') {
    const hasPresetContent = formData?.preset_ini && formData.preset_ini.trim().length > 0
    const isCustomPath = value && value.toString().trim().length > 0

    let badge = null
    let callout = null

    if (!isCustomPath && hasPresetContent) {
      badge = <Badge variant="secondary" className="ml-2">Auto</Badge>
      callout = (
        <Alert className="border-primary/40 bg-primary/10 text-primary">
          <FileCode />
          <AlertTitle>Preset defined</AlertTitle>
          <AlertDescription className="text-muted-foreground">
            Will be auto-set to the preset.ini created in Preset Editor
          </AlertDescription>
        </Alert>
      )
    } else if (isCustomPath) {
      badge = <Badge variant="outline" className="ml-2">Custom</Badge>
    }

    return (
      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center">
            <Label htmlFor={fieldKey}>Models Preset Path</Label>
            {badge}
          </div>
          {onOpenPresetDialog && (
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={onOpenPresetDialog}
              className="flex items-center gap-2"
            >
              <FileCode className="h-4 w-4" />
              Edit Preset
            </Button>
          )}
        </div>
        <Input
          id={fieldKey}
          value={value?.toString() || ''}
          onChange={(e) => onChange(fieldKey, e.target.value)}
