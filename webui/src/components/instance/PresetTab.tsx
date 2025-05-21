import type React from 'react'
import { useRef } from 'react'
import { Button } from '@/components/ui/button'
import { Upload, Download } from 'lucide-react'
import IniEditor from './IniEditor'
import type { CreateInstanceOptions } from '@/schemas/instanceOptions'

interface PresetTabProps {
  formData: CreateInstanceOptions
  onChange: (key: keyof CreateInstanceOptions, value: unknown) => void
}

const PresetTab: React.FC<PresetTabProps> = ({ formData, onChange }) => {
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleImport = () => {
    fileInputRef.current?.click()
  }

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (event) => {
        const content = event.target?.result as string
        onChange('preset_ini', content)
      }
      reader.readAsText(file)
    }
  }

  const handleExport = () => {
    const content = formData.preset_ini || ''
    const blob = new Blob([content], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'preset.ini'
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }

  return (
    <div className="h-full flex flex-col space-y-4">
      <div>
        <h3 className="text-lg font-semibold">Model Presets (preset.ini)</h3>
