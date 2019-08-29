import { nodesApi } from '@/lib/api';
import type { NodeHealth } from '@/types/node';

type NodeHealthCallback = (health: NodeHealth) => void;

const POLL_INTERVAL = 10000;
const CACHE_TTL = 2000;

// Polls node health while the Nodes tab is mounted. First subscriber starts
// polling, last unsubscribe stops it.
class NodeHealthService {
  private intervals: Map<string, ReturnType<typeof setInterval>> = new Map();
  private callbacks: Map<string, Set<NodeHealthCallback>> = new Map();
  private healthCache: Map<string, { health: NodeHealth; timestamp: number }> =
    new Map();

  async performHealthCheck(nodeName: string): Promise<NodeHealth> {
    const cached = this.healthCache.get(nodeName);
    if (cached && Date.now() - cached.timestamp < CACHE_TTL) {
      return cached.health;
    }

    let health: NodeHealth;
    try {
      const response = await nodesApi.getHealth(nodeName);
