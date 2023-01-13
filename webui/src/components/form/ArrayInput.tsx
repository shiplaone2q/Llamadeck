import type React from 'react'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

interface ArrayInputProps {
  id: string
  label: string
  value: string[] | undefined
  onChange: (value: string[] | undefined) => void
  placeholder?: string
  description?: string
  disabled?: boolean
  className?: string
}

const ArrayInput: React.FC<ArrayInputProps> = ({
  id,
  label,
  value,
  onChange,
  placeholder = "item1, item2, item3",
  description,
  disabled = false,
  className
}) => {
  const handleChange = (inputValue: string) => {
    if (inputValue === '') {
      onChange(undefined)
      return
    }

    const arrayValue = inputValue
      .split(',')
      .map(s => s.trim())
      .filter(Boolean)

    onChange(arrayValue.length > 0 ? arrayValue : undefined)
  }

  return (
    <div className="grid gap-2">
      <Label htmlFor={id}>
        {label}
      </Label>
