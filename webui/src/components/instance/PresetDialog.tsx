import type React from 'react'
import { useRef } from 'react'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Upload, Download } from 'lucide-react'
import IniEditor from './IniEditor'
import type { CreateInstanceOptions } from '@/schemas/instanceOptions'

interface PresetDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  formData: CreateInstanceOptions
  onChange: (key: keyof CreateInstanceOptions, value: unknown) => void
}

const PresetDialog: React.FC<PresetDialogProps> = ({
  open,
  onOpenChange,
  formData,
  onChange
}) => {
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleImport = () => {
    fileInputRef.current?.click()
  }

