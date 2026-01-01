# Contributing

Thank you for your interest in contributing to openapi-merge!

## Getting Started

### Prerequisites

- Go 1.23 or later
- Git

### Clone the Repository

```bash
git clone https://github.com/rperez95/openapi-merge.git
cd openapi-merge
```

### Install Dependencies

```bash
go mod download
```

### Build

```bash
go build -o openapi-merge .
```

### Run Tests

```bash
go test ./... -v
```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes

- Follow Go conventions and best practices
- Add tests for new functionality
- Update documentation if needed

### 3. Run Tests

```bash
# Run all tests
go test ./... -v

# Run with coverage
go test ./... -cover

# Run specific tests
go test ./internal/merger/... -v
```

### 4. Lint

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

### 5. Commit

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```bash
git commit -m "feat: add support for X"
git commit -m "fix: resolve issue with Y"
git commit -m "docs: update configuration guide"
```

### 6. Push and Create PR

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

## Code Style

### Go Conventions

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting
- Add comments for exported functions

### Example

```go
// MergeSpecs merges multiple OpenAPI specifications into one.
// It handles Swagger 2.0 to OpenAPI 3.0 conversion automatically.
func MergeSpecs(inputs []InputConfig) (*openapi3.T, error) {
    // Implementation
}
```

## Testing

### Unit Tests

Add tests in `*_test.go` files:

```go
func TestMerger_PathModification(t *testing.T) {
    // Setup
    cfg := &config.Config{
        Inputs: []config.InputConfig{
            {
                InputFile: "test.json",
                PathModification: &config.PathModificationConfig{
                    StripStart: "/v1",
                    Prepend:    "/api",
                },
            },
        },
    }

    // Execute
    result, err := merge(cfg)

    // Assert
    assert.NoError(t, err)
    assert.Contains(t, result.Paths, "/api/users")
}
```

### Test Data

Place test OpenAPI files in `testdata/` directory:

```
internal/merger/testdata/
├── valid-openapi3.json
├── valid-swagger2.json
└── invalid.json
```

## Documentation

### MkDocs

Documentation is in `docs/` using MkDocs Material.

```bash
# Install MkDocs
pip install mkdocs-material mkdocs-minify-plugin

# Serve locally
mkdocs serve

# Build
mkdocs build
```

### Adding Pages

1. Create markdown file in `docs/`
2. Add to `nav` in `mkdocs.yml`

## Reporting Issues

### Bug Reports

Include:

- openapi-merge version
- Go version
- Operating system
- Configuration file (sanitized)
- Input files (if possible)
- Error message
- Steps to reproduce

### Feature Requests

Describe:

- Use case
- Expected behavior
- Example configuration

## Pull Request Guidelines

- [ ] Tests pass
- [ ] Linter passes
- [ ] Documentation updated (if needed)
- [ ] Follows conventional commits
- [ ] Single purpose (one feature/fix per PR)

## License

By contributing, you agree that your contributions will be licensed under the Apache 2.0 License.

## Questions?

- Open a [GitHub Issue](https://github.com/rperez95/openapi-merge/issues)
- Start a [Discussion](https://github.com/rperez95/openapi-merge/discussions)

Thank you for contributing! 🎉

