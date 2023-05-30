import type React from 'react'
import { Label } from '@/components/ui/label'

interface SelectOption {
  value: string
  label: string
}

interface SelectInputProps {
  id: string
  label: string
  value: string | undefined
  onChange: (value: string | undefined) => void
  options: SelectOption[]
  description?: string
  disabled?: boolean
  className?: string
}

const SelectInput: React.FC<SelectInputProps> = ({
  id,
  label,
  value,
  onChange,
  options,
  description,
  disabled = false,
  className
}) => {
  return (
    <div className="grid gap-2">
      <Label htmlFor={id}>
        {label}
      </Label>
      <select
        id={id}
        value={value || ''}
        onChange={(e) => onChange(e.target.value || undefined)}
