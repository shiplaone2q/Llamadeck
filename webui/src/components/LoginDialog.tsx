import type React from 'react'
import { useState, useEffect } from 'react'
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
import { AlertCircle, Fingerprint, Key, Loader2 } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'

interface LoginDialogProps {
  open: boolean
  onOpenChange?: (open: boolean) => void
}

const LoginDialog: React.FC<LoginDialogProps> = ({
  open,
  onOpenChange,
}) => {
  const { login, isLoading, error, clearError, oidcEnabled, oidcError, clearOIDCError } = useAuth()
  const [apiKey, setApiKey] = useState('')
  const [localLoading, setLocalLoading] = useState(false)
  const [showKeyLogin, setShowKeyLogin] = useState(!oidcEnabled)

  // Clear form and errors when dialog opens/closes
  useEffect(() => {
    if (open) {
      setApiKey('')
      setShowKeyLogin(!oidcEnabled)
      clearError()
    }
  }, [open, oidcEnabled, clearError])

  // Clear errors when user starts typing
  useEffect(() => {
    if ((error || oidcError) && apiKey) {
      clearError()
      clearOIDCError()
    }
  }, [apiKey, error, oidcError, clearError, clearOIDCError])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!apiKey.trim()) {
      return
    }

    setLocalLoading(true)
    
    try {
      await login(apiKey.trim())
      // Login successful - dialog will close automatically when auth state changes
      setApiKey('')
    } catch (err) {
      // Error is handled by the AuthContext
      console.error('Login failed:', err)
    } finally {
      setLocalLoading(false)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && !isSubmitDisabled) {
