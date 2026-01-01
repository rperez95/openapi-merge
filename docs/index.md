# openapi-merge

<p align="center">
  <strong>Merge multiple OpenAPI specifications into a single unified API</strong>
</p>

<p align="center">
  <a href="getting-started/installation/">Getting Started</a> •
  <a href="configuration/overview/">Configuration</a> •
  <a href="examples/basic/">Examples</a> •
  <a href="cli/">CLI Reference</a>
</p>

---

## What is openapi-merge?

**openapi-merge** is a CLI tool that merges multiple OpenAPI 2.0 (Swagger) and OpenAPI 3.0/3.1 specifications into a single valid OpenAPI 3.0 file.

This is primarily used for **API Gateways** where multiple microservices need to be exposed under a single unified schema.

## Key Features

<div class="grid cards" markdown>

-   :material-merge:{ .lg .middle } __Merge Multiple Specs__

    ---

    Combine any number of OpenAPI 2.0 or 3.0 specifications into one unified API document.

-   :material-swap-horizontal:{ .lg .middle } __Auto-Convert Swagger 2.0__

    ---

    Automatically converts Swagger 2.0 specs to OpenAPI 3.0 during the merge process.

-   :material-road-variant:{ .lg .middle } __Path Modification__

    ---

    Strip prefixes, add service-specific paths, and reorganize your API structure.

-   :material-filter:{ .lg .middle } __Operation Filtering__

    ---

    Include or exclude operations by tags, paths, or HTTP methods with glob pattern support.

-   :material-shield-lock:{ .lg .middle } __Security Schemes__

    ---

    Define OAS3-compliant security schemes (Cookie, Bearer, OAuth2, OpenID Connect).

-   :material-puzzle:{ .lg .middle } __Conflict Resolution__

    ---

    Handle schema name collisions with automatic prefixing and reference updates.

-   :material-github:{ .lg .middle } __GitHub Actions__

    ---

    Native CI/CD integration with the official GitHub Action for automated merging.

</div>

## Quick Example

```yaml title="merge-config.yaml"
info:
  title: "Unified Platform API"
  version: "1.0.0"

servers:
  - url: "https://api.example.com"

securitySchemes:
  cookieAuth:
    type: apiKey
    in: cookie
    name: session_id

security:
  - cookieAuth: []

inputs:
  - inputFile: apis/users.json
    pathModification:
      prepend: "/users-service/v1"

  - inputFile: apis/orders.json
    pathModification:
      prepend: "/orders-service/v1"

output: unified-api.json
```

```bash
openapi-merge merge --config merge-config.yaml
```

## Why openapi-merge?

| Challenge | Solution |
|-----------|----------|
| Multiple microservices with separate OpenAPI specs | Merge into single spec for API Gateway |
| Mixed Swagger 2.0 and OpenAPI 3.0 files | Auto-convert to OpenAPI 3.0 |
| Schema name collisions across services | Automatic prefixing with dispute resolution |
| Need to filter internal endpoints | Tag and path-based filtering |
| Complex authentication requirements | Full OAS3 security scheme support |

## Installation

=== "GitHub Action"

    ```yaml
    - uses: rperez95/openapi-merge@v0
      with:
        config: merge-config.yaml
    ```

=== "Go Install"

    ```bash
    go install github.com/rperez95/openapi-merge@latest
    ```

=== "From Source"

    ```bash
    git clone https://github.com/rperez95/openapi-merge.git
    cd openapi-merge
    go build -o openapi-merge .
    ```

=== "Binary Release"

    Download the latest release from [GitHub Releases](https://github.com/rperez95/openapi-merge/releases).

## Next Steps

- [Installation Guide](getting-started/installation.md) - Detailed installation instructions
- [Quick Start](getting-started/quickstart.md) - Get up and running in 5 minutes
- [Configuration Reference](configuration/overview.md) - Full configuration options
- [GitHub Actions](github-actions.md) - CI/CD integration guide
- [Examples](examples/basic.md) - Real-world usage examples

