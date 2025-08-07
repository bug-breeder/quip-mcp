# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.4.0] - 2025-01-08

### Added
- **Comprehensive integration testing** against live Quip API
- **Full document lifecycle support**: create, read, edit, delete operations
- **Recent documents access** via `get_recent_threads` tool
- **Document editing capabilities** with append/prepend/replace operations
- **Document deletion** with confirmation safeguards
- **Robust API response handling** for complex nested Quip API structures
- **Cursor IDE development rules** for consistent code quality
- **Pre-commit workflow** with automated quality checks
- **Security scanning** with gosec integration
- **Comprehensive test coverage** (unit + integration tests)
- **Development automation** with enhanced Makefile targets

### Enhanced
- **API client reliability** with proper JSON unmarshaling for complex responses
- **Error handling** with detailed error messages and proper HTTP status codes
- **Markdown support** throughout all document operations
- **Documentation** with concise and comprehensive README
- **Testing infrastructure** with CI-ready test automation
- **Code quality** with linting, formatting, and security checks

### Fixed
- **Recent threads API** parsing of map-based response structures
- **Document creation** handling of nested `thread` objects in API responses
- **Document retrieval** proper extraction of HTML content from complex responses
- **Edit/Delete operations** correct API endpoints and parameter passing
- **Environment variable handling** in configuration and testing
- **API versioning support** for both v1 and v2 endpoints where appropriate

### Technical Improvements
- **RecentThreadsResponse/RecentThreadData** structs for complex API responses
- **Fallback unmarshaling** strategies for API compatibility
- **Form-encoded requests** with proper parameter mapping
- **Endpoint corrections** for edit and delete operations
- **Content type handling** for markdown vs HTML formats
- **Integration test cleanup** with automatic resource management

### Developer Experience
- **make pre-commit** - comprehensive quality checks before commits
- **make test-all** - unified unit and integration testing
- **make dev-setup** - one-command development environment setup
- **Enhanced error messages** for better debugging
- **Structured logging** throughout the application
- **Development documentation** with clear setup instructions

### Breaking Changes
- Unit tests updated to match corrected API behavior (no user-facing impact)
- Some internal API response structures changed (no user-facing impact)

---

## Previous Versions

For earlier versions, see the [GitHub releases page](https://github.com/bug-breeder/quip-mcp/releases).
