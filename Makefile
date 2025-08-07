.PHONY: build test clean install help release-test release
.PHONY: test-unit test-integration test-all format lint vet tidy
.PHONY: docs security bench coverage pre-commit

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

# Run integration tests (requires QUIP_API_TOKEN env var)
test-integration:
	@if [ -z "$$QUIP_API_TOKEN" ]; then \
		echo "❌ QUIP_API_TOKEN environment variable is required for integration tests"; \
		echo "Set it with: export QUIP_API_TOKEN=your-token-here"; \
		echo "Get your token from: https://quip.com/dev/token"; \
		exit 1; \
	fi
	@echo "🧪 Running integration tests against real Quip API..."
	go test -v -tags=integration ./pkg/quip -run TestIntegration

# Run specific integration test
test-integration-single:
	@if [ -z "$$QUIP_API_TOKEN" ]; then \
		echo "❌ QUIP_API_TOKEN environment variable is required for integration tests"; \
		exit 1; \
	fi
	@if [ -z "$$TEST" ]; then \
		echo "❌ TEST variable is required. Example: make test-integration-single TEST=GetRecentThreads"; \
		exit 1; \
	fi
	@echo "🧪 Running integration test: $$TEST"
	go test -v -tags=integration ./pkg/quip -run TestIntegration_$$TEST

# Run integration benchmarks
test-integration-bench:
	@if [ -z "$$QUIP_API_TOKEN" ]; then \
		echo "❌ QUIP_API_TOKEN environment variable is required for integration tests"; \
		exit 1; \
	fi
	@echo "📊 Running integration benchmarks..."
	go test -v -tags=integration -bench=BenchmarkIntegration ./pkg/quip

# Run all tests (unit + integration)
test-all:
	@echo "🧪 Running unit tests..."
	go test -v ./...
	@echo ""
	@if [ -n "$$QUIP_API_TOKEN" ]; then \
		echo "🧪 Running integration tests..."; \
		go test -v -tags=integration ./pkg/quip -run TestIntegration; \
	else \
		echo "⚠️  Skipping integration tests (QUIP_API_TOKEN not set)"; \
	fi

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
format:
	@echo "🎨 Formatting code..."
	gofmt -w .
	go fmt ./...

# Alias for format (backwards compatibility)
fmt: format

# Lint code (requires golangci-lint)
lint:
	@echo "🔍 Linting code..."
	golangci-lint run

# Vet code
vet:
	@echo "🔬 Vetting code..."
	go vet ./...

# Tidy dependencies
tidy:
	@echo "🧹 Tidying dependencies..."
	go mod tidy
	go mod verify

# Security check (requires gosec)
security:
	@echo "🔒 Running security check..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "⚠️  gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Run benchmarks
bench:
	@echo "📊 Running benchmarks..."
	go test -bench=. -benchmem ./...

# Generate coverage report
coverage:
	@echo "📊 Generating coverage report..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Generate documentation
docs:
	@echo "📚 Generating documentation..."
	@mkdir -p docs
	go doc -all > docs/api.md
	@echo "Documentation generated: docs/api.md"

# Alias for test (backwards compatibility)  
test-unit: test

# Pre-commit workflow - run all quality checks
pre-commit:
	@echo "🚀 Running pre-commit checks..."
	@echo ""
	@$(MAKE) format
	@echo ""
	@$(MAKE) lint  
	@echo ""
	@$(MAKE) vet
	@echo ""
	@$(MAKE) test-all
	@echo ""
	@$(MAKE) build
	@echo ""
	@$(MAKE) tidy
	@echo ""
	@$(MAKE) docs
	@echo ""
	@echo "✅ All pre-commit checks passed! Ready to commit."

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

# Install gosec (for security scanning)
install-security:
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Setup development environment
dev-setup: install-goreleaser install-lint install-security deps
	@echo "🎉 Development environment setup complete!"
	@echo "Available tools:"
	@echo "  - golangci-lint: Code linting"
	@echo "  - gosec: Security scanning"
	@echo "  - goreleaser: Release management"
	@echo ""
	@echo "💡 Run 'make pre-commit' before committing changes"

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
	@echo ""
	@echo "🚀 DEVELOPMENT WORKFLOW:"
	@echo "  pre-commit              - Run all quality checks (format, lint, vet, test, build, tidy, docs)"
	@echo "  dev-setup               - Setup development environment"
	@echo ""
	@echo "🏗️  BUILD & RUN:"
	@echo "  build                   - Build the binary"
	@echo "  build-all               - Build for multiple platforms"
	@echo "  run                     - Run the binary"
	@echo "  inspect                 - Test with MCP Inspector"
	@echo ""
	@echo "🧪 TESTING:"
	@echo "  test                    - Run unit tests"
	@echo "  test-unit               - Run unit tests (alias)"
	@echo "  test-integration        - Run integration tests (requires QUIP_API_TOKEN)"
	@echo "  test-integration-single - Run specific integration test (TEST=TestName)"
	@echo "  test-integration-bench  - Run integration benchmarks"
	@echo "  test-all                - Run unit + integration tests"
	@echo "  bench                   - Run benchmarks"
	@echo "  coverage                - Generate coverage report"
	@echo ""
	@echo "🔍 CODE QUALITY:"
	@echo "  format                  - Format code (gofmt + go fmt)"
	@echo "  fmt                     - Format code (alias)"
	@echo "  lint                    - Lint code (requires golangci-lint)"
	@echo "  vet                     - Vet code (go vet)"
	@echo "  security                - Security check (requires gosec)"
	@echo ""
	@echo "📚 DOCUMENTATION & DEPENDENCIES:"
	@echo "  docs                    - Generate documentation"
	@echo "  deps                    - Install dependencies"
	@echo "  tidy                    - Tidy dependencies"
	@echo ""
	@echo "🗑️  CLEANUP:"
	@echo "  clean                   - Clean build artifacts"
	@echo ""
	@echo "🚢 RELEASE:"
	@echo "  release-test            - Test release process locally"
	@echo "  release                 - Create a new release"
	@echo "  release-check           - Check if ready for release"
	@echo "  tag-patch               - Create patch version tag"
	@echo "  tag-minor               - Create minor version tag"
	@echo "  tag-major               - Create major version tag"
	@echo ""
	@echo "❓ HELP:"
	@echo "  help                    - Show this help" 