# Operation Filtering

Filter which operations are included in the merged output.

## Overview

You can include or exclude operations based on:

- **Tags** - Filter by operation tags
- **Paths** - Filter by path patterns (with glob support)
- **Parameters** - Add or remove parameters from operations

```yaml
inputs:
  - inputFile: api.json
    operationSelection:
      includeTags: ["Public"]
      excludeTags: ["Internal", "Admin"]
      includePaths:
        - path: "/api/*"
          method: "GET"
      excludePaths:
        - path: "/internal/*"
```

## Tag Filtering

### Include Tags

Only include operations with specified tags:

```yaml
operationSelection:
  includeTags:
    - "Users"
    - "Orders"
    - "Products"
```

!!! note
    If `includeTags` is empty or not specified, all tags are included.

### Exclude Tags

Exclude operations with specified tags:

```yaml
operationSelection:
  excludeTags:
    - "Internal"
    - "Admin"
    - "Deprecated"
```

### Combined Tag Filtering

```yaml
operationSelection:
  includeTags: ["Users", "Orders"]  # Must have one of these
  excludeTags: ["Internal"]          # But not this
```

An operation is included if:
1. It has at least one tag from `includeTags` (if specified)
2. It does NOT have any tag from `excludeTags`

## Path Filtering

### Include Paths

Whitelist specific paths:

```yaml
operationSelection:
  includePaths:
    - path: "/users"
    - path: "/users/*"
    - path: "/orders"
```

### Exclude Paths

Blacklist specific paths:

```yaml
operationSelection:
  excludePaths:
    - path: "/internal/*"
    - path: "/admin/*"
    - path: "/debug/*"
```

### Filter by Method

Specify HTTP method for finer control:

```yaml
operationSelection:
  includePaths:
    - path: "/users"
      method: "GET"     # Only GET /users
    - path: "/users/*"
      method: "GET"     # Only GET /users/{id}
  
  excludePaths:
    - path: "/users"
      method: "DELETE"  # Exclude DELETE /users
```

!!! tip
    If `method` is not specified, the filter applies to all HTTP methods.

## Glob Pattern Support

Path filters support glob patterns:

| Pattern | Matches | Does Not Match |
|---------|---------|----------------|
| `/users` | `/users` | `/users/123` |
| `/users/*` | `/users/123`, `/users/profile` | `/users/123/orders` |
| `/users/**` | `/users/123`, `/users/123/orders` | - |
| `/api/*/items` | `/api/v1/items`, `/api/v2/items` | `/api/items` |

### Examples

```yaml
operationSelection:
  includePaths:
    # Exact match
    - path: "/health"
    
    # Single segment wildcard
    - path: "/api/*/users"
    
    # Multi-segment wildcard
    - path: "/public/**"
  
  excludePaths:
    # All internal endpoints
    - path: "/internal/**"
    
    # All admin operations
    - path: "/**/admin/**"
```

## Parameter Filtering

### Include Extra Parameters

Inject parameters into all operations from an input:

```yaml
includeExtraParameters:
  - name: "X-Request-ID"
    in: header
    description: "Request tracking ID"
    required: false
    schema:
      type: string
  
  - name: "X-Tenant-ID"
    in: header
    description: "Tenant identifier"
    required: true
    schema:
      type: string
```

### Parameter Properties

| Property | Type | Description |
|----------|------|-------------|
| `name` | `string` | Parameter name |
| `in` | `string` | Location: `header`, `query`, `path`, `cookie` |
| `description` | `string` | Parameter description |
| `required` | `boolean` | Whether the parameter is required |
| `schema` | `object` | JSON Schema for the parameter |

### Exclude Parameters

Remove parameters matching filters:

```yaml
excludeParameters:
  - name: "X-Internal-Auth"
    in: header
  
  - name: "debug"
    in: query
  
  - name: "X-User-Id"
    # in: not specified = matches any location
```

!!! info "Matching Logic"
    - `name` is required and must match exactly
    - `in` is optional; if not specified, matches any location

## Complete Example

```yaml title="config.yaml"
inputs:
  - inputFile: apis/users.json
    pathModification:
      prepend: "/users-service/v1"
    
    operationSelection:
      # Only include public and user-facing tags
      includeTags:
        - "Users"
        - "Profile"
        - "Authentication"
      
      # Exclude internal and admin tags
      excludeTags:
        - "Internal"
        - "Admin"
        - "Debug"
      
      # Only include these paths
      includePaths:
        - path: "/users"
        - path: "/users/*"
        - path: "/auth/**"
      
      # Exclude these paths
      excludePaths:
        - path: "/users/*/admin"
        - path: "/internal/**"
    
    # Add tracking headers
    includeExtraParameters:
      - name: "X-Request-ID"
        in: header
        description: "Request tracking ID"
        required: false
        schema:
          type: string
          format: uuid
    
    # Remove internal headers
    excludeParameters:
      - name: "X-Internal-Auth"
        in: header
      - name: "X-Debug-Mode"
        in: header
```

## Use Cases

### Public API Only

Include only public-facing operations:

```yaml
operationSelection:
  includeTags: ["Public"]
  excludeTags: ["Internal", "Admin", "Deprecated"]
  excludePaths:
    - path: "/internal/**"
    - path: "/admin/**"
```

### Read-Only API

Include only GET operations:

```yaml
operationSelection:
  includePaths:
    - path: "/**"
      method: "GET"
```

### Exclude Health Checks

Remove health and readiness endpoints:

```yaml
operationSelection:
  excludePaths:
    - path: "/health"
    - path: "/health/*"
    - path: "/ready"
    - path: "/live"
```

### Service Mesh Headers

Add service mesh tracing headers:

```yaml
includeExtraParameters:
  - name: "X-Request-ID"
    in: header
    description: "Unique request identifier"
    schema:
      type: string
  
  - name: "X-B3-TraceId"
    in: header
    description: "Zipkin trace ID"
    schema:
      type: string
  
  - name: "X-B3-SpanId"
    in: header
    description: "Zipkin span ID"
    schema:
      type: string
```

## Next Steps

- [Security Schemes](security.md) - Authentication configuration
- [Examples](../examples/advanced.md) - Advanced configuration examples

