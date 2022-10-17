import type React from 'react'
import { Checkbox } from '@/components/ui/checkbox'
import { Label } from '@/components/ui/label'

interface CheckboxInputProps {
  id: string
  label: string
  value: boolean | undefined
  onChange: (value: boolean) => void
  description?: string
  disabled?: boolean
  className?: string
}

const CheckboxInput: React.FC<CheckboxInputProps> = ({
  id,
  label,
  value,
  onChange,
  description,
  disabled = false,
  className
