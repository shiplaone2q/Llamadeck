import React, { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { llamaCppApi } from '@/lib/api'
import { RefreshCw, Loader2, AlertCircle } from 'lucide-react'

interface ModelsDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  instanceName: string
  isRunning: boolean
}

interface Model {
  id: string
  object: string
  owned_by: string
  created: number
  in_cache: boolean
  path: string
  status?: {
    value: string // "loaded" | "loading" | "unloaded"
    args: string[]
  }
}

const StatusIcon: React.FC<{ status: string }> = ({ status }) => {
  switch (status) {
    case 'loaded':
      return (
        <div className="h-2 w-2 rounded-full bg-green-500" />
      )
    case 'loading':
      return (
        <Loader2
          className="h-3 w-3 animate-spin text-yellow-500"
        />
      )
    case 'unloaded':
      return (
        <div className="h-2 w-2 rounded-full bg-gray-400" />
      )
    default:
      return null
  }
}

const ModelsDialog: React.FC<ModelsDialogProps> = ({
  open,
  onOpenChange,
  instanceName,
  isRunning,
}) => {
  const [models, setModels] = useState<Model[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [loadingModels, setLoadingModels] = useState<Set<string>>(new Set())
