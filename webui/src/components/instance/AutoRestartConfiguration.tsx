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

  return (
    <div className="space-y-4">
      {/* biome-ignore lint/a11y/useSemanticElements: toggle contains an <h3>, so a <button> would be invalid HTML; keyboard support provided via role/tabIndex/onKeyDown */}
      <div
        role="button"
        tabIndex={0}
        className="flex items-center justify-between cursor-pointer"
        onClick={() => setIsExpanded(!isExpanded)}
        onKeyDown={(e) => {
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault()
            setIsExpanded(!isExpanded)
          }
        }}
      >
        <h3 className="text-lg font-medium">Auto Restart Configuration</h3>
        <div className="flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground">
          {isExpanded ? (
            <ChevronDown className="h-4 w-4" />
          ) : (
            <ChevronRight className="h-4 w-4" />
          )}
          {getSummaryText()}
        </div>
      </div>
