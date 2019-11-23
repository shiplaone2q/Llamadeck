import type React from 'react'
import { useState } from 'react'
import type { CreateInstanceOptions } from '@/types/instance'
import CheckboxInput from '@/components/form/CheckboxInput'
import NumberInput from '@/components/form/NumberInput'
import { ChevronDown, ChevronRight } from 'lucide-react'

interface AutoRestartConfigurationProps {
  formData: CreateInstanceOptions
  onChange: <K extends keyof CreateInstanceOptions>(key: K, value: CreateInstanceOptions[K]) => void
}

const AutoRestartConfiguration: React.FC<AutoRestartConfigurationProps> = ({
  formData,
  onChange
}) => {
  const [isExpanded, setIsExpanded] = useState(false)
  const isAutoRestartEnabled = formData.auto_restart === true

  const getSummaryText = () => {
    if (!isAutoRestartEnabled) {
      return 'Disabled'
    }
    const maxRestarts = formData.max_restarts ?? 3
    const restartDelay = formData.restart_delay ?? 5
    const restartsText = maxRestarts === 0 ? 'unlimited' : `up to ${maxRestarts}`
    return `Enabled (restarts: ${restartsText}, delay: ${restartDelay}s)`
  }

