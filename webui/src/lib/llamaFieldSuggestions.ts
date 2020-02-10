import {
  getAllLlamaCppFieldKeys,
  getLlamaCppFieldType,
  getAllLlamaCppAltKeys,
  getLlamaCppAltKeyType
} from '@/schemas/instanceOptions'

export interface FieldSuggestion {
  name: string
  type: 'text' | 'number' | 'boolean' | 'array'
