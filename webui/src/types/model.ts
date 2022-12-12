export enum ModelFormat {
  GGUF = 'gguf',
  SAFETENSORS = 'safetensors'
}

export interface DownloadJob {
  id: string
  repo: string
  tag: string
  status: 'queued' | 'downloading' | 'completed' | 'failed' | 'cancelled'
  progress: {
    bytes_downloaded: number
    total_bytes: number
    current_file: string
  }
  error: string | null
  created_at: number
  completed_at: number | null
}

export interface CachedModel {
  node: string
  repo: string
  tag: string
