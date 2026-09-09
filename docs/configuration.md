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
- `./llamadeck.yaml` or `./config.yaml` (current directory)  
- `$HOME/.config/llamadeck/config.yaml`  
- `/etc/llamadeck/config.yaml`  

**macOS:**  
- `./llamadeck.yaml` or `./config.yaml` (current directory)  
- `$HOME/Library/Application Support/llamadeck/config.yaml`  
- `/Library/Application Support/llamadeck/config.yaml`  

**Windows:**  
- `./llamadeck.yaml` or `./config.yaml` (current directory)  
- `%APPDATA%\llamadeck\config.yaml`  
- `%USERPROFILE%\llamadeck\config.yaml`  
- `%PROGRAMDATA%\llamadeck\config.yaml`  

You can specify the path to config file with `LLAMADECK_CONFIG_PATH` environment variable.

### Environment Variable Expansion

Config files support `${VAR}` and `${VAR:-default}` placeholders, resolved from the environment before parsing. Unset variables with no default are left as-is. Only `${VAR}` syntax is supported (not `$VAR`).

```yaml
auth:
  management_keys:
    - ${OPENAI_API_KEY}

backends:
  llama-cpp:
    environment:
      CUDA_VISIBLE_DEVICES: ${CUDA_DEVICE:-0}
```

### Dotenv File Loading

A `.env` file is automatically loaded before config parsing. It's searched in the same directories as config files (first match wins), plus `./.env`. Existing environment variables are never overwritten.

```env
OPENAI_API_KEY=sk-proj-abc123...
CUDA_DEVICE=0
```

## Configuration Options

### Server Configuration

```yaml
server:
  host: "0.0.0.0"         # Server host to bind to (default: "0.0.0.0")
  port: 8080              # Server port to bind to (default: 8080)
  allowed_origins: ["*"]  # CORS allowed origins (default: ["*"])
  allowed_headers: ["*"]  # CORS allowed headers (default: ["*"])
  enable_swagger: false   # Enable Swagger UI (default: false)
```

**Environment Variables:**
- `LLAMADECK_HOST` - Server host
- `LLAMADECK_PORT` - Server port
- `LLAMADECK_ALLOWED_ORIGINS` - Comma-separated CORS origins
- `LLAMADECK_ENABLE_SWAGGER` - Enable Swagger UI (true/false)

### Backend Configuration
```yaml
backends:
  llama-cpp:
    command: "llama-server"
    args: []
    environment: {}              # Environment variables for the backend process
    docker:
      enabled: false             # Enable Docker runtime (default: false)
      image: "ghcr.io/ggml-org/llama.cpp:server"
      args: ["run", "--rm", "--network", "host", "--gpus", "all"]
      environment: {}
    response_headers: {}         # Additional response headers to send with responses

  vllm:
    command: "vllm"
    args: ["serve"]
    environment: {}              # Environment variables for the backend process
    docker:
      enabled: false             # Enable Docker runtime (default: false)
      image: "vllm/vllm-openai:latest"
      args: ["run", "--rm", "--network", "host", "--gpus", "all", "--shm-size", "1g"]
      environment: {}
    response_headers: {}         # Additional response headers to send with responses

  mlx:
    command: "mlx_lm.server"
    args: []
    environment: {}              # Environment variables for the backend process
    # MLX does not support Docker
    response_headers: {}         # Additional response headers to send with responses

  custom: {}                     # Named custom backends, see "Custom Backends" below
```

**Backend Configuration Fields:**
- `command`: Executable name/path for the backend
- `args`: Default arguments prepended to all instances
- `environment`: Environment variables for the backend process (optional)
- `response_headers`: Additional response headers to send with responses (optional)
- `health_path`: Path polled for readiness; any 2xx response counts as healthy (default: `/health`)
- `docker`: Docker-specific configuration (optional)
  - `enabled`: Boolean flag to enable Docker runtime
  - `image`: Docker image to use
  - `args`: Additional arguments passed to `docker run`
  - `environment`: Environment variables for the container (optional)

