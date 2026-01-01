.PHONY: all build test lint fmt vet clean install run help docs docs-serve

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
BINARY_NAME=openapi-merge

# Build info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

all: lint test build

## build: Build the binary
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

## install: Install binary to GOPATH/bin
install:
	$(GOCMD) install $(LDFLAGS) .

## test: Run tests
test:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

## test-short: Run tests without race detector
test-short:
	$(GOTEST) -v ./...

## coverage: Show test coverage
coverage: test
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## lint: Run linters
lint: fmt vet golangci-lint

## fmt: Format code
fmt:
	$(GOFMT) ./...

## vet: Run go vet
vet:
	$(GOVET) ./...

## golangci-lint: Run golangci-lint
golangci-lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

## tidy: Tidy go modules
tidy:
	$(GOMOD) tidy

## clean: Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	rm -rf dist/

## run: Run the merge command with example config
run: build
	./$(BINARY_NAME) merge --config examples/merge-config.yaml -v

## docs: Build documentation
docs:
	cd docs && pip install -r requirements.txt
	mkdocs build --strict

## docs-serve: Serve documentation locally
docs-serve:
	mkdocs serve

## release-dry: Test goreleaser without publishing
release-dry:
	goreleaser release --snapshot --clean

## help: Show this help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/  /'

# Default target
.DEFAULT_GOAL := help

