# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Security
- **BREAKING**: `Reset()` now requires `GO_ENV=test` and panics in production environments
  - Prevents accidental state corruption in production
  - Add `GO_ENV=test` when running tests that use `Reset()`
- Updated Go version requirement to 1.24.6+ (fixes 4 stdlib CVEs in Go 1.24.0)
- Added 5-second timeouts to all git command executions
  - Prevents hanging if git is unresponsive or unavailable
  - Uses `context.WithTimeout` for timeout enforcement
- Added HTTP request size limits to all handlers
  - `MaxBytesReader` with 1KB limit (defense in depth)
  - Applied to `Handler()` and `HealthHandler()`

### Fixed
- Fixed `loadedAt` race condition by setting timestamp earlier in initialization
  - Now set immediately after Info struct creation for proper immutability
- Improved Reset() documentation with strong security warnings
- Enhanced godoc comments for security-critical functions

### Added
- Security section in README.md covering:
  - Command execution safety (timeouts, injection prevention)
  - HTTP request protection (size limits)
  - Production safeguards (GO_ENV requirement)
  - Best practices for secure deployment

### Planned
- Core version information types and structures
- Multi-dimensional version manifest (YAML/JSON)
- Thread-safe singleton pattern
- Build-time injection via ldflags
- Repository layer (file, embed, buildinfo)
- Service layer with validation
- HTTP handlers (/version, /health)
- CLI tool (show, validate, bump, check commands)
- Comprehensive test suite
- Documentation and examples

## [0.1.0] - 2025-10-11

### Added
- Initial project structure
- Go module initialization
- Core dependencies (go-cuserr, zap, yaml.v3)
- Makefile with build, test, lint targets
- golangci-lint configuration
- CLAUDE.md for development guidance
- Project plan documentation
- MIT License
- VERSION file
- Directory structure for implementation

[Unreleased]: https://github.com/itsatony/go-version/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/itsatony/go-version/releases/tag/v0.1.0
