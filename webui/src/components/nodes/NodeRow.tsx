// ui/src/components/nodes/NodeRow.tsx

import { CheckCircle, Loader2, XCircle } from 'lucide-react';
import { memo } from 'react';
import { TableCell, TableRow } from '@/components/ui/table';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { useNodeHealth } from '@/hooks/useNodeHealth';
import type { NodeHealth } from '@/types/node';

interface NodeRowProps {
  name: string;
  address: string;
  isLocal: boolean;
}