> If llamadeck is behind an NGINX proxy, `X-Accel-Buffering: no` response header may be required for NGINX to properly stream the responses without buffering.

**Environment Variables:**

**LlamaCpp Backend:**
- `LLAMADECK_LLAMACPP_COMMAND` - LlamaCpp executable command
- `LLAMADECK_LLAMACPP_ARGS` - Space-separated default arguments
- `LLAMADECK_LLAMACPP_ENV` - Environment variables in format "KEY1=value1,KEY2=value2"
- `LLAMADECK_LLAMACPP_DOCKER_ENABLED` - Enable Docker runtime (true/false)
- `LLAMADECK_LLAMACPP_DOCKER_IMAGE` - Docker image to use
- `LLAMADECK_LLAMACPP_DOCKER_ARGS` - Space-separated Docker arguments
- `LLAMADECK_LLAMACPP_DOCKER_ENV` - Docker environment variables in format "KEY1=value1,KEY2=value2"
- `LLAMADECK_LLAMACPP_RESPONSE_HEADERS` - Response headers in format "KEY1=value1;KEY2=value2"

**VLLM Backend:**
- `LLAMADECK_VLLM_COMMAND` - VLLM executable command
- `LLAMADECK_VLLM_ARGS` - Space-separated default arguments
- `LLAMADECK_VLLM_ENV` - Environment variables in format "KEY1=value1,KEY2=value2"
- `LLAMADECK_VLLM_DOCKER_ENABLED` - Enable Docker runtime (true/false)
- `LLAMADECK_VLLM_DOCKER_IMAGE` - Docker image to use
- `LLAMADECK_VLLM_DOCKER_ARGS` - Space-separated Docker arguments
- `LLAMADECK_VLLM_DOCKER_ENV` - Docker environment variables in format "KEY1=value1,KEY2=value2"
- `LLAMADECK_VLLM_RESPONSE_HEADERS` - Response headers in format "KEY1=value1;KEY2=value2"

**MLX Backend:**
- `LLAMADECK_MLX_COMMAND` - MLX executable command
- `LLAMADECK_MLX_ARGS` - Space-separated default arguments
- `LLAMADECK_MLX_ENV` - Environment variables in format "KEY1=value1,KEY2=value2"
- `LLAMADECK_MLX_RESPONSE_HEADERS` - Response headers in format "KEY1=value1;KEY2=value2"

### Custom Backends

Custom backends are user-defined named entries under `backends.custom` for managing inference servers without native support. Each entry accepts the same fields as the built-in backends, and `args` support the `{port}` and `{model}` placeholders; see the [Custom Backends](custom-backends.md) guide for how they work.

```yaml
backends:
  custom:
    my-engine:
      command: "my-server"
      args: []                     # Default arguments prepended to instance args
      environment: {}              # Environment variables for the backend process
      health_path: "/health"       # Path polled for readiness (default: /health)
      docker:
        enabled: false             # Enable Docker runtime (default: false)
        image: "my-server:latest"
        args: ["run", "--rm", "--network", "host", "--gpus", "all"]
        environment: {}
      response_headers: {}         # Additional response headers to send with responses
```

### Data Directory Configuration

```yaml
data_dir: "~/.local/share/llamadeck"  # Main data directory for database, instances, and logs (default varies by OS)
```

**Environment Variables:**
- `LLAMADECK_DATA_DIRECTORY` - Main data directory path

**Default Data Directory by Platform:**
- **Linux**: `~/.local/share/llamadeck`
- **macOS**: `~/Library/Application Support/llamadeck`
- **Windows**: `%LOCALAPPDATA%\llamadeck` or `%PROGRAMDATA%\llamadeck`

### Instance Configuration

