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
