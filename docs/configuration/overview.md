# Configuration Overview

openapi-merge uses a YAML or JSON configuration file to control the merge process.

## Configuration File Structure

```yaml
# API metadata override
info:
  title: "API Title"
  description: "API Description"
  version: "1.0.0"

# Server definitions
servers:
  - url: "https://api.example.com"
    description: "Production"

# Global base path (optional)
basePath: "/v1"

# Security scheme definitions
securitySchemes:
  bearerAuth:
    type: http
    scheme: bearer

# Global security requirements
security:
  - bearerAuth: []

# Tag ordering
tagOrder:
  - "Users"
  - "Orders"

# Path ordering (priority paths first)
pathsOrder:
  - "/v1/users"
  - "/v1/orders"

# Input files to merge
inputs:
  - inputFile: "path/to/spec.json"
    # ... input configuration

# Output file path
output: "merged-api.json"
```

## Top-Level Properties

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `inputs` | `[]InputConfig` | ✅ | List of input files to merge |
| `output` | `string` | ✅ | Path to save the merged file |
| `info` | `InfoConfig` | ❌ | Override API metadata |
| `servers` | `[]ServerConfig` | ❌ | Server definitions |
| `basePath` | `string` | ❌ | Global prefix for all paths |
| `securitySchemes` | `map[string]SecurityScheme` | ❌ | Security scheme definitions |
| `security` | `[]SecurityRequirement` | ❌ | Global security requirements |
| `tagOrder` | `[]string` | ❌ | Tag ordering in output |
| `pathsOrder` | `[]string` | ❌ | High-priority paths (appear first) |

## Info Configuration

Override the merged API's info section:

```yaml
info:
  title: "Platform API"
  description: |
    # Platform API
    
    Full API documentation for the platform.
    
    ## Authentication
    Use Bearer token authentication.
  version: "2.0.0"
  termsOfService: "https://example.com/terms"
  contact:
    name: "API Support"
    url: "https://example.com/support"
    email: "support@example.com"
  license:
    name: "Apache 2.0"
    url: "https://www.apache.org/licenses/LICENSE-2.0"
```

## Server Configuration

Define API servers:

```yaml
servers:
  - url: "https://api.example.com"
    description: "Production server"
  
  - url: "https://staging-api.example.com"
    description: "Staging server"
  
  - url: "http://localhost:8080"
    description: "Local development"
  
  # With variables
  - url: "https://{environment}.api.example.com"
    description: "Dynamic environment"
    variables:
      environment:
        default: "prod"
        enum: ["prod", "staging", "dev"]
        description: "Environment name"
```

## Global Base Path

Prepend a path prefix to all merged paths:

```yaml
basePath: "/api/v1"
```

!!! example "Path Transformation"
    With `basePath: "/api/v1"`:
    
    - `/users` → `/api/v1/users`
    - `/orders` → `/api/v1/orders`

## Tag and Path Ordering

Control the order of tags and paths in the output:

```yaml
# Tags appear in this order (unlisted tags are appended alphabetically)
tagOrder:
  - "Authentication"
  - "Users"
  - "Orders"
  - "Admin"

# High-priority paths appear first
pathsOrder:
  - "/api/v1/auth/login"
  - "/api/v1/users"
  - "/api/v1/orders"
```

## Output Format

The output format is determined by the file extension:

| Extension | Format |
|-----------|--------|
| `.json` | JSON |
| `.yaml`, `.yml` | YAML |

```yaml
# JSON output
output: merged-api.json

# YAML output
output: merged-api.yaml
```

!!! tip "CLI Override"
    Override the output path with the `-o` flag:
    ```bash
    openapi-merge merge --config config.yaml -o custom-output.yaml
    ```

## Environment Variables

Configuration values can reference environment variables (coming soon):

```yaml
servers:
  - url: "${API_BASE_URL}"
    description: "Production"
```

## Validation

The configuration is validated before merge:

- At least one input file is required
- Output path is required
- All input files must exist
- Input files must be valid OpenAPI 2.0 or 3.0

## Next Steps

- [Input Files Configuration](inputs.md)
- [Path Modification](paths.md)
- [Security Schemes](security.md)
- [Operation Filtering](filtering.md)

