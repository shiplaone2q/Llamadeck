import type React from 'react'
import { useState, useEffect, useCallback } from 'react'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import SelectInput from '@/components/form/SelectInput'
import {
  RefreshCw,
  AlertCircle,
  Loader2,
  ChevronDown,
  ChevronRight,
  Monitor,
  HelpCircle,
  Info
} from 'lucide-react'
import { serverApi } from '@/lib/api'
import { BackendType, type BackendTypeValue } from '@/types/instance'

// Helper to get version from environment
const getAppVersion = (): string => {
  try {
    return (import.meta.env as Record<string, string>).VITE_APP_VERSION || 'unknown'
  } catch {
    return 'unknown'
  }
}

interface SystemInfoDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

interface BackendInfo {
  version: string
  devices: string
  help: string
}

const BACKEND_OPTIONS = [
  { value: BackendType.LLAMA_CPP, label: 'Llama Server' },
  { value: BackendType.MLX_LM, label: 'MLX LM' },
  { value: BackendType.VLLM, label: 'vLLM' },
]

const SystemInfoDialog: React.FC<SystemInfoDialogProps> = ({
  open,
  onOpenChange
}) => {
  const [selectedBackend, setSelectedBackend] = useState<BackendTypeValue>(BackendType.LLAMA_CPP)
  const [backendInfo, setBackendInfo] = useState<BackendInfo | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [showHelp, setShowHelp] = useState(false)

  // Fetch backend info
  const fetchBackendInfo = useCallback(async (backend: BackendTypeValue) => {
    if (backend !== BackendType.LLAMA_CPP) {
      setBackendInfo(null)
      setError(null)
      return
    }

    setLoading(true)
    setError(null)

    try {
      const [version, devices, help] = await Promise.all([
        serverApi.getVersion(),
        serverApi.getDevices(),
        serverApi.getHelp()
      ])

      setBackendInfo({ version, devices, help })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch backend info')
    } finally {
      setLoading(false)
    }
  }, [])

  // Load data when dialog opens or backend changes
  useEffect(() => {
    if (open) {
      void fetchBackendInfo(selectedBackend)
    }
  }, [open, selectedBackend, fetchBackendInfo])

  const handleBackendChange = (value: string) => {
    setSelectedBackend(value as BackendTypeValue)
    setShowHelp(false) // Reset help section when switching backends
  }

  const renderBackendSpecificContent = () => {
    if (selectedBackend !== BackendType.LLAMA_CPP) {
      return (
        <div className="flex items-center justify-center py-8">
          <div className="text-center space-y-3">
            <Info className="h-8 w-8 text-gray-400 mx-auto" />
            <div>
              <h3 className="font-semibold text-gray-700">Backend Info Not Available</h3>
              <p className="text-sm text-gray-500 mt-1">
                Information for {BACKEND_OPTIONS.find(b => b.value === selectedBackend)?.label} backend is not yet implemented.
              </p>
            </div>
          </div>
        </div>
      )
    }

    if (loading && !backendInfo) {
      return (
