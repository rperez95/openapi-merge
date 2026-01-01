# Quick Start

Get up and running with openapi-merge in 5 minutes.

## Step 1: Create Your OpenAPI Files

Let's say you have two microservices with their own OpenAPI specs:

=== "users-api.json"

    ```json
    {
      "openapi": "3.0.0",
      "info": {
        "title": "Users API",
        "version": "1.0.0"
      },
      "paths": {
        "/users": {
          "get": {
            "summary": "List users",
            "tags": ["Users"],
            "responses": {
              "200": { "description": "Success" }
            }
          }
        }
      }
    }
    ```

=== "orders-api.json"

    ```json
    {
      "openapi": "3.0.0",
      "info": {
        "title": "Orders API",
        "version": "1.0.0"
      },
      "paths": {
        "/orders": {
          "get": {
            "summary": "List orders",
            "tags": ["Orders"],
            "responses": {
              "200": { "description": "Success" }
            }
          }
        }
      }
    }
    ```

## Step 2: Create Configuration File

Create a `merge-config.yaml` file:

```yaml title="merge-config.yaml"
# Output API metadata
info:
  title: "Platform API"
  description: "Unified API for all platform services"
  version: "1.0.0"

# Server configuration
servers:
  - url: "https://api.example.com"
    description: "Production"

# Input files to merge
inputs:
  - inputFile: users-api.json
    pathModification:
      prepend: "/users-service"

  - inputFile: orders-api.json
    pathModification:
      prepend: "/orders-service"

# Output file
output: platform-api.json
```

## Step 3: Run the Merge

```bash
openapi-merge merge --config merge-config.yaml
```

Output:
```
Successfully merged 2 specifications into platform-api.json
```

## Step 4: Verify the Result

The merged `platform-api.json` will contain:

```json
{
  "openapi": "3.0.3",
  "info": {
    "title": "Platform API",
    "description": "Unified API for all platform services",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "https://api.example.com",
      "description": "Production"
    }
  ],
  "paths": {
    "/users-service/users": {
      "get": {
        "summary": "List users",
        "tags": ["Users"],
        "responses": {
          "200": { "description": "Success" }
        }
      }
    },
    "/orders-service/orders": {
      "get": {
        "summary": "List orders",
        "tags": ["Orders"],
        "responses": {
          "200": { "description": "Success" }
        }
      }
    }
  }
}
```

## Adding Authentication

Add security to your merged API:

```yaml title="merge-config.yaml" hl_lines="12-18 20-21"
info:
  title: "Platform API"
  version: "1.0.0"

servers:
  - url: "https://api.example.com"

# Define security schemes
securitySchemes:
  bearerAuth:
    type: http
    scheme: bearer
    bearerFormat: JWT

# Apply globally
security:
  - bearerAuth: []

inputs:
  - inputFile: users-api.json
    pathModification:
      prepend: "/users-service"

  - inputFile: orders-api.json
    pathModification:
      prepend: "/orders-service"

output: platform-api.json
```

## Using Verbose Mode

See detailed progress during merge:

```bash
openapi-merge merge --config merge-config.yaml --verbose
```

Output:
```
Starting merge with 2 input files
Output file: /path/to/platform-api.json
Processing input 1: users-api.json
Processing input 2: orders-api.json
Successfully merged 2 specifications into platform-api.json
```

## Specifying Output File

Override the config output with `-o` flag:

```bash
openapi-merge merge --config merge-config.yaml -o custom-output.yaml
```

!!! tip "Output Format"
    The output format (JSON or YAML) is automatically determined by the file extension.

## Next Steps

- [Configuration Overview](../configuration/overview.md) - Learn all configuration options
- [Path Modification](../configuration/paths.md) - Advanced path transformations
- [Security Schemes](../configuration/security.md) - Authentication configuration
- [Examples](../examples/basic.md) - More real-world examples

