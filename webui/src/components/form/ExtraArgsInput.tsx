import type React from 'react'
import KeyValueInput from './KeyValueInput'

interface ExtraArgsInputProps {
  id: string
  label: string
  value: Record<string, string> | undefined
  onChange: (value: Record<string, string> | undefined) => void
  description?: string
  disabled?: boolean
  className?: string
}

const ExtraArgsInput: React.FC<ExtraArgsInputProps> = (props) => {
