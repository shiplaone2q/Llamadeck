import { useState, useEffect } from 'react'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import SelectInput from '@/components/form/SelectInput'
import { useModels } from '@/contexts/ModelsContext'
import { nodesApi, type NodesMap } from '@/lib/api'
import { ModelFormat } from '@/types/model'

interface DownloadDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

const formatOptions = [
  { value: ModelFormat.GGUF, label: 'GGUF' },
  { value: ModelFormat.SAFETENSORS, label: 'Safetensors' },
]

export default function DownloadDialog({ open, onOpenChange }: DownloadDialogProps) {
  const [model, setModel] = useState('')
  const [modelError, setModelError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [nodes, setNodes] = useState<NodesMap>({})
  const [loadingNodes, setLoadingNodes] = useState(true)
  const [selectedNode, setSelectedNode] = useState<string | undefined>(undefined)
  const [format, setFormat] = useState<ModelFormat>(ModelFormat.GGUF)
  const { startDownload } = useModels()

  // biome-ignore lint/correctness/useExhaustiveDependencies: open intentionally included to refetch nodes when the dialog opens
  useEffect(() => {
    const fetchNodes = async () => {
      try {
        const fetchedNodes = await nodesApi.list()
        setNodes(fetchedNodes)

        const nodeNames = Object.keys(fetchedNodes)
        if (nodeNames.length > 0) {
          if (!selectedNode || !nodeNames.includes(selectedNode)) {
            setSelectedNode(nodeNames[0])
          }
        } else if (selectedNode) {
          setSelectedNode(undefined)
        }
      } catch (error) {
        console.error('Failed to fetch nodes:', error)
      } finally {
        setLoadingNodes(false)
      }
    }

    void fetchNodes()
  }, [open, selectedNode])

  const validateModel = (value: string): boolean => {
    if (!value) {
      setModelError('Model is required')
      return false
    }
    if (format === ModelFormat.GGUF && !value.includes('/')) {
      setModelError('Format must be: org/model-name or org/model-name:tag')
      return false
    }
    if (!value.includes('/')) {
      setModelError('Format must be: org/model-name or org/model-name:tag')
      return false
    }
    setModelError('')
