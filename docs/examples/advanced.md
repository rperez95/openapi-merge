# Advanced Configuration

Advanced examples for complex merge scenarios.

## Schema Collision Resolution

When multiple APIs have schemas with the same name:

```yaml
inputs:
  # Both APIs have a "User" schema
  - inputFile: auth-api.json
    dispute:
      prefix: "Auth_"  # Auth_User, Auth_Session, etc.
    pathModification:
      prepend: "/auth"

  - inputFile: profile-api.json
    dispute:
      prefix: "Profile_"  # Profile_User, Profile_Settings, etc.
    pathModification:
      prepend: "/profile"
```

All references are automatically updated:

```yaml
# Before
$ref: "#/components/schemas/User"

# After (auth-api)
$ref: "#/components/schemas/Auth_User"

# After (profile-api)
$ref: "#/components/schemas/Profile_User"
```

## OAuth2 with Multiple Scopes

```yaml
securitySchemes:
  oauth2:
    type: oauth2
    description: OAuth 2.0 authentication
    flows:
      authorizationCode:
        authorizationUrl: https://auth.example.com/oauth/authorize
        tokenUrl: https://auth.example.com/oauth/token
        refreshUrl: https://auth.example.com/oauth/refresh
        scopes:
          users:read: Read user information
          users:write: Modify user information
          orders:read: Read orders
          orders:write: Create and modify orders
          admin: Full administrative access

# Require specific scopes
security:
  - oauth2:
      - users:read
      - orders:read
```

## Multi-Environment Configuration

Create environment-specific configs:

=== "config.production.yaml"

    ```yaml
    info:
      title: "Platform API"
      version: "1.0.0"
    
    servers:
      - url: "https://api.example.com"
        description: "Production"
    
    inputs:
      - inputFile: apis/users.json
        operationSelection:
          excludeTags: ["Internal", "Debug", "Beta"]
    
    output: dist/production-api.json
    ```

=== "config.staging.yaml"

    ```yaml
    info:
      title: "Platform API (Staging)"
      version: "1.0.0-staging"
    
    servers:
      - url: "https://staging-api.example.com"
        description: "Staging"
    
    inputs:
      - inputFile: apis/users.json
        operationSelection:
          excludeTags: ["Internal", "Debug"]  # Include Beta
    
    output: dist/staging-api.json
    ```

=== "config.development.yaml"

    ```yaml
    info:
      title: "Platform API (Development)"
      version: "1.0.0-dev"
    
    servers:
      - url: "http://localhost:8080"
        description: "Local"
    
    inputs:
      - inputFile: apis/users.json
        # Include everything in dev
    
    output: dist/development-api.json
    ```

Build all environments:

```bash
openapi-merge merge --config config.production.yaml
openapi-merge merge --config config.staging.yaml
openapi-merge merge --config config.development.yaml
```

## Filtering Internal Endpoints

Create a public API by filtering out internal endpoints:

```yaml
inputs:
  - inputFile: full-api.json
    operationSelection:
      # Only include public tags
      includeTags:
        - "Public"
        - "Users"
        - "Products"
      
      # Exclude internal tags
      excludeTags:
        - "Internal"
        - "Admin"
        - "Debug"
        - "Metrics"
      
      # Exclude internal paths
      excludePaths:
        - path: "/internal/**"
        - path: "/admin/**"
        - path: "/debug/**"
        - path: "/metrics"
        - path: "/health"
        - path: "/ready"
    
    # Remove internal headers
    excludeParameters:
      - name: X-Internal-Auth
        in: header
      - name: X-Debug-Mode
        in: header
      - name: X-Admin-Override
        in: header
```

## Adding Gateway Headers

Inject headers required by your API Gateway:

```yaml
inputs:
  - inputFile: api.json
    includeExtraParameters:
      # Request tracking
      - name: X-Request-ID
        in: header
        description: "Unique request identifier for tracing"
        required: false
        schema:
          type: string
          format: uuid
      
      # Tenant identification
      - name: X-Tenant-ID
        in: header
        description: "Tenant identifier for multi-tenancy"
        required: true
        schema:
          type: string
      
      # API versioning
      - name: X-API-Version
        in: header
        description: "API version override"
        required: false
        schema:
          type: string
          default: "2024-01"
      
      # Distributed tracing (Zipkin/Jaeger)
      - name: X-B3-TraceId
        in: header
        description: "Trace ID for distributed tracing"
        required: false
        schema:
          type: string
      
      - name: X-B3-SpanId
        in: header
        description: "Span ID for distributed tracing"
        required: false
        schema:
          type: string
```

## Version Migration

Merge old and new API versions with clear separation:

```yaml
basePath: "/api"

inputs:
  # Legacy v1 API
  - inputFile: apis/v1/legacy.json
    dispute:
      prefix: "V1_"
    pathModification:
      prepend: "/v1"
    description:
      append: true
      title:
        value: "API v1 (Deprecated)"
        headingLevel: 2
  
  # New v2 API
  - inputFile: apis/v2/modern.json
    dispute:
      prefix: "V2_"
    pathModification:
      prepend: "/v2"
    description:
      append: true
      title:
        value: "API v2 (Current)"
        headingLevel: 2

output: combined-api.json
```

Result:
```
/api/v1/users      # Legacy
/api/v1/orders     # Legacy
/api/v2/users      # Modern
/api/v2/orders     # Modern
```

## Complex Tag Ordering

Organize documentation with logical tag grouping:

```yaml
tagOrder:
  # Authentication first
  - "Authentication"
  - "Authorization"
  
  # Core resources
  - "Users"
  - "Groups"
  - "Organizations"
  
  # Business logic
  - "Orders"
  - "Products"
  - "Inventory"
  
  # Infrastructure
  - "Webhooks"
  - "Events"
  
  # Utility
  - "Health"
  - "Metrics"
```

## Makefile Integration

```makefile title="Makefile"
.PHONY: api-docs api-docs-prod api-docs-staging clean

# Build all API documentation
api-docs: api-docs-prod api-docs-staging api-docs-dev

api-docs-prod:
	openapi-merge merge --config config/production.yaml -o dist/api-production.json

api-docs-staging:
	openapi-merge merge --config config/staging.yaml -o dist/api-staging.json

api-docs-dev:
	openapi-merge merge --config config/development.yaml -o dist/api-development.json

# Validate merged specs
validate:
	npx @stoplight/spectral-cli lint dist/*.json

# Clean output
clean:
	rm -rf dist/*.json

# Watch for changes (requires entr)
watch:
	find apis/ config/ -name '*.json' -o -name '*.yaml' | entr -r make api-docs
```

## Next Steps

- [CLI Reference](../cli.md) - All command-line options
- [Configuration Reference](../configuration/overview.md) - Full configuration options

