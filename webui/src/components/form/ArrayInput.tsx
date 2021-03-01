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
