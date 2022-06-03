import { useState } from 'react'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import type { ModelRow, DownloadJob, CachedModel, ModelFile } from '@/types/model'
import { formatBytes } from '@/lib/utils'
import { useModels } from '@/contexts/ModelsContext'
import {
  Download,
  X,
  Trash2,
  Loader2,
  ChevronDown,
  ChevronRight,
  FileCode,
  FileText,
  Image,
  RefreshCw,
} from 'lucide-react'
import StatusBadge from './StatusBadge'

interface ModelsTableProps {
  rows: ModelRow[]
}

export default function ModelsTable({ rows }: ModelsTableProps) {
  const [expandedRows, setExpandedRows] = useState<Set<string>>(new Set())

  const toggleRow = (key: string) => {
    setExpandedRows(prev => {
      const newSet = new Set(prev)
      if (newSet.has(key)) {
        newSet.delete(key)
      } else {
        newSet.add(key)
      }
      return newSet
    })
  }

  if (rows.length === 0) {
    return (
      <div className="text-center py-12 text-muted-foreground">
        <Download className="h-12 w-12 mx-auto mb-4 opacity-50" />
        <p className="text-lg mb-2">No models yet</p>
        <p className="text-sm">Download your first model to get started</p>
      </div>
    )
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Repository</TableHead>
          <TableHead>Tag</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Node</TableHead>
          <TableHead>Size</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {rows.map((row) => {
          if (row.type === 'downloading') {
            return <DownloadingRow key={row.job.id} job={row.job} />
          } else {
            const key = `${row.model.repo}:${row.model.tag}`
            const expanded = expandedRows.has(key)
            return (
              <>
                <CachedModelRow
                  key={key}
                  model={row.model}
                  expanded={expanded}
                  onToggle={() => toggleRow(key)}
                />
                {expanded && (
                  <ExpandedFilesRow key={`${key}-files`} files={row.model.files} />
                )}
              </>
            )
          }
        })}
      </TableBody>
    </Table>
  )
}

function DownloadingRow({ job }: { job: DownloadJob }) {
  const { cancelDownload, deleteDownload, startDownload } = useModels()
  const [cancelling, setCancelling] = useState(false)
  const [deleting, setDeleting] = useState(false)
