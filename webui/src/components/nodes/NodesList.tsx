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
        className="flex items-center justify-center py-12"
        role="status"
        aria-label="Loading"
      >
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Loading nodes...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <div className="text-destructive mb-4">
          <p className="text-lg font-semibold">Error loading nodes</p>
          <p className="text-sm">{error}</p>
        </div>
      </div>
    );
  }

  const entries = nodes ? Object.entries(nodes) : [];
