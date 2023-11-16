import { useModels } from '@/contexts/ModelsContext'
import { Alert, AlertDescription } from '@/components/ui/alert'
import ModelsHeader from './models/ModelsHeader'
import ModelsTable from './models/ModelsTable'
import type { ModelRow } from '@/types/model'

export default function ModelsList() {
  const { models, activeJobs, loading, error, clearError } = useModels()

  // Combine active jobs and cached models into unified rows
  const rows: ModelRow[] = [
    ...(activeJobs || []).map((job) => ({ type: 'downloading' as const, job })),
    ...(models || []).map((model) => ({ type: 'cached' as const, model })),
  ]

  if (loading) {
    return (
      <div className="space-y-4">
        <ModelsHeader />
        <div className="space-y-2">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-16 bg-muted animate-pulse rounded" />
          ))}
        </div>
      </div>
    )
  }

  return (
