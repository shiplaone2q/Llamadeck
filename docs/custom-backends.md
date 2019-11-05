# Custom Backends

Custom backends provide support for inference servers without native llamadeck integration. Any server can be managed (launched, health-checked, proxied, and evicted) by defining it as a named entry in the config file.

## How It Works

A custom backend is a named entry under `backends.custom` in the config file; instances reference the entry by name, and llamadeck runs the configured command with the instance's arguments. The available fields are documented in [Backend Configuration](configuration.md#custom-backends).
