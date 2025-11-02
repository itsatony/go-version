# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2025-11-02

### Breaking Changes
- **CRITICAL**: Fixed YAML struct tags to use lowercase keys (aligning with YAML conventions)
  - **Migration required**: Update your `versions.yaml` files to use lowercase keys
  - Change `Schemas:` → `schemas:`
  - Change `APIs:` → `apis:`
  - Change `Components:` → `components:`
  - Change `Custom:` → `custom:`
  - This fixes a bug where the implementation used non-standard capitalized YAML keys
  - Template file (_templates/versions.yaml.tmpl) and documentation were already using correct lowercase format
  - **Why this change**: YAML convention uses lowercase/snake_case keys (like `manifest_version` and `project`), not mixed case
  - **Impact**: All existing manifest files must be updated to use lowercase keys
- Info struct map fields now unexported (use GetSchemas(), GetAPIs(), GetComponents(), GetCustom() instead)
  - True immutability with defensive copies
  - Prevents external mutation of internal state
- `Reset()` now uses `testing.Testing()` instead of GO_ENV (automatically detects test environment)

### Security
- Git binary PATH validation (whitelists /usr/bin, /usr/local/bin, /opt/homebrew/bin)
- Command injection protection with constant arguments
- Commit hash validation (hexadecimal format checks)
- Git command timeouts (5s) to prevent hangs
- Updated Go version requirement to 1.24.6+ (fixes 4 stdlib CVEs in Go 1.24.0)
- HTTP request size limits using `MaxBytesReader` (1KB limit, defense in depth)
- Production safeguards: `Reset()` panics if called outside test environment

### Added
- Core version information types and structures (Info, Manifest, ProjectVersion, GitInfo, BuildInfo)
- Multi-dimensional version manifest support (project, schemas, APIs, components, custom)
- Thread-safe singleton pattern with lock-free atomic reads
- Build-time injection via ldflags (GitCommit, GitTag, BuildTime, BuildUser)
- File-based manifest loading (versions.yaml)
- Embedded manifest support (go:embed)
- Build info extraction from runtime/debug
- Git info extraction (commit, tag, tree state)
- Context propagation to validators (WithContext option)
- Public semantic version API (ParseSemVer, MustParseSemVer, CompareVersions, IsNewerVersion)
- SemVer struct with full comparison methods (Compare, LessThan, GreaterThan, Equal, etc.)
- IsInitialized() to check state without triggering auto-init
- WithStrictMode() for production validation requirements
- Defensive copy getters (GetSchemas(), GetAPIs(), GetComponents(), GetCustom())
- Custom JSON marshaling for unexported fields
- Validation framework with Validator interface
- Built-in validators (NewSchemaValidator, NewAPIValidator, NewComponentValidator)
- HTTP handlers (/version endpoint, /health endpoint)
- HTTP middleware (adds version headers to responses)
- CLI tool with multiple output formats (full, compact, JSON, schemas, apis, components, git, build)
- Structured logging support with LogFields() for zap
- Comprehensive test suite (140+ tests, 90.5% coverage)
- Race detection verified (all tests pass with -race flag)
- 11 benchmark tests for performance tracking
- Security tests for git validation and command injection protection
- Complete documentation in README.md with examples
- YAML manifest template (_templates/versions.yaml.tmpl)
- Working examples (simple and advanced)

### Fixed
- YAML struct tags now use lowercase keys (standards compliance)
- All example and test YAML files updated to use lowercase keys
- Thread safety with true immutability (defensive copies)
- Reset() now properly detects test environment automatically

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

[Unreleased]: https://github.com/itsatony/go-version/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/itsatony/go-version/compare/v0.1.0...v1.0.0
[0.1.0]: https://github.com/itsatony/go-version/releases/tag/v0.1.0
