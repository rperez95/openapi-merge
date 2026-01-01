# GitHub Actions

Automate OpenAPI specification merging in your CI/CD pipelines using the official GitHub Action.

## Quick Start

Add this step to your workflow:

```yaml
- name: Merge OpenAPI specs
  uses: rperez95/openapi-merge@v0
  with:
    config: merge-config.yaml
```

## Action Reference

### Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `config` | Path to the merge configuration file (YAML or JSON) | ✅ | — |
| `output` | Output file path (overrides config file setting) | ❌ | — |
| `verbose` | Enable verbose output (`true` or `false`) | ❌ | `false` |
| `version` | Version of openapi-merge to use | ❌ | `latest` |

### Outputs

| Output | Description |
|--------|-------------|
| `output-file` | Path to the generated merged OpenAPI file |

## Usage Examples

### Basic Merge

The simplest usage — merge specs whenever APIs change:

```yaml
name: Merge OpenAPI Specs

on:
  push:
    branches: [main]
    paths:
      - 'apis/**'
      - 'merge-config.yaml'

jobs:
  merge:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Merge OpenAPI specifications
        uses: rperez95/openapi-merge@v0
        with:
          config: merge-config.yaml
```

### Upload as Artifact

Save the merged spec as a build artifact:

```yaml
name: Build API Spec

on:
  push:
    branches: [main]

jobs:
  merge:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Merge OpenAPI specifications
        uses: rperez95/openapi-merge@v0
        with:
          config: merge-config.yaml
          output: dist/api-spec.json
          verbose: 'true'

      - name: Upload merged spec
        uses: actions/upload-artifact@v4
        with:
          name: api-specification
          path: dist/api-spec.json
          retention-days: 30
```

### Validate with Redocly

Merge and validate the output spec:

```yaml
name: API Validation

on:
  pull_request:
    paths:
      - 'apis/**'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Merge OpenAPI specifications
        id: merge
        uses: rperez95/openapi-merge@v0
        with:
          config: merge-config.yaml
          output: merged-api.json

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Validate merged spec
        run: npx @redocly/cli lint ${{ steps.merge.outputs.output-file }}
```

### Generate Documentation

Merge specs and generate API documentation:

```yaml
name: API Documentation

on:
  push:
    branches: [main]

jobs:
  docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Merge OpenAPI specifications
        id: merge
        uses: rperez95/openapi-merge@v0
        with:
          config: docs/merge-config.yaml
          output: docs/openapi.json

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Generate Redoc documentation
        run: |
          npx @redocly/cli build-docs ${{ steps.merge.outputs.output-file }} \
            --output docs/api.html \
            --title "API Documentation"

      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v4
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./docs
```

### Multi-Environment Specs

Generate different specs for different environments:

```yaml
name: Multi-Environment API Specs

on:
  push:
    branches: [main]

jobs:
  merge:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        env: [production, staging, development]
    steps:
      - uses: actions/checkout@v4

      - name: Merge OpenAPI specifications
        uses: rperez95/openapi-merge@v0
        with:
          config: config/merge-${{ matrix.env }}.yaml
          output: dist/api-${{ matrix.env }}.json
          verbose: 'true'

      - name: Upload spec
        uses: actions/upload-artifact@v4
        with:
          name: api-spec-${{ matrix.env }}
          path: dist/api-${{ matrix.env }}.json
```

### Pin to Specific Version

For reproducible builds, pin to a specific version:

```yaml
- name: Merge OpenAPI specifications
  uses: rperez95/openapi-merge@v0
  with:
    config: merge-config.yaml
    version: 'v1.0.0'  # Pin to specific version
```

### Commit Merged Spec

Automatically commit the merged spec back to the repository:

