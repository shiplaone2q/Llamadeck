# Configuration

llamadeck can be configured via configuration files or environment variables. Configuration is loaded in the following order of precedence:

```
Defaults < .env file < Configuration file (with env var expansion) < LLAMADECK_* environment variables
```

llamadeck works out of the box with sensible defaults, but you can customize the behavior to suit your needs.

## Default Configuration

Here's the default configuration with all available options:

```yaml
server:
  host: "0.0.0.0"                # Server host to bind to
  port: 8080                     # Server port to bind to
  allowed_origins: ["*"]         # Allowed CORS origins (default: all)
  allowed_headers: ["*"]         # Allowed CORS headers (default: all)
  enable_swagger: false          # Enable Swagger UI for API docs

backends:
  llama-cpp:
    command: "llama-server"
    args: []
    environment: {}              # Environment variables for the backend process
    docker:
      enabled: false
      image: "ghcr.io/ggml-org/llama.cpp:server"
      args: ["run", "--rm", "--network", "host", "--gpus", "all"]
      environment: {}
    response_headers: {}         # Additional response headers to send with responses

  vllm:
    command: "vllm"
    args: ["serve"]
    environment: {}              # Environment variables for the backend process
    docker:
      enabled: false
      image: "vllm/vllm-openai:latest"
      args: ["run", "--rm", "--network", "host", "--gpus", "all", "--shm-size", "1g"]
      environment: {}
    response_headers: {}         # Additional response headers to send with responses

  mlx:
    command: "mlx_lm.server"
    args: []
    environment: {}              # Environment variables for the backend process
    response_headers: {}         # Additional response headers to send with responses

data_dir: ~/.local/share/llamadeck  # Main data directory (database, instances, logs), default varies by OS

instances:
  port_range: [8000, 9000]         # Port range for instances
  configs_dir: data_dir/instances  # Instance configs directory
  logs_dir: data_dir/logs          # Logs directory
  auto_create_dirs: true           # Auto-create data/config/logs dirs if missing
  max_instances: -1                # Max instances (-1 = unlimited)
  max_running_instances: -1        # Max running instances (-1 = unlimited)
  enable_lru_eviction: true        # Enable LRU eviction for idle instances
  default_idle_timeout: 30         # Default idle timeout in minutes (0 = no timeout)
  default_auto_restart: true       # Auto-restart new instances by default
  default_max_restarts: 3          # Max restarts for new instances
  default_restart_delay: 5         # Restart delay (seconds) for new instances
  default_on_demand_start: true    # Default on-demand start setting
  on_demand_start_timeout: 120     # Default on-demand start timeout in seconds
  timeout_check_interval: 5        # Idle instance timeout check in minutes
  group_limits: {}                 # Per-group running instance limits (e.g., {large: 1, small: 3})

database:
  path: data_dir/llamadeck.db              # Database file path
  max_open_connections: 25       # Maximum open database connections
  max_idle_connections: 5        # Maximum idle database connections
  connection_max_lifetime: 5m    # Connection max lifetime

auth:
  require_inference_auth: true   # Require auth for inference endpoints
  require_management_auth: true  # Require auth for management endpoints
  management_keys: []            # Keys for management endpoints

local_node: "main"               # Name of the local node (default: "main")
nodes:                           # Node configuration for multi-node deployment
  main:                          # Default local node (empty config)
```

## Configuration Files

### Configuration File Locations

Configuration files are searched in the following locations (in order of precedence, first found is used):

**Linux:**  
