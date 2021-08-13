export interface BackendSettings {
  command: string
  args: string[]
  environment?: Record<string, string>
  docker?: DockerSettings
  response_headers?: Record<string, string>
}

export interface DockerSettings {
  enabled: boolean
  image: string
  args: string[]
  environment?: Record<string, string>
}

export interface BackendConfig {
  'llama-cpp': BackendSettings
  vllm: BackendSettings
  mlx: BackendSettings
  custom?: Record<string, BackendSettings>
}

export interface ServerConfig {
  host: string
  port: number
  allowed_origins: string[]
  allowed_headers: string[]
  enable_swagger: boolean
  response_headers?: Record<string, string>
}

