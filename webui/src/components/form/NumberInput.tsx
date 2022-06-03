import type React from 'react'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

interface NumberInputProps {
  id: string
  label: string
  value: number | undefined
  onChange: (value: number | undefined) => void
  placeholder?: string
  description?: string
  disabled?: boolean
  className?: string
}

const NumberInput: React.FC<NumberInputProps> = ({
  id,
  label,
  value,
  onChange,
  placeholder,
  description,
  disabled = false,
  className
}) => {
  const handleChange = (inputValue: string) => {
    if (inputValue === '') {
      onChange(undefined)
      return
    }

    const numValue = parseFloat(inputValue)
    if (!Number.isNaN(numValue)) {
      onChange(numValue)
    }
  }

  return (
