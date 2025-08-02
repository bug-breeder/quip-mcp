.PHONY: build test clean install help release-test release

# Default target
all: test build

# Build the binary
build:
	go build -o quip-mcp .

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -f quip-mcp quip-mcp-* coverage.out coverage.html

# Install dependencies
deps:
	go mod download
	go mod tidy

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o quip-mcp-linux .
	GOOS=darwin GOARCH=amd64 go build -o quip-mcp-darwin .
	GOOS=windows GOARCH=amd64 go build -o quip-mcp-windows.exe .

# Run the binary (requires QUIP_API_TOKEN env var)
run:
	./quip-mcp

# Test with MCP Inspector (requires npm)
inspect:
	npx @modelcontextprotocol/inspector ./quip-mcp

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Test release process locally
release-test:
	goreleaser release --snapshot --clean

# Create a new release (requires tag)
release:
	goreleaser release --clean

# Install GoReleaser (for development)
install-goreleaser:
	go install github.com/goreleaser/goreleaser@latest

# Install golangci-lint (for development)
install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Setup development environment
dev-setup: install-goreleaser install-lint deps
	@echo "Development environment setup complete!"

# Check if ready for release
release-check:
	@echo "Checking release readiness..."
	@go test ./...
	@golangci-lint run
	@goreleaser check
	@echo "✅ Ready for release!"

# Version management helpers
tag-patch:
	@echo "Creating patch version tag..."
	@git tag $$(git describe --tags $$(git rev-list --tags --max-count=1) | awk -F. '{$$NF = $$NF + 1;} 1' | sed 's/ /./g')

tag-minor:
	@echo "Creating minor version tag..."
	@git tag $$(git describe --tags $$(git rev-list --tags --max-count=1) | awk -F. '{$$(NF-1) = $$(NF-1) + 1; $$NF = 0;} 1' | sed 's/ /./g')

tag-major:
	@echo "Creating major version tag..."
	@git tag $$(git describe --tags $$(git rev-list --tags --max-count=1) | awk -F. '{$$(NF-2) = $$(NF-2) + 1; $$(NF-1) = 0; $$NF = 0;} 1' | sed 's/ /./g')

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  run           - Run the binary"
	@echo "  inspect       - Test with MCP Inspector"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  release-test  - Test release process locally"
	@echo "  release       - Create a new release"
	@echo "  release-check - Check if ready for release"
	@echo "  dev-setup     - Setup development environment"
	@echo "  tag-patch     - Create patch version tag"
	@echo "  tag-minor     - Create minor version tag"
	@echo "  tag-major     - Create major version tag"
	@echo "  help          - Show this help" 