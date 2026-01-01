# Input Files Configuration

Each input file can be individually configured with path modifications, filtering, and parameter handling.

## Input Configuration Structure

```yaml
inputs:
  - inputFile: "path/to/api.json"
    
    # Conflict resolution
    dispute:
      prefix: "ServiceName_"
    
    # Path transformations
    pathModification:
      stripStart: "/v1"
      prepend: "/service-name/v1"
    
    # Operation filtering
    operationSelection:
      includeTags: []
      excludeTags: ["Internal"]
      includePaths: []
      excludePaths: []
    
    # Parameter modifications
    includeExtraParameters:
      - name: "X-Request-ID"
        in: header
    excludeParameters:
      - name: "X-Internal-Header"
        in: header
    
    # Description handling
    description:
      append: true
      title:
        value: "Service Name API"
        headingLevel: 2
```

## InputConfig Properties

| Property | Type | Description |
|----------|------|-------------|
| `inputFile` | `string` | Path to the OpenAPI file (JSON or YAML) |
| `dispute` | `DisputeConfig` | Conflict resolution settings |
| `pathModification` | `PathModificationConfig` | Path transformation rules |
| `operationSelection` | `OperationSelectionConfig` | Operation filtering rules |
| `includeExtraParameters` | `[]ParameterConfig` | Parameters to inject |
| `excludeParameters` | `[]ParamFilter` | Parameters to remove |
| `description` | `DescriptionConfig` | Description handling |

## File Path Resolution

Paths can be absolute or relative to the config file:

```yaml
inputs:
  # Relative path (resolved from config file location)
  - inputFile: apis/users.json
  
  # Absolute path
  - inputFile: /opt/apis/orders.json
```

## Swagger 2.0 Support

Swagger 2.0 files are automatically converted to OpenAPI 3.0:

```yaml
inputs:
  # Swagger 2.0 file - auto-converted
  - inputFile: legacy-api.json
    pathModification:
      prepend: "/legacy"
```

!!! info "Automatic Detection"
    The tool detects the OpenAPI version by checking for `swagger: "2.0"` or `openapi: "3.x.x"` in the file.

## Conflict Resolution (Dispute)

When multiple files have components with the same name, use the `dispute` prefix:

```yaml
inputs:
  - inputFile: users-api.json
    dispute:
      prefix: "Users_"
  
  - inputFile: orders-api.json
    dispute:
      prefix: "Orders_"
```

### What Gets Prefixed

- Schemas: `User` → `Users_User`
- Responses: `ErrorResponse` → `Users_ErrorResponse`
- Parameters: `PageSize` → `Users_PageSize`
- Security Schemes: `BearerAuth` → `Users_BearerAuth`
- Request Bodies
- Headers
- Links
- Callbacks

### Reference Updates

All `$ref` references are automatically updated:

```yaml
# Before
$ref: "#/components/schemas/User"

# After (with prefix "Users_")
$ref: "#/components/schemas/Users_User"
```

!!! warning "No Prefix = Error on Collision"
    If two files have the same schema name and no dispute prefix is set, the merge will fail with a collision error.

## Description Handling

Append input API descriptions to the merged output:

```yaml
inputs:
  - inputFile: users-api.json
    description:
      append: true
      title:
        value: "Users Service"
        headingLevel: 2
```

This adds to the merged description:

```markdown
## Users Service

[Original description from users-api.json]
```

### Description Properties

| Property | Type | Description |
|----------|------|-------------|
| `append` | `boolean` | Whether to append the description |
| `title.value` | `string` | Title for the description section |
| `title.headingLevel` | `integer` | Markdown heading level (1-6) |

## Multiple Inputs Example

```yaml
inputs:
  # IAM Service
  - inputFile: apis/iam.json
    dispute:
      prefix: ""  # No prefix needed (no collisions)
    pathModification:
      stripStart: "/v1"
      prepend: "/iam-service/v1"
    excludeParameters:
      - name: X-Internal-Auth
        in: header
    description:
      append: true
      title:
        value: "Identity & Access Management"
        headingLevel: 2

  # Orders Service  
  - inputFile: apis/orders.json
    dispute:
      prefix: "Orders_"
    pathModification:
      stripStart: "/api"
      prepend: "/orders-service/v1"
    operationSelection:
      excludeTags: ["Internal", "Admin"]
    description:
      append: true
      title:
        value: "Orders Service"
        headingLevel: 2

  # Legacy Service (Swagger 2.0)
  - inputFile: apis/legacy.yaml
    dispute:
      prefix: "Legacy_"
    pathModification:
      prepend: "/legacy/v1"
    description:
      append: true
      title:
        value: "Legacy API (Deprecated)"
        headingLevel: 2
```

## Next Steps

- [Path Modification](paths.md) - Detailed path transformation guide
- [Operation Filtering](filtering.md) - Include/exclude operations
- [Security Schemes](security.md) - Authentication configuration

