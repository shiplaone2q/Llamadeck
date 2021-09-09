# Quick Start

This guide will help you get Llamadeck up and running in just a few minutes.

**Before you begin:** Ensure you have at least one backend installed (llama.cpp, MLX, or vLLM). See the [Installation Guide](installation.md#prerequisites) for backend setup.

## Core Concepts

Before you start, let's clarify a few key terms:

- **Instance**: A running backend server that serves a specific model. Each instance has a unique name and runs independently.
- **Backend**: The inference engine that actually runs the model (llama.cpp, MLX, or vLLM). You need at least one backend installed before creating instances.
- **Node**: In multi-machine setups, a node represents one machine. Most users will just use the default "main" node for single-machine deployments.
- **Proxy Architecture**: Llamadeck acts as a proxy in front of your instances. You make requests to llamadeck (e.g., `http://localhost:8080/v1/chat/completions`), and it routes them to the appropriate backend instance. This means you don't need to track individual instance ports or endpoints.

## Authentication

Llamadeck uses two types of API keys:

- **Management API Key**: Used to authenticate with the Llamadeck management API and web UI. If not configured, one is auto-generated at startup and printed to the terminal.
- **Inference API Key**: Used to authenticate requests to the OpenAI-compatible endpoints (`/v1/chat/completions`, `/v1/completions`, etc.). These are created and managed via the web UI.

By default, authentication is required for both management and inference endpoints. You can configure custom management keys or disable authentication in the [Configuration](configuration.md) guide.

## Start Llamadeck

Start the Llamadeck server:

```bash
llamadeck
```

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  MANAGEMENT AUTHENTICATION REQUIRED
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔑  Generated Management API Key:

    sk-management-...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  IMPORTANT
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
• This key is auto-generated and will change on restart
• For production, add explicit management_keys to your configuration
• Copy this key before it disappears from the terminal
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Llamadeck server listening on 0.0.0.0:8080
```

Copy the **Management API Key** from the terminal - you'll need it to access the web UI.

By default, Llamadeck will start on `http://localhost:8080`.

## Access the Web UI

Open your web browser and navigate to:

```
http://localhost:8080
```

Login with the management API key from the terminal output.

You should see the Llamadeck web interface.

## Create Your First Instance

