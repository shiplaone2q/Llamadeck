import React, { useState, useEffect, useRef } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { instancesApi } from '@/lib/api'
import { 
  RefreshCw, 
  Download, 
  Copy, 
  CheckCircle, 
  AlertCircle,
  Loader2,
  Settings
} from 'lucide-react'

interface LogsDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  instanceName: string
  isRunning: boolean
}

const LogsDialog: React.FC<LogsDialogProps> = ({
  open,
  onOpenChange,
  instanceName,
  isRunning
}) => {
  const [logs, setLogs] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [lineCount, setLineCount] = useState(100)
  const [autoRefresh, setAutoRefresh] = useState(false)
  const [copied, setCopied] = useState(false)
  const [showSettings, setShowSettings] = useState(false)
  
  const logContainerRef = useRef<HTMLDivElement>(null)
  const refreshIntervalRef = useRef<NodeJS.Timeout | null>(null)

  // Fetch logs function
  const fetchLogs = React.useCallback(
    async (lines?: number) => {
      if (!instanceName) return
      
      setLoading(true)
      setError(null)
      
      try {
        const logText = await instancesApi.getLogs(instanceName, lines)
        setLogs(logText)
        
        // Auto-scroll to bottom
        setTimeout(() => {
          if (logContainerRef.current) {
            logContainerRef.current.scrollTop = logContainerRef.current.scrollHeight
          }
        }, 100)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch logs')
      } finally {
        setLoading(false)
      }
    },
    [instanceName]
  )

  // Initial load when dialog opens
  useEffect(() => {
    if (open && instanceName) {
      void fetchLogs(lineCount)
    }
  }, [open, instanceName, fetchLogs, lineCount])