```yaml
name: Update API Spec

on:
  push:
    branches: [main]
    paths:
      - 'apis/**'

jobs:
  update-spec:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4

      - name: Merge OpenAPI specifications
        uses: rperez95/openapi-merge@v0
        with:
          config: merge-config.yaml
          output: openapi.json

      - name: Commit changes
        uses: stefanzweifel/git-auto-commit-action@v5
        with:
          commit_message: 'chore: update merged OpenAPI spec'
          file_pattern: 'openapi.json'
```

### Release Workflow

Include merged spec in releases:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4

      - name: Merge OpenAPI specifications
        uses: rperez95/openapi-merge@v0
        with:
          config: merge-config.yaml
          output: api-spec.json

      - name: Create Release
        uses: softprops/action-gh-release@v2
        with:
          files: |
            api-spec.json
          generate_release_notes: true
```

## Complete Example Workflow

Here's a comprehensive workflow that merges, validates, generates docs, and deploys:

```yaml
name: API Pipeline

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  merge:
    runs-on: ubuntu-latest
    outputs:
      output-file: ${{ steps.merge.outputs.output-file }}
    steps:
      - uses: actions/checkout@v4

      - name: Merge OpenAPI specifications
        id: merge
        uses: rperez95/openapi-merge@v0
        with:
          config: merge-config.yaml
          output: dist/openapi.json
          verbose: 'true'

      - name: Upload merged spec
        uses: actions/upload-artifact@v4
        with:
          name: openapi-spec
          path: dist/openapi.json

  validate:
    needs: merge
    runs-on: ubuntu-latest
    steps:
      - name: Download merged spec
        uses: actions/download-artifact@v4
        with:
          name: openapi-spec

      - name: Validate with Spectral
        uses: stoplightio/spectral-action@v0.8.11
        with:
          file_glob: 'openapi.json'

  generate-docs:
    needs: [merge, validate]
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Download merged spec
        uses: actions/download-artifact@v4
        with:
          name: openapi-spec
          path: docs/

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Generate documentation
        run: |
          npx @redocly/cli build-docs docs/openapi.json \
            --output docs/index.html

      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v4
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./docs
```

## Troubleshooting

### Config file not found

```
Error: failed to load configuration: open merge-config.yaml: no such file or directory
```

!!! tip "Solution"
    Ensure you have `actions/checkout@v4` before the merge step, and verify the config path is correct relative to the repository root.

### Permission denied writing output

```
Error: failed to write output file: permission denied
```

!!! tip "Solution"
    The output directory must exist. Create it in a previous step:
    ```yaml
    - run: mkdir -p dist
    - uses: rperez95/openapi-merge@v0
      with:
        config: merge-config.yaml
        output: dist/api.json
    ```

### Input files not found

```
Error: failed to load apis/service.json: no such file or directory
```

!!! tip "Solution"
    Input file paths in the config are relative to the config file location. Verify paths are correct after checkout.

### Using output in subsequent steps

Access the merged file path in later steps:

```yaml
- name: Merge specs
  id: merge
  uses: rperez95/openapi-merge@v0
  with:
    config: merge-config.yaml
    output: api.json

- name: Use output
  run: |
    echo "Merged file: ${{ steps.merge.outputs.output-file }}"
    cat ${{ steps.merge.outputs.output-file }} | head -20
```

## Self-Hosted Runners

The action works on all runner types:

=== "GitHub-hosted"

    ```yaml
    runs-on: ubuntu-latest
    ```

=== "Self-hosted Linux"

    ```yaml
    runs-on: self-hosted
    # Requires Go 1.22+ or the action will install it
    ```

=== "Self-hosted macOS"

    ```yaml
    runs-on: macos-latest
    # Works on both Intel and Apple Silicon
    ```

=== "Self-hosted Windows"

    ```yaml
    runs-on: windows-latest
    # Uses PowerShell-compatible commands
    ```

## See Also

- [CLI Reference](cli.md) — Command-line usage
- [Configuration Overview](configuration/overview.md) — Config file format
- [Examples](examples/basic.md) — More usage examples