```yaml
instances:
  port_range: [8000, 9000]      # Port range for instances (default: [8000, 9000])
  configs_dir: "instances"      # Directory for instance configs, default: data_dir/instances
  logs_dir: "logs"              # Directory for instance logs, default: data_dir/logs
  auto_create_dirs: true        # Automatically create data/config/logs directories (default: true)
  max_instances: -1             # Maximum instances (-1 = unlimited)
  max_running_instances: -1     # Maximum running instances (-1 = unlimited)
  enable_lru_eviction: true        # Enable LRU eviction for idle instances
  default_idle_timeout: 30         # Default idle timeout in minutes (0 = no timeout)
  default_auto_restart: true       # Default auto-restart setting
  default_max_restarts: 3          # Default maximum restart attempts
  default_restart_delay: 5         # Default restart delay in seconds
  default_on_demand_start: true    # Default on-demand start setting
  on_demand_start_timeout: 120     # Default on-demand start timeout in seconds
  timeout_check_interval: 5        # Default instance timeout check interval in minutes
  group_limits: {}                 # Per-group running instance limits (e.g., {large: 1, small: 3})
  log_rotation_enabled: true    # Enable log rotation (default: true)
  log_rotation_max_size: 100    # Max log file size in MB before rotation (default: 100)
  log_rotation_compress: false  # Compress rotated log files (default: false)
```

**Environment Variables:**
- `LLAMADECK_INSTANCE_PORT_RANGE` - Port range (format: "8000-9000" or "8000,9000")
- `LLAMADECK_INSTANCES_DIR` - Instance configs directory path
- `LLAMADECK_LOGS_DIR` - Log directory path
- `LLAMADECK_AUTO_CREATE_DATA_DIR` - Auto-create data/config/logs directories (true/false)
- `LLAMADECK_MAX_INSTANCES` - Maximum number of instances  
- `LLAMADECK_MAX_RUNNING_INSTANCES` - Maximum number of running instances
- `LLAMADECK_ENABLE_LRU_EVICTION` - Enable LRU eviction for idle instances
- `LLAMADECK_DEFAULT_IDLE_TIMEOUT` - Default idle timeout in minutes (0 = no timeout)
- `LLAMADECK_DEFAULT_AUTO_RESTART` - Default auto-restart setting (true/false)  
- `LLAMADECK_DEFAULT_MAX_RESTARTS` - Default maximum restarts  
- `LLAMADECK_DEFAULT_RESTART_DELAY` - Default restart delay in seconds  
- `LLAMADECK_DEFAULT_ON_DEMAND_START` - Default on-demand start setting (true/false)  
- `LLAMADECK_ON_DEMAND_START_TIMEOUT` - Default on-demand start timeout in seconds
- `LLAMADECK_TIMEOUT_CHECK_INTERVAL` - Default instance timeout check interval in minutes
- `LLAMADECK_GROUP_LIMITS` - Per-group running instance limits (format: "group1=2,group2=1")
- `LLAMADECK_LOG_ROTATION_ENABLED` - Enable log rotation (true/false)
- `LLAMADECK_LOG_ROTATION_MAX_SIZE` - Max log file size in MB
- `LLAMADECK_LOG_ROTATION_COMPRESS` - Compress rotated logs (true/false)

### Database Configuration

```yaml
database:
  path: "llamadeck.db"              # Database file path, default: data_dir/llamadeck.db
  max_open_connections: 25         # Maximum open database connections (default: 25)
  max_idle_connections: 5          # Maximum idle database connections (default: 5)
  connection_max_lifetime: 5m      # Connection max lifetime (default: 5m)
```

**Environment Variables:**
- `LLAMADECK_DATABASE_PATH` - Database file path (relative to data_dir or absolute)
- `LLAMADECK_DATABASE_MAX_OPEN_CONNECTIONS` - Maximum open database connections
- `LLAMADECK_DATABASE_MAX_IDLE_CONNECTIONS` - Maximum idle database connections
- `LLAMADECK_DATABASE_CONN_MAX_LIFETIME` - Connection max lifetime (e.g., "5m", "1h")

### Authentication Configuration

llamadeck supports three authentication mechanisms:

