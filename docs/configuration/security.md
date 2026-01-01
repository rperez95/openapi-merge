# Security Schemes

Configure OAS3-compliant authentication for your merged API.

## Overview

Security in OpenAPI 3.0 consists of two parts:

1. **Security Schemes** - Define authentication methods (`components.securitySchemes`)
2. **Security Requirements** - Apply schemes to operations (`security`)

```yaml
# 1. Define security schemes
securitySchemes:
  bearerAuth:
    type: http
    scheme: bearer
    bearerFormat: JWT

# 2. Apply globally
security:
  - bearerAuth: []
```

## Supported Security Types

| Type | Use Case |
|------|----------|
| `apiKey` | API keys in header, query, or cookie |
| `http` | HTTP Basic or Bearer authentication |
| `oauth2` | OAuth 2.0 flows |
| `openIdConnect` | OpenID Connect Discovery |

## API Key Authentication

### Header API Key

```yaml
securitySchemes:
  apiKeyAuth:
    type: apiKey
    in: header
    name: X-API-Key
    description: API key for authentication

security:
  - apiKeyAuth: []
```

### Query Parameter API Key

```yaml
securitySchemes:
  apiKeyQuery:
    type: apiKey
    in: query
    name: api_key
    description: API key as query parameter
```

### Cookie Authentication

```yaml
securitySchemes:
  cookieAuth:
    type: apiKey
    in: cookie
    name: session_id
    description: Session cookie
```

!!! example "Ory Kratos Session"
    ```yaml
    securitySchemes:
      kratosSession:
        type: apiKey
        in: cookie
        name: ory_kratos_session
        description: Ory Kratos session cookie
    
    security:
      - kratosSession: []
    ```

## HTTP Authentication

### Bearer Token (JWT)

```yaml
securitySchemes:
  bearerAuth:
    type: http
    scheme: bearer
    bearerFormat: JWT
    description: JWT Bearer token authentication
```

The `bearerFormat` is optional and used for documentation purposes.

### Basic Authentication

```yaml
securitySchemes:
  basicAuth:
    type: http
    scheme: basic
    description: HTTP Basic authentication
```

## OAuth 2.0

### Authorization Code Flow

```yaml
securitySchemes:
  oauth2Auth:
    type: oauth2
    description: OAuth 2.0 Authorization Code
    flows:
      authorizationCode:
        authorizationUrl: https://auth.example.com/oauth/authorize
        tokenUrl: https://auth.example.com/oauth/token
        refreshUrl: https://auth.example.com/oauth/refresh
        scopes:
          read: Read access
          write: Write access
          admin: Admin access
```

### Implicit Flow

```yaml
securitySchemes:
  oauth2Implicit:
    type: oauth2
    flows:
      implicit:
        authorizationUrl: https://auth.example.com/oauth/authorize
        scopes:
          read: Read access
          write: Write access
```

### Client Credentials Flow

```yaml
securitySchemes:
  oauth2Client:
    type: oauth2
    flows:
      clientCredentials:
        tokenUrl: https://auth.example.com/oauth/token
        scopes:
          api: API access
```

### Password Flow

```yaml
securitySchemes:
  oauth2Password:
    type: oauth2
    flows:
      password:
        tokenUrl: https://auth.example.com/oauth/token
        scopes:
          read: Read access
          write: Write access
```

## OpenID Connect

```yaml
securitySchemes:
  oidcAuth:
    type: openIdConnect
    openIdConnectUrl: https://auth.example.com/.well-known/openid-configuration
    description: OpenID Connect authentication
```

## Applying Security

### Global Security

Apply to all operations:

```yaml
security:
  - bearerAuth: []
```

### Multiple Options (OR)

Allow any of these authentication methods:

```yaml
security:
  - bearerAuth: []    # Option 1: Bearer token
  - apiKeyAuth: []    # Option 2: API key
  - cookieAuth: []    # Option 3: Cookie
```

### Combined Requirements (AND)

Require multiple schemes together:

```yaml
security:
  - bearerAuth: []
    apiKeyAuth: []   # Both required
```

### OAuth Scopes

Specify required scopes:

```yaml
security:
  - oauth2Auth:
      - read
      - write
```

## Complete Example

```yaml title="config.yaml"
info:
  title: "Platform API"
  version: "1.0.0"
  description: |
    ## Authentication
    
    This API supports multiple authentication methods:
    
    - **Cookie**: Use `ory_kratos_session` cookie
    - **Bearer Token**: Use `Authorization: Bearer <token>` header
    - **API Key**: Use `X-API-Key` header

servers:
  - url: "https://api.example.com"

securitySchemes:
  # Cookie-based session (Ory Kratos)
  cookieAuth:
    type: apiKey
    in: cookie
    name: ory_kratos_session
    description: Ory Kratos session cookie

  # JWT Bearer token
  bearerAuth:
    type: http
    scheme: bearer
    bearerFormat: JWT
    description: JWT access token

  # API Key for service-to-service
  apiKeyAuth:
    type: apiKey
    in: header
    name: X-API-Key
    description: API key for service authentication

# Allow any of these authentication methods
security:
  - cookieAuth: []
  - bearerAuth: []
  - apiKeyAuth: []

inputs:
  - inputFile: apis/users.json
    pathModification:
      prepend: "/users"

output: platform-api.json
```

## Output Structure

The generated OpenAPI will include:

```json
{
  "components": {
    "securitySchemes": {
      "cookieAuth": {
        "type": "apiKey",
        "in": "cookie",
        "name": "ory_kratos_session",
        "description": "Ory Kratos session cookie"
      },
      "bearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT",
        "description": "JWT access token"
      }
    }
  },
  "security": [
    { "cookieAuth": [] },
    { "bearerAuth": [] }
  ]
}
```

## Tool Support

!!! warning "Cookie Auth in Tools"
    Some tools (like Postman) may not fully support cookie authentication when importing OpenAPI specs. You may need to manually configure cookies in the tool's cookie manager.

| Tool | Cookie Support | Bearer Support | OAuth2 Support |
|------|---------------|----------------|----------------|
| Swagger UI | ⚠️ Limited | ✅ Full | ✅ Full |
| Postman | ⚠️ Manual | ✅ Full | ✅ Full |
| Insomnia | ✅ Full | ✅ Full | ✅ Full |
| curl | ✅ `-b` flag | ✅ `-H` header | ⚠️ Manual |

## Next Steps

- [Operation Filtering](filtering.md) - Filter operations by tags and paths
- [Examples](../examples/advanced.md) - Advanced configuration examples

