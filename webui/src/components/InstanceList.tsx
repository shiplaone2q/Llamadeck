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
