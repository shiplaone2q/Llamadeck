import { type ReactNode, createContext, useContext, useState, useEffect, useCallback } from 'react'
import type { CreateInstanceOptions, Instance } from '@/types/instance'
import { instancesApi } from '@/lib/api'
import { useAuth } from '@/contexts/AuthContext'
import { healthService } from '@/lib/healthService'

interface InstancesContextState {
  instances: Instance[]
  loading: boolean
  error: string | null
}

interface InstancesContextActions {
  fetchInstances: () => Promise<void>
  createInstance: (name: string, options: CreateInstanceOptions) => Promise<void>
  updateInstance: (name: string, options: CreateInstanceOptions) => Promise<void>
  startInstance: (name: string) => Promise<void>
  stopInstance: (name: string) => Promise<void>
  restartInstance: (name: string) => Promise<void>
  deleteInstance: (name: string) => Promise<void>
  clearError: () => void
}

type InstancesContextType = InstancesContextState & InstancesContextActions

const InstancesContext = createContext<InstancesContextType | undefined>(undefined)

interface InstancesProviderProps {
  children: ReactNode
}

export const InstancesProvider = ({ children }: InstancesProviderProps) => {
  const { isAuthenticated, isLoading: authLoading } = useAuth()
  const [instancesMap, setInstancesMap] = useState<Map<string, Instance>>(new Map())
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Convert map to array for consumers
  const instances = Array.from(instancesMap.values())

  const clearError = useCallback(() => {
    setError(null)
  }, [])

  const fetchInstances = useCallback(async () => {
    if (!isAuthenticated) {
      setLoading(false)
      return
    }

    try {
      setLoading(true)
      setError(null)
      const data = await instancesApi.list()
      
      // Convert array to map
      const newMap = new Map<string, Instance>()
      data.forEach(instance => {
        newMap.set(instance.name, instance)
      })
      setInstancesMap(newMap)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch instances')
    } finally {
