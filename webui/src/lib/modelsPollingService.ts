import type { DownloadJob } from '@/types/model'
import { llamaCppModelsApi } from '@/lib/api'

type JobsCallback = (jobs: DownloadJob[]) => void

class ModelsPollingService {
  private interval: NodeJS.Timeout | null = null
  private callbacks: Set<JobsCallback> = new Set()
  private readonly POLL_INTERVAL = 2000 // 2 seconds
  private expectingJobs = false // Flag to keep polling when a download just started
  private emptyPollCount = 0 // Count consecutive empty polls
  private node?: string // Optional node to poll jobs from

  /**
   * Subscribe to job updates
   */
  subscribe(callback: JobsCallback, node?: string): () => void {
    this.callbacks.add(callback)
    this.node = node

    // Start polling if this is the first subscriber
    if (this.callbacks.size === 1) {
      this.startPolling()
    }

    // Return unsubscribe function
    return () => {
      this.callbacks.delete(callback)

      // Stop polling if no more subscribers
      if (this.callbacks.size === 0) {
        this.stopPolling()
      }
    }
  }

  /**
   * Poll for active jobs
   */
  private async pollJobs(): Promise<void> {
    try {
      const { jobs } = await llamaCppModelsApi.listJobs(this.node)

      // Filter to active and recently completed/failed jobs (exclude only completed)
      const activeJobs = jobs.filter(j =>
        j.status === 'downloading' || j.status === 'queued' || j.status === 'failed' || j.status === 'cancelled'
      )

      // Notify all callbacks
      this.callbacks.forEach(cb => { cb(activeJobs) })

      // Handle stopping logic
      if (activeJobs.length === 0) {
