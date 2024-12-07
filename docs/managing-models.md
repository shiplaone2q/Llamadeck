# Managing Models

Llamadeck provides built-in model management capabilities for downloading models directly from HuggingFace without manually managing files. It supports both **GGUF** models (for llama.cpp) and **Safetensors** models (for vLLM and MLX backends).

## Overview

The model downloader:
- Downloads models directly from HuggingFace repositories in two formats: **GGUF** and **Safetensors**
- Caches models in a local directory using HuggingFace Hub-compatible layout for reuse
- Supports split files (models split into multiple parts) for both formats
- Downloads additional files like mmproj (for multimodal GGUF models) and preset.ini
- Downloads tokenizer, config, and all weight files for safetensors models
- Tracks download progress and manages jobs
- Works with both public and private HuggingFace repositories
- Supports resume of interrupted downloads

## Cache Directory

Downloaded models are cached in the following default locations:

- **Linux**: `~/.cache/llama.cpp/` (or `$XDG_CACHE_HOME/llama.cpp/` if XDG_CACHE_HOME is set)
- **macOS**: `~/Library/Caches/llama.cpp/`
- **Windows**: `%LOCALAPPDATA%\llama.cpp\cache\`

The cache directory can be customized using:
- Configuration file: `backends.llama_cpp.cache_dir`
- Environment variable: `LLAMADECK_LLAMACPP_CACHE_DIR`
- Environment variable: `LLAMA_CACHE` (standard llama.cpp convention)

## Supported Formats

Llamadeck supports downloading models in two formats:

### GGUF (for llama.cpp)

GGUF is the native format for llama.cpp. GGUF models are pre-quantized and contain everything needed for inference in a single file (or a set of split files).

- Use with the **llama.cpp** backend
- Specify a quantization tag to select a specific variant (e.g., `Q4_K_M`, `Q8_0`)
- Example repos: `bartowski/Llama-3.2-3B-Instruct-GGUF`, `TheBloke/Mistral-7B-Instruct-v0.2-GGUF`

### Safetensors (for vLLM and MLX)

Safetensors is the standard model format used by HuggingFace. Downloads include all weight files, tokenizer files, config files, and other mandatory files needed to load the model.

- Use with **vLLM** or **MLX** backends
- Downloads the full model directory (all safetensors weight shards, tokenizer, config, etc.)
- Falls back to `.bin` files if no safetensors files are found in the repo
- Example repos: `meta-llama/Llama-3.2-3B`, `mistralai/Mistral-7B-Instruct-v0.3`

## Downloading Models

### Using the Web UI

1. Open the Llamadeck web interface
2. Navigate to the **Models** tab
3. Click **Download Model**
4. Select the **Format**: **GGUF** or **Safetensors**
5. Enter the model identifier:
   - **GGUF**: `org/model-name:tag` (e.g., `bartowski/Llama-3.2-3B-Instruct-GGUF:Q4_K_M`)
   - **Safetensors**: `org/model-name` (e.g., `meta-llama/Llama-3.2-3B`)
   - The `:tag` part is optional - if omitted, the default branch is used
6. Click **Start Download**
7. Monitor the download progress in the UI

The Models tab shows:
- **Downloaded Bytes**: Progress of current file download
- **Current File**: Name of the file being downloaded
- **Status**: Current job status (queued, downloading, completed, failed, cancelled)

### Using the API

Download a GGUF model:

```bash
curl -X POST http://localhost:8080/api/v1/models/download \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_MANAGEMENT_KEY" \
  -d '{
    "repo": "bartowski/Llama-3.2-3B-Instruct-GGUF:Q4_K_M",
    "format": "gguf"
  }'
```

Download a safetensors model:

```bash
curl -X POST http://localhost:8080/api/v1/models/download \
  -H "Content-Type: application/json" \
