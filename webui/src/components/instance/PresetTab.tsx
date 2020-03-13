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
