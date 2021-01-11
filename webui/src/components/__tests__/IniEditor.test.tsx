import { describe, it, expect, vi } from 'vitest'
import { useState } from 'react'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import IniEditor from '@/components/instance/IniEditor'
import { getLlamaFieldSuggestions } from '@/lib/llamaFieldSuggestions'

const [first, second] = getLlamaFieldSuggestions('m')

const SELECTED = 'bg-accent text-accent-foreground'

function rowOf(name: string): HTMLElement {
  return screen.getByText(name).closest('div') as HTMLElement
}

function renderEditor() {
  const onChange = vi.fn()
  const Harness = () => {
    const [value, setValue] = useState('')
    return <IniEditor value={value} onChange={(v) => { onChange(v); setValue(v) }} />
  }
