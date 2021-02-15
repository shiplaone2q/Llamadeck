import type React from 'react'
import { useState, useEffect } from 'react'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { X, Plus } from 'lucide-react'

interface KeyValueInputProps {
  id: string
  label: string
  value: Record<string, string> | undefined
  onChange: (value: Record<string, string> | undefined) => void
  description?: string
  disabled?: boolean
  className?: string
  keyPlaceholder?: string
  valuePlaceholder?: string
  addButtonText?: string
  helperText?: string
  allowEmptyValues?: boolean // If true, entries with empty values are considered valid
}

interface KeyValuePair {
  key: string
  value: string
}

const KeyValueInput: React.FC<KeyValueInputProps> = ({
  id,
  label,
  value,
  onChange,
  description,
  disabled = false,
  className,
  keyPlaceholder = 'Key',
  valuePlaceholder = 'Value',
  addButtonText = 'Add Entry',
  helperText,
  allowEmptyValues = false
}) => {
  // Convert the value object to an array of key-value pairs for editing
  const pairsFromValue = value
    ? Object.entries(value).map(([key, val]) => ({ key, value: val }))
    : []

  const [pairs, setPairs] = useState<KeyValuePair[]>(
    pairsFromValue.length > 0 ? pairsFromValue : [{ key: '', value: '' }]
  )

  // Sync internal state when value prop changes
  useEffect(() => {
    const newPairsFromValue = value
      ? Object.entries(value).map(([key, val]) => ({ key, value: val }))
      : []

