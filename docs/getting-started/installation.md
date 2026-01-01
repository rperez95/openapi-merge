# Installation

## Requirements

- **Go 1.22+** (for building from source)
- No runtime dependencies

## Installation Methods

### GitHub Action

For CI/CD pipelines, use the official GitHub Action:

```yaml
- uses: rperez95/openapi-merge@v0
  with:
    config: merge-config.yaml
```

See the [GitHub Actions](../github-actions.md) documentation for detailed usage.

### Using Go Install

The easiest way to install openapi-merge locally:

```bash
go install github.com/rperez95/openapi-merge@latest
```

This installs the binary to your `$GOPATH/bin` directory.

### Building from Source

Clone and build the project:

```bash
# Clone the repository
git clone https://github.com/rperez95/openapi-merge.git
cd openapi-merge

# Build
go build -o openapi-merge .

# Optional: Install to PATH
sudo mv openapi-merge /usr/local/bin/
```

### Download Binary

Download pre-built binaries from the [GitHub Releases](https://github.com/rperez95/openapi-merge/releases) page.

=== "Linux (amd64)"

    ```bash
    curl -LO https://github.com/rperez95/openapi-merge/releases/latest/download/openapi-merge-linux-amd64.tar.gz
    tar -xzf openapi-merge-linux-amd64.tar.gz
    chmod +x openapi-merge
    sudo mv openapi-merge /usr/local/bin/
    ```

=== "macOS (amd64)"

    ```bash
    curl -LO https://github.com/rperez95/openapi-merge/releases/latest/download/openapi-merge-darwin-amd64.tar.gz
    tar -xzf openapi-merge-darwin-amd64.tar.gz
    chmod +x openapi-merge
    sudo mv openapi-merge /usr/local/bin/
    ```

=== "macOS (arm64)"

    ```bash
    curl -LO https://github.com/rperez95/openapi-merge/releases/latest/download/openapi-merge-darwin-arm64.tar.gz
    tar -xzf openapi-merge-darwin-arm64.tar.gz
    chmod +x openapi-merge
    sudo mv openapi-merge /usr/local/bin/
    ```

=== "Windows"

    ```powershell
    # Download from releases page and add to PATH
    ```

## Verify Installation

```bash
openapi-merge --help
```

Expected output:

```
openapi-merge is a CLI tool that merges multiple OpenAPI 2.0 (Swagger) 
and OpenAPI 3.0/3.1 specifications into a single valid OpenAPI 3.0 file.

Usage:
  openapi-merge [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  merge       Merge OpenAPI specifications based on config

Flags:
      --config string   config file (required)
  -h, --help            help for openapi-merge
  -v, --verbose         enable verbose output

Use "openapi-merge [command] --help" for more information about a command.
```

## Shell Completion

Generate shell completion scripts for your shell:

=== "Bash"

    ```bash
    openapi-merge completion bash > /etc/bash_completion.d/openapi-merge
    ```

=== "Zsh"

    ```bash
    openapi-merge completion zsh > "${fpath[1]}/_openapi-merge"
    ```

=== "Fish"

    ```bash
    openapi-merge completion fish > ~/.config/fish/completions/openapi-merge.fish
    ```

## Next Steps

Now that you have openapi-merge installed, check out the [Quick Start](quickstart.md) guide to create your first merged API specification.

