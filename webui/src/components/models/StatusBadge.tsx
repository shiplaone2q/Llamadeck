import { Badge } from "@/components/ui/badge"
import type { DownloadJob } from "@/types/model"

interface StatusBadgeProps {
  status: DownloadJob['status']
}

export default function StatusBadge({ status }: StatusBadgeProps) {
