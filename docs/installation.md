# Installation

This guide will walk you through installing Llamadeck on your system.

## Prerequisites

### Backend Dependencies

llamadeck supports multiple backends. Install at least one:

**For llama.cpp backend (all platforms):**

You need `llama-server` from [llama.cpp](https://github.com/ggml-org/llama.cpp) installed:

```bash
# Homebrew (macOS/Linux)
brew install llama.cpp
# Winget (Windows)
winget install llama.cpp
```

Or build from source - see llama.cpp docs

**For MLX backend (macOS only):**

MLX provides optimized inference on Apple Silicon. Install MLX-LM:

```bash
# Install via pip (requires Python 3.8+)
pip install mlx-lm

# Or in a virtual environment (recommended)
python -m venv mlx-env
source mlx-env/bin/activate
pip install mlx-lm
```

Note: MLX backend is only available on macOS with Apple Silicon (M1, M2, M3, etc.)

**For vLLM backend:**

vLLM provides high-throughput distributed serving for LLMs. Install vLLM:

```bash
# Install in a virtual environment
python -m venv vllm-env
source vllm-env/bin/activate
pip install vllm
```

## Installation Methods

### Option 1: Download Binary (Recommended)

Download the latest release from the [GitHub releases page](https://github.com/lordmathis/llamadeck/releases):

```bash
# Linux/macOS - Get latest version and download
LATEST_VERSION=$(curl -s https://api.github.com/repos/lordmathis/llamadeck/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
curl -L https://github.com/lordmathis/llamadeck/releases/download/${LATEST_VERSION}/llamadeck-${LATEST_VERSION}-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m).tar.gz | tar -xz
sudo mv llamadeck /usr/local/bin/

# Or download manually from:
# https://github.com/lordmathis/llamadeck/releases/latest

# Windows - Download from releases page
```

### Option 2: Docker

llamadeck provides Dockerfiles for creating Docker images with backends pre-installed. The resulting images include the latest llamadeck release with the respective backend.

**Available Dockerfiles (CUDA):**
- **llamadeck with llama.cpp CUDA**: `docker/Dockerfile.llamacpp` (based on `ghcr.io/ggml-org/llama.cpp:server-cuda`)
- **llamadeck with vLLM CUDA**: `docker/Dockerfile.vllm` (based on `vllm/vllm-openai:latest`)
- **llamadeck built from source**: `docker/Dockerfile.source` (multi-stage build with webui)

**Note:** These Dockerfiles are configured for CUDA. For other platforms (CPU, ROCm, Vulkan, etc.), adapt the base image. For llama.cpp, see available tags at [llama.cpp Docker docs](https://github.com/ggml-org/llama.cpp/blob/master/docs/docker.md). For vLLM, check [vLLM docs](https://docs.vllm.ai/en/v0.6.5/serving/deploying_with_docker.html).

**Using Docker Compose**

```bash
# Clone the repository
git clone https://github.com/lordmathis/llamadeck.git
cd llamadeck

# Create directories for data and models
mkdir -p data/llamacpp data/vllm models

# Start llamadeck with llama.cpp backend
docker-compose -f docker/docker-compose.yml up llamadeck-llamacpp -d

# Or start llamadeck with vLLM backend
docker-compose -f docker/docker-compose.yml up llamadeck-vllm -d
```

Access the dashboard at:
- llamadeck with llama.cpp: http://localhost:8080
- llamadeck with vLLM: http://localhost:8081

**Using Docker Build and Run**

1. llamadeck with llama.cpp CUDA:
```bash
docker build -f docker/Dockerfile.llamacpp -t llamadeck:llamacpp-cuda .
docker run -d \
  --name llamadeck-llamacpp \
  --gpus all \
  -p 8080:8080 \
  -v ~/.cache/llama.cpp:/root/.cache/llama.cpp \
  llamadeck:llamacpp-cuda
```

2. llamadeck with vLLM CUDA:
```bash
docker build -f docker/Dockerfile.vllm -t llamadeck:vllm-cuda .
docker run -d \
  --name llamadeck-vllm \
  --gpus all \
  -p 8080:8080 \
  -v ~/.cache/huggingface:/root/.cache/huggingface \
  llamadeck:vllm-cuda
```

3. llamadeck built from source:
```bash
docker build -f docker/Dockerfile.source -t llamadeck:source .
docker run -d \
  --name llamadeck \
  -p 8080:8080 \
  llamadeck:source
```

### Option 3: Build from Source

Requirements:
- Go 1.24 or later
- Node.js 22 or later
- Git

If you prefer to build from source:

```bash
# Clone the repository
git clone https://github.com/lordmathis/llamadeck.git
cd llamadeck

# Build the web UI
cd webui && npm ci && npm run build && cd ..

# Build the application
go build -o llamadeck ./cmd/server
```

## Remote Node Installation

For deployments with remote nodes:
- Install llamadeck on each node using any of the methods above
- Configure API keys for authentication between nodes
- Ensure node names are consistent across all configurations

## Verification

Verify your installation by checking the version:

```bash
llamadeck --version
```

## Next Steps

Now that Llamadeck is installed, continue to the [Quick Start](quick-start.md) guide to get your first instance running!

For remote node deployments, see the [Configuration Guide](configuration.md) for node setup instructions.
