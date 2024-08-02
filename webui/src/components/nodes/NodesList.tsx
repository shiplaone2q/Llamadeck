// ui/src/components/nodes/NodesList.tsx

import { HardDrive } from 'lucide-react';
import { useEffect, useState } from 'react';
import NodeRow from '@/components/nodes/NodeRow';
import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { TooltipProvider } from '@/components/ui/tooltip';
import { useConfig } from '@/contexts/ConfigContext';
import { type NodesMap, nodesApi } from '@/lib/api';

function NodesList() {
  const { config } = useConfig();
  const [nodes, setNodes] = useState<NodesMap | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    const loadNodes = async () => {
      try {
        const data = await nodesApi.list();
        if (!cancelled) {
          setNodes(data);
          setError(null);
        }
      } catch (err) {
        if (!cancelled) {
          setError(
            err instanceof Error ? err.message : 'Failed to fetch nodes',
          );
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    };

    void loadNodes();
    return () => {
      cancelled = true;
    };
  }, []);

  if (loading) {
    return (
      <div
