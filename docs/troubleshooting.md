# Troubleshooting

Issues specific to Llamadeck deployment and operation.

## Configuration Issues

### Invalid Configuration

**Problem:** Invalid configuration preventing startup

**Solutions:**
1. Use minimal configuration:
   ```yaml
   server:
     host: "0.0.0.0"
     port: 8080
   instances:
     port_range: [8000, 9000]
   ```

2. Check data directory permissions:
   ```bash
   # Ensure data directory is writable (default: ~/.local/share/llamadeck)
   mkdir -p ~/.local/share/llamadeck/{instances,logs}
   ```

## Instance Management Issues

### Instance Fails to Start

**Problem:** Instance fails to start or immediately stops

**Solutions:**

1. **Check instance logs** to see the actual error:
   ```bash
   curl http://localhost:8080/api/v1/instances/{name}/logs
   # Or check log files directly
   tail -f ~/.local/share/llamadeck/logs/{instance-name}.log
   ```

2. **Verify backend is installed:**  
     - **llama.cpp**: Ensure `llama-server` is in PATH
     - **MLX**: Ensure `mlx-lm` Python package is installed
     - **vLLM**: Ensure `vllm` Python package is installed

3. **Check model path and format:**
     - Use absolute paths to model files
     - Verify model format matches backend (GGUF for llama.cpp, etc.)

4. **Verify backend command configuration:**
     - Check that the backend `command` is correctly configured in the global config
     - For virtual environments, specify the full path to the command (e.g., `/path/to/venv/bin/mlx_lm.server`)
     - See the [Configuration Guide](configuration.md) for backend configuration details
