import type React from 'react'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

interface TextInputProps {
  id: string
  label: string
  value: string | number | undefined
  onChange: (value: string | undefined) => void
  placeholder?: string
  description?: string
  disabled?: boolean
  className?: string
}

const TextInput: React.FC<TextInputProps> = ({
  id,
  label,
  value,