- **Management API Keys**: For accessing the web UI and management API (creating/managing instances). These can be configured in the config file or via environment variables.
- **Inference API Keys**: For accessing the OpenAI-compatible inference endpoints. These are managed via the web UI (Settings → API Keys) and stored in the database.
- **OIDC / SSO Sessions**: Optional web UI login through an OpenID Connect provider, configured via the `auth.oidc` block below.

See [Authentication](authentication.md) for inference key permissions and SSO setup.

```yaml
auth:
  require_inference_auth: true           # Require API key for OpenAI endpoints (default: true)
  require_management_auth: true          # Require API key for management endpoints (default: true)
  management_keys: []                    # List of valid management API keys
```

Inference API keys are created and managed via the web UI or management API

**Environment Variables:**
- `LLAMADECK_REQUIRE_INFERENCE_AUTH` - Require auth for OpenAI endpoints (true/false)
- `LLAMADECK_REQUIRE_MANAGEMENT_AUTH` - Require auth for management endpoints (true/false)
- `LLAMADECK_MANAGEMENT_KEYS` - Comma-separated management API keys

#### OIDC / SSO Login

The web UI can authenticate users through an OpenID Connect provider (Authentik, Keycloak, Authelia, etc.) instead of a management API key. OIDC login is enabled when both `issuer_url` and `client_id` are set; management API keys keep working alongside it, and inference endpoints always use API keys.

```yaml
auth:
  oidc:
    issuer_url: "https://authentik.example.com/application/o/llamadeck/"  # IdP issuer URL
    client_id: "llamadeck"                  # OAuth2 client ID registered with the IdP
    client_secret: "your-client-secret"    # OAuth2 client secret
    # redirect_url: ""                     # Optional; derived from the request when empty
    # scopes: [openid, profile, email]     # Default scopes
    # allowed_groups: []                   # Restrict login to these groups (empty = all IdP users)
    # groups_claim: "groups"               # ID-token claim carrying group names
    # session_ttl: 12h                     # Session lifetime (default: 12h)
    # secure_cookie: true                  # Set false only for plain-HTTP LAN deployments
```

**Environment Variables:**
- `LLAMADECK_AUTH_OIDC_ISSUER_URL` - IdP issuer URL
- `LLAMADECK_AUTH_OIDC_CLIENT_ID` - OAuth2 client ID
- `LLAMADECK_AUTH_OIDC_CLIENT_SECRET` - OAuth2 client secret
- `LLAMADECK_AUTH_OIDC_REDIRECT_URL` - Redirect URL override
- `LLAMADECK_AUTH_OIDC_SCOPES` - Comma-separated scopes
- `LLAMADECK_AUTH_OIDC_ALLOWED_GROUPS` - Comma-separated groups allowed to log in
- `LLAMADECK_AUTH_OIDC_GROUPS_CLAIM` - Name of the ID-token claim carrying groups (default: `groups`)
- `LLAMADECK_AUTH_OIDC_SESSION_TTL` - Session lifetime (Go duration, e.g. `12h`)
- `LLAMADECK_AUTH_OIDC_SECURE_COOKIE` - Mark the session cookie Secure (true/false)

See [Authentication](authentication.md#oidc-sso-login) for IdP client setup.

### Remote Node Configuration

llamadeck supports remote node deployments. Configure remote nodes to deploy instances on remote hosts and manage them centrally.

```yaml
local_node: "main"               # Name of the local node (default: "main")
nodes:                           # Node configuration map
  main:                          # Local node (empty address means local)
    address: ""                  # Not used for local node
    api_key: ""                  # Not used for local node
  worker1:                       # Remote worker node
    address: "http://192.168.1.10:8080"
    api_key: "worker1-api-key"   # Management API key for authentication
```

**Node Configuration Fields:**
- `local_node`: Specifies which node in the `nodes` map represents the local node. Must match exactly what other nodes call this node.
- `nodes`: Map of node configurations
  - `address`: HTTP/HTTPS URL of the remote node (empty for local node)
  - `api_key`: Management API key for authenticating with the remote node

**Environment Variables:**
- `LLAMADECK_LOCAL_NODE` - Name of the local node
