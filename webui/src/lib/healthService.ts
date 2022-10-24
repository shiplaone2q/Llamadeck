import type { HealthStatus, InstanceStatus, HealthState } from '@/types/instance'
import { instancesApi } from '@/lib/api'

type HealthCallback = (health: HealthStatus) => void

// Polling intervals based on health state (in milliseconds)
const POLLING_INTERVALS: Record<HealthState, number> = {
  'starting': 5000,      // 5 seconds - frequent during startup
  'restarting': 5000,    // 5 seconds - restart in progress
  'shutting_down': 3000, // 3 seconds - monitor shutdown progress
  'ready': 60000,        // 60 seconds - stable state
  'stopped': 0,          // No polling
  'failed': 0,           // No polling
}

class HealthService {
  private intervals: Map<string, NodeJS.Timeout> = new Map()
  private callbacks: Map<string, Set<HealthCallback>> = new Map()
  private lastHealthState: Map<string, HealthState> = new Map()
  private healthCache: Map<string, { health: HealthStatus; timestamp: number }> = new Map()
  private readonly CACHE_TTL = 2000 // 2 seconds cache

  /**
   * Performs a two-tier health check:
   * 1. Get instance status from backend (authoritative)
   * 2. If running, perform HTTP health check
   */
  async performHealthCheck(instanceName: string): Promise<HealthStatus> {
    // Check cache first
    const cached = this.healthCache.get(instanceName)
    if (cached && Date.now() - cached.timestamp < this.CACHE_TTL) {
      return cached.health
    }

    try {
      // Step 1: Get instance details (includes status)
      const instance = await instancesApi.get(instanceName)

      // Step 2: If running, attempt HTTP health check
      if (instance.status === 'running') {
        try {
          await instancesApi.getHealth(instanceName)

          // HTTP health check succeeded - instance is ready
          const health: HealthStatus = {
            state: 'ready',
            instanceStatus: 'running',
            lastChecked: new Date(),
            source: 'http'
          }

          this.updateCache(instanceName, health)
          return health

        } catch (httpError) {
          // HTTP health check failed - instance is still starting
          // Any error (503, connection refused, timeout, etc.) means "starting"
          const health: HealthStatus = {
            state: 'starting',
            instanceStatus: 'running',
            lastChecked: new Date(),
            error: httpError instanceof Error ? httpError.message : 'Health check failed',
            source: 'http'
          }

          this.updateCache(instanceName, health)
          return health
        }
      } else {
        // Instance not running - map backend status directly
        const health: HealthStatus = {
          state: this.mapStatusToHealthState(instance.status),
          instanceStatus: instance.status,
