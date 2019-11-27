import type React from 'react'
import KeyValueInput from './KeyValueInput'

interface EnvVarsInputProps {
  id: string
  label: string
  value: Record<string, string> | undefined
