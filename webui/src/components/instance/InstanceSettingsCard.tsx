import type React from 'react'
import { useState, useEffect } from 'react'
import type { CreateInstanceOptions } from '@/types/instance'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import AutoRestartConfiguration from '@/components/instance/AutoRestartConfiguration'
import NumberInput from '@/components/form/NumberInput'
import CheckboxInput from '@/components/form/CheckboxInput'
import SelectInput from '@/components/form/SelectInput'
import { nodesApi, type NodesMap } from '@/lib/api'

interface InstanceSettingsCardProps {
  instanceName: string
  nameError: string
  isEditing: boolean
  formData: CreateInstanceOptions
  onNameChange: (name: string) => void
  onChange: (key: keyof CreateInstanceOptions, value: unknown) => void
}

const InstanceSettingsCard: React.FC<InstanceSettingsCardProps> = ({
  instanceName,
  nameError,
  isEditing,
  formData,
  onNameChange,
  onChange
}) => {
  const [nodes, setNodes] = useState<NodesMap>({})
  const [loadingNodes, setLoadingNodes] = useState(true)

  // biome-ignore lint/correctness/useExhaustiveDependencies: fetch nodes once on mount; intentionally not re-run when formData/onChange change
  useEffect(() => {
    const fetchNodes = async () => {
      try {
