import {
  getAllLlamaCppFieldKeys,
  getLlamaCppFieldType,
  getAllLlamaCppAltKeys,
  getLlamaCppAltKeyType
} from '@/schemas/instanceOptions'

export interface FieldSuggestion {
  name: string
  type: 'text' | 'number' | 'boolean' | 'array'
}

function snakeToKebab(snake: string): string {
  return snake.replace(/_/g, '-')
}

export function getLlamaFieldSuggestions(input: string): FieldSuggestion[] {
  const allFields = getAllLlamaCppFieldKeys()
  const allAltKeys = getAllLlamaCppAltKeys()
