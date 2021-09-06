import { type ReactNode, createContext, useContext, useState, useEffect, useCallback } from 'react'
import type { CachedModel, DownloadJob, ModelFormat } from '@/types/model'
import { llamaCppModelsApi } from '@/lib/api'
import { useAuth } from '@/contexts/AuthContext'
import { modelsPollingService } from '@/lib/modelsPollingService'

interface ModelsContextType {
  // State
  models: CachedModel[]
  activeJobs: DownloadJob[]
  loading: boolean
  error: string | null

  // Actions
  fetchModels: (node?: string) => Promise<void>
  startDownload: (repo: string, tag?: string, node?: string, format?: ModelFormat) => Promise<void>
  cancelDownload: (jobId: string, node?: string) => Promise<void>
  deleteDownload: (jobId: string, node?: string) => Promise<void>
  deleteModel: (repo: string, tag?: string, node?: string) => Promise<void>
  clearError: () => void
}

const ModelsContext = createContext<ModelsContextType | undefined>(undefined)

interface ModelsProviderProps {
  children: ReactNode
}

export const ModelsProvider = ({ children }: ModelsProviderProps) => {
  const { isAuthenticated, isLoading: authLoading } = useAuth()
  const [models, setModels] = useState<CachedModel[]>([])
  const [activeJobs, setActiveJobs] = useState<DownloadJob[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const clearError = useCallback(() => {
    setError(null)
  }, [])

  const fetchModels = useCallback(async (node?: string) => {
    if (!isAuthenticated) {
      setLoading(false)
      return
    }

    try {
      setLoading(true)
      setError(null)
      const modelsData = await llamaCppModelsApi.listModels(node)
      setModels(modelsData)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch models')
    } finally {
      setLoading(false)
    }
  }, [isAuthenticated])

