# Custom Backends

Custom backends provide support for inference servers without native llamadeck integration. Any server can be managed (launched, health-checked, proxied, and evicted) by defining it as a named entry in the config file.

## How It Works

A custom backend is a named entry under `backends.custom` in the config file; instances reference the entry by name, and llamadeck runs the configured command with the instance's arguments. The available fields are documented in [Backend Configuration](configuration.md#custom-backends).

The managed server must meet three requirements:

- It starts an HTTP server that listens on a TCP port.
- It serves inference endpoints under `/v1/`. All `POST /v1/*` requests are proxied to the instance, so OpenAI-style (`/v1/chat/completions`, `/v1/responses`) and Anthropic-style (`/v1/messages`) conventions all work.
- It exposes a health endpoint that returns a 2xx response when the server is ready (see [Health Checks](#health-checks)).
