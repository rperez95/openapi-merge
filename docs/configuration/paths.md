# Path Modification

Transform API paths during the merge process.

## Path Modification Options

```yaml
inputs:
  - inputFile: api.json
    pathModification:
      stripStart: "/v1"      # Remove from beginning
      prepend: "/service/v1" # Add to beginning
```

## Strip Start

Remove a prefix from all paths:

```yaml
pathModification:
  stripStart: "/v1"
```

| Original Path | After stripStart |
|---------------|------------------|
| `/v1/users` | `/users` |
| `/v1/users/{id}` | `/users/{id}` |
| `/v1/orders` | `/orders` |
| `/health` | `/health` (unchanged) |

!!! note
    Only removes the prefix if the path **starts** with the specified string.

## Prepend

Add a prefix to all paths:

```yaml
pathModification:
  prepend: "/users-service"
```

| Original Path | After prepend |
|---------------|---------------|
| `/users` | `/users-service/users` |
| `/users/{id}` | `/users-service/users/{id}` |

## Combined Transformation

Use both options together for full control:

```yaml
pathModification:
  stripStart: "/api/v1"
  prepend: "/gateway/users/v2"
```

| Original Path | Step 1 (strip) | Step 2 (prepend) |
|---------------|----------------|------------------|
| `/api/v1/users` | `/users` | `/gateway/users/v2/users` |
| `/api/v1/profile` | `/profile` | `/gateway/users/v2/profile` |

## Global Base Path

In addition to per-input path modification, you can set a global base path:

```yaml
# Global configuration
basePath: "/api"

inputs:
  - inputFile: users.json
    pathModification:
      prepend: "/users-service/v1"
```

The global `basePath` is applied **after** all input-specific modifications:

| Step | Path |
|------|------|
| Original | `/users` |
| After input prepend | `/users-service/v1/users` |
| After global basePath | `/api/users-service/v1/users` |

## Real-World Examples

### Microservices Gateway

```yaml
basePath: "/api/v1"

inputs:
  - inputFile: apis/users.json
    pathModification:
      stripStart: "/v1"
      prepend: "/users"
  
  - inputFile: apis/orders.json
    pathModification:
      stripStart: "/v1"
      prepend: "/orders"
  
  - inputFile: apis/payments.json
    pathModification:
      stripStart: "/api"
      prepend: "/payments"
```

Result:
```
/api/v1/users/...
/api/v1/orders/...
/api/v1/payments/...
```

### Version Migration

Moving from v1 to v2 API structure:

```yaml
inputs:
  - inputFile: legacy-v1-api.json
    pathModification:
      stripStart: "/v1"
      prepend: "/v2"
```

### Service Namespacing

Isolate services under unique namespaces:

```yaml
inputs:
  - inputFile: team-a/api.json
    pathModification:
      prepend: "/team-a"
  
  - inputFile: team-b/api.json
    pathModification:
      prepend: "/team-b"
```

### Swagger 2.0 basePath Handling

Swagger 2.0 files have a `basePath` field. When converted to OpenAPI 3.0, this is incorporated into the paths. Use `stripStart` to normalize:

```json title="legacy.json (Swagger 2.0)"
{
  "swagger": "2.0",
  "basePath": "/v1",
  "paths": {
    "/users": { ... }
  }
}
```

```yaml title="config.yaml"
inputs:
  - inputFile: legacy.json
    pathModification:
      stripStart: "/v1"
      prepend: "/legacy-service/v1"
```

## Path Ordering

Control the order of paths in the output:

```yaml
pathsOrder:
  - "/api/v1/auth/login"
  - "/api/v1/auth/logout"
  - "/api/v1/users"
  - "/api/v1/orders"
```

!!! tip
    Paths not in `pathsOrder` are sorted alphabetically after the priority paths.

## Edge Cases

### Trailing Slashes

```yaml
pathModification:
  stripStart: "/v1/"  # Includes trailing slash
  prepend: "/api/"    # Includes trailing slash
```

Be careful with trailing slashes to avoid double slashes (`//`).

### Empty Paths

If `stripStart` results in an empty path, it becomes `/`:

```yaml
pathModification:
  stripStart: "/health"
```

| Original | Result |
|----------|--------|
| `/health` | `/` |
| `/health/live` | `/live` |

### Path Starts With

`stripStart` only works if the path **starts** with the specified string:

```yaml
pathModification:
  stripStart: "/users"
```

| Original | Result |
|----------|--------|
| `/users/list` | `/list` ✅ |
| `/api/users` | `/api/users` (unchanged) ❌ |

## Next Steps

- [Operation Filtering](filtering.md) - Include/exclude specific operations
- [Security Schemes](security.md) - Authentication configuration

