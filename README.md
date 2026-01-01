# openapi-merge

[![CI](https://github.com/rperez95/openapi-merge/actions/workflows/ci.yml/badge.svg)](https://github.com/rperez95/openapi-merge/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/rperez95/openapi-merge/graph/badge.svg)](https://codecov.io/gh/rperez95/openapi-merge)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

📚 **[Documentation](https://rperez95.github.io/openapi-merge/)**

A CLI tool for merging multiple OpenAPI 2.0 (Swagger) and OpenAPI 3.0/3.1 specifications into a single unified OpenAPI 3.0 file.

This tool is a Go implementation inspired by [robertmassaioli/openapi-merge](https://github.com/robertmassaioli/openapi-merge).

## Features

- 🔀 Merge multiple OpenAPI specs (Swagger 2.0 & OpenAPI 3.x)
- 🔄 Auto-convert Swagger 2.0 to OpenAPI 3.0
- 🛣️ Path modification (strip/prepend prefixes)
- 🏷️ Filter operations by tags, paths, or methods
- 🔐 Full OAS3 security scheme support
- ⚡ Conflict resolution with dispute prefixes
- 🤖 Native GitHub Actions integration

## Quick Start

### GitHub Action

```yaml
- uses: rperez95/openapi-merge@v0
  with:
    config: examples/merge-config.yaml
    output: examples/platform-api.json
```

### CLI

```bash
go install github.com/rperez95/openapi-merge@latest

openapi-merge merge --config examples/merge-config.yaml -o examples/platform-api.json
```

## License

Apache 2.0 — see [LICENSE](LICENSE) for details.
