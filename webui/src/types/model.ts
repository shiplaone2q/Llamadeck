export enum ModelFormat {
  GGUF = 'gguf',
  SAFETENSORS = 'safetensors'
}

export interface DownloadJob {
  id: string
  repo: string
  tag: string
  status: 'queued' | 'downloading' | 'completed' | 'failed' | 'cancelled'
