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

function NodeStatus({ health }: { health: NodeHealth | undefined }) {
  if (!health) {
    return (
      <span className="flex items-center gap-1.5 text-muted-foreground">
        <Loader2 className="h-3.5 w-3.5 animate-spin" />
        Checking
      </span>
    );
  }

  if (health.state === 'healthy') {
    return (
      <span className="flex items-center gap-1.5 text-green-600 dark:text-green-400">
        <CheckCircle className="h-3.5 w-3.5" />
        Healthy
      </span>
    );
  }

  return (
    <span
      className="flex items-center gap-1.5 text-destructive"
      title={health.error || undefined}
    >
      <XCircle className="h-3.5 w-3.5" />
      Unreachable
    </span>
  );
}

function NodeRow({ name, address, isLocal }: NodeRowProps) {
  const health = useNodeHealth(name);

  // local node is never pinged; latency is meaningless for it
  const latency =
    !isLocal && health?.state === 'healthy' && health.latencyMs !== undefined
      ? `${health.latencyMs}ms`
      : '—';

  return (
    <TableRow>
      <TableCell className="font-medium">{name}</TableCell>
      <TableCell>
        {isLocal ? (
          <span className="text-muted-foreground">Local</span>
        ) : (
          <Tooltip>
