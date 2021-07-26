import { type ReactNode, createContext, useContext, useEffect, useState, useRef } from 'react'
import { serverApi } from '@/lib/api'
import type { AppConfig } from '@/types/config'
import { useAuth } from './AuthContext'

interface ConfigContextType {
  config: AppConfig | null
  isLoading: boolean
  error: string | null
}

const ConfigContext = createContext<ConfigContextType | undefined>(undefined)

interface ConfigProviderProps {
  children: ReactNode
}

export const ConfigProvider = ({ children }: ConfigProviderProps) => {
  const { isAuthenticated } = useAuth()
  const [config, setConfig] = useState<AppConfig | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const loadedRef = useRef(false)

  useEffect(() => {
    if (!isAuthenticated || loadedRef.current) {
      setIsLoading(false)
      return
