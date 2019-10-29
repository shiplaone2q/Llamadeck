// ui/src/components/InstanceList.tsx
import { useInstances } from '@/contexts/InstancesContext'
import InstanceCard from '@/components/InstanceCard'
import type { Instance } from '@/types/instance'
import { memo } from 'react'

interface InstanceListProps {
  editInstance: (instance: Instance) => void
}

// Memoize InstanceCard to prevent re-renders when other instances change
const MemoizedInstanceCard = memo(InstanceCard)

function InstanceList({ editInstance }: InstanceListProps) {
  const { instances, loading, error, startInstance, stopInstance, deleteInstance } = useInstances()

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12" role="status" aria-label="Loading">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
