# Authentication

Llamadeck controls access with three mechanisms:

- **Management API Key** — authenticates the web UI and management API (creating/stopping instances, managing keys). Configured in the config file or via environment variables; see the [Configuration](configuration.md#authentication-configuration) guide.
- **Inference API Key** — authenticates OpenAI-compatible inference requests (`/v1/chat/completions`, `/v1/completions`, etc.). Created and managed via the web UI or management API and stored in the database.
- **OIDC / SSO Session** — optional browser login for the web UI through an OpenID Connect provider; see [OIDC / SSO Login](#oidc-sso-login) below.

!!! warning "Management access is host-administrator access"
    Anyone who can access the management dashboard or management API can configure commands, arguments, environment variables, containers, custom backends, and model downloads. Grant management API keys and OIDC access only to users trusted to administer the host. Use `allowed_groups` to restrict which identity-provider users receive this access.

The rest of this page covers **inference API keys** and their per-instance permissions, followed by **SSO setup**.

## Inference API Keys

### Permission Modes

When you create an inference key, you choose a permission mode:

- **Full Access** (`allow_all`): The key can use every instance with full lifecycle rights. This is the simplest option.
- **Per-Instance Access** (`per_instance`): You explicitly grant access to specific instances and choose an access level for each.

Management API keys bypass all permission checks entirely.

### Access Levels

For each instance granted to a per-instance key, you pick one of three access levels. A key can always **use a running instance** — the level only affects what happens when the instance is **stopped** at request time.

| Level | What the key can do |
|---|---|
| **Use running only** | Send requests only while the instance is already running. Never starts it. |
| **Can start on demand** | Auto-start the instance when needed, but won't disrupt others — fails if there's no free capacity. |
| **Can start and evict others** | Auto-start the instance and may evict other instances to free up room. |

!!! note "What &quot;evict&quot; means"
    The *Can start and evict others* level allows the key to **evict other instances** to make room for the one it's starting. It does **not** protect this instance from being evicted when a different key starts something else. Eviction always targets the least-recently-used instance.

Both flags default to enabled, so omitting them is equivalent to *Can start and evict others* (backward compatible with keys created before this feature).

### Behavior Reference

How a request is handled depends on the instance's state and the key's access level for it:

| Instance state | Access level | Result |
|---|---|---|
| Running | any | Request proxied normally |
| Stopped, free capacity | Use running only | `503 instance_not_running` |
| Stopped, free capacity | Can start (either) | Instance started, request proxied |
| Stopped, at capacity | Use running only | `503 instance_not_running` |
| Stopped, at capacity | Can start on demand | `503 max_instances_reached` |
| Stopped, at capacity | Can start and evict | Evicts LRU instance, starts, request proxied |

