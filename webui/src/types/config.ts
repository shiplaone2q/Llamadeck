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

export interface InstancesConfig {
  port_range: [number, number]
  logs_dir: string
  auto_create_dirs: boolean
  max_instances: number
  max_running_instances: number
  enable_lru_eviction: boolean
  default_idle_timeout: number
  default_auto_restart: boolean
  default_max_restarts: number
  default_restart_delay: number
  default_on_demand_start: boolean
  on_demand_start_timeout: number
  timeout_check_interval: number
}

export interface DatabaseConfig {
