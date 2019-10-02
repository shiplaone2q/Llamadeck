import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Download } from 'lucide-react'
import { useModels } from '@/contexts/ModelsContext'
import DownloadDialog from './DownloadDialog'

export default function ModelsHeader() {
  const { models, activeJobs } = useModels()
  const [downloadDialogOpen, setDownloadDialogOpen] = useState(false)

