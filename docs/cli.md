# CLI Reference

Complete command-line interface reference for openapi-merge.

## Global Options

```bash
openapi-merge [command] [flags]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--config` | | Configuration file path (required) |
| `--verbose` | `-v` | Enable verbose output |
| `--help` | `-h` | Show help |

## Commands

### merge

Merge OpenAPI specifications based on configuration file.

```bash
openapi-merge merge --config <config-file> [flags]
```

#### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--config` | | Configuration file path (required) |
| `--output` | `-o` | Override output file path |
| `--verbose` | `-v` | Enable verbose output |

#### Examples

```bash
# Basic merge
openapi-merge merge --config config.yaml

# With verbose output
openapi-merge merge --config config.yaml --verbose

# Override output file
openapi-merge merge --config config.yaml -o custom-output.json

# Output as YAML
openapi-merge merge --config config.yaml -o output.yaml
```

### completion

Generate shell completion scripts.

```bash
openapi-merge completion [bash|zsh|fish|powershell]
```

#### Bash

```bash
# Generate completion script
openapi-merge completion bash > /etc/bash_completion.d/openapi-merge

# Or for current session
source <(openapi-merge completion bash)
```

#### Zsh

```bash
# Generate completion script
openapi-merge completion zsh > "${fpath[1]}/_openapi-merge"

# Or add to .zshrc
echo 'source <(openapi-merge completion zsh)' >> ~/.zshrc
```

#### Fish

```bash
openapi-merge completion fish > ~/.config/fish/completions/openapi-merge.fish
```

#### PowerShell

```powershell
openapi-merge completion powershell | Out-String | Invoke-Expression
```

### help

Show help for any command.

```bash
openapi-merge help [command]
```

#### Examples

```bash
# General help
openapi-merge help

# Help for merge command
openapi-merge help merge

# Alternative syntax
openapi-merge merge --help
```

### version

Show version information.

```bash
openapi-merge version
```

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | General error (invalid config, file not found, etc.) |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `NO_COLOR` | Disable colored output |

## Usage Examples

### Basic Usage

```bash
# Merge with config file
openapi-merge merge --config merge-config.yaml
```

### Verbose Output

```bash
# See detailed progress
openapi-merge merge --config merge-config.yaml -v
```

Output:
```
Starting merge with 3 input files
Output file: /path/to/output.json
Processing input 1: apis/users.json
Processing input 2: apis/orders.json
  Detected Swagger 2.0, converting to OpenAPI 3.0
Processing input 3: apis/products.json
Applied global basePath: /api/v1
Successfully merged 3 specifications into /path/to/output.json
```

### Custom Output Path

```bash
# Override config output
openapi-merge merge --config config.yaml -o /tmp/api.json

# Output as YAML
openapi-merge merge --config config.yaml -o /tmp/api.yaml
```

### Pipeline Usage

```bash
# Use in scripts
if openapi-merge merge --config config.yaml; then
    echo "Merge successful"
    # Continue with deployment
else
    echo "Merge failed"
    exit 1
fi
```

### Multiple Configurations

```bash
# Build multiple outputs
for env in production staging development; do
    openapi-merge merge --config "config.${env}.yaml" -o "dist/api-${env}.json"
done
```

### Watch Mode (with external tool)

```bash
# Using entr (install: brew install entr / apt install entr)
find apis/ -name '*.json' | entr -r openapi-merge merge --config config.yaml

# Using nodemon
npx nodemon --watch apis/ --ext json,yaml --exec "openapi-merge merge --config config.yaml"
```

### Docker Usage

```dockerfile
FROM golang:1.22-alpine AS builder
RUN go install github.com/rperez95/openapi-merge@latest

FROM alpine:latest
COPY --from=builder /go/bin/openapi-merge /usr/local/bin/
ENTRYPOINT ["openapi-merge"]
```

```bash
# Build and run
docker build -t openapi-merge .
docker run -v $(pwd):/work -w /work openapi-merge merge --config config.yaml
```

## Configuration File Formats

The tool accepts both YAML and JSON configuration files:

=== "YAML"

    ```yaml
    info:
      title: "API"
      version: "1.0.0"
    
    inputs:
      - inputFile: api.json
    
    output: merged.json
    ```

=== "JSON"

    ```json
    {
      "info": {
        "title": "API",
        "version": "1.0.0"
      },
      "inputs": [
        { "inputFile": "api.json" }
      ],
      "output": "merged.json"
    }
    ```

## Troubleshooting

### Config File Not Found

```
Error: failed to load configuration: open config.yaml: no such file or directory
```

**Solution**: Check the config file path is correct and the file exists.

### Input File Not Found

```
Error: failed to load apis/missing.json: open apis/missing.json: no such file or directory
```

**Solution**: Verify the input file path in your config. Paths are relative to the config file location.

### Invalid OpenAPI Spec

```
Error: failed to load apis/invalid.json: failed to parse JSON: ...
```

**Solution**: Validate your OpenAPI file with a tool like [Swagger Editor](https://editor.swagger.io/).

### Schema Collision

```
Error: schema collision for 'User' without dispute prefix
```

**Solution**: Add a `dispute.prefix` to one or both conflicting inputs:

```yaml
inputs:
  - inputFile: api1.json
    dispute:
      prefix: "Api1_"
```

### Permission Denied

```
Error: failed to write output file: permission denied
```

**Solution**: Check write permissions for the output directory.

