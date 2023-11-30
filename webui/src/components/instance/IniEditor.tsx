import type React from 'react'
import { useState, useRef, useEffect, type KeyboardEvent } from 'react'
import { Textarea } from '@/components/ui/textarea'
import { getLlamaFieldSuggestions, type FieldSuggestion } from '@/lib/llamaFieldSuggestions'

interface IniEditorProps {
  value: string
  onChange: (value: string) => void
  className?: string
}

const IniEditor: React.FC<IniEditorProps> = ({ value, onChange, className }) => {
  const textareaRef = useRef<HTMLTextAreaElement>(null)
  // keyup listener must not re-run updateSuggestions for keys the keydown handler already
  // consumed (arrow navigation, Enter/Tab apply, Escape dismiss) — it would reset selectedIndex
  // and resurrect a dismissed popup
  const consumedKeyRef = useRef<string | null>(null)
  const [cursorPosition, setCursorPosition] = useState(0)
  const [suggestions, setSuggestions] = useState<FieldSuggestion[]>([])
  const [selectedIndex, setSelectedIndex] = useState(0)
  const [showSuggestions, setShowSuggestions] = useState(false)
  const [suggestionPosition, setSuggestionPosition] = useState({ top: 0, left: 0 })

  const updateSuggestions = (text: string, cursorPos: number) => {
    const textBeforeCursor = text.substring(0, cursorPos)
    const currentLine = textBeforeCursor.split('\n').pop() || ''

    const trimmedLine = currentLine.trim()

    if (trimmedLine.startsWith('[') || trimmedLine.startsWith(';') || trimmedLine === '') {
      setShowSuggestions(false)
      return
    }

    const equalIndex = trimmedLine.indexOf('=')
    if (equalIndex !== -1) {
      setShowSuggestions(false)
      return
    }

    const fieldValue = trimmedLine
    if (fieldValue.length > 0) {
      const allSuggestions = getLlamaFieldSuggestions(fieldValue)
      setSuggestions(allSuggestions)
      setSelectedIndex(0)
