# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`go-version` is a multi-dimensional versioning library for Go applications that tracks project versions, database schemas, API versions, and component versions. It follows vAudience.AI's strict coding standards with a focus on excellence, zero dependencies for core functionality, and thread-safe concurrent access.

**Motto:** Excellence. Always.

## Build Commands

```bash
# Initialize and install dependencies
go mod download

# Run all tests (unit + integration + race detection)
make test

# Run only unit tests
make test-unit

# Run integration tests
make test-integration

# Run race detector
make test-race

# Generate coverage report
make test-coverage

# Run linter
make lint

# Build CLI tool
make build-cli

# Build all binaries
make build-all

# Clean build artifacts
make clean

# Docker build
make docker-build
```

## Architecture

### Three-Layer Design

```
┌─────────────────────────────────────────────┐
│   Presentation Layer (HTTP, CLI, Metrics)   │
├─────────────────────────────────────────────┤
│        Service Layer (Business Logic)       │
├─────────────────────────────────────────────┤
│  Repository Layer (File, Embed, BuildInfo)  │
└─────────────────────────────────────────────┘
```

- **Repository Layer**: Handles loading version data from files, embedded manifests, and runtime build info
- **Service Layer**: Orchestrates loaders, validators, caching, and thread-safe access
- **Presentation Layer**: Exposes version info via HTTP endpoints, CLI commands, and structured logging

## Critical Coding Standards

### 1. File Naming Convention (NON-NEGOTIABLE)

**Pattern:** `{project}.{type}.{module}.{framework}.go`

```
version.go                    # Core public API
version.info.go              # Info struct and methods
version.manifest.go          # YAML/JSON parsing
version.buildinfo.go         # runtime/debug integration
version.ldflags.go           # Build-time injection points
version.http.go              # HTTP handlers
version.http.middleware.go   # HTTP middleware
version.validate.go          # Validation logic
version.constants.go         # ALL string constants
version.errors.go            # Sentinel error definitions
version.service.go           # Service layer
version.repository.file.go   # File-based repository
```

### 2. Constants Management (CRITICAL)

**EVERY string literal MUST be a constant.** No exceptions.

All constants live in `version.constants.go`:
- Manifest filenames, keys, and version
- Build info keys
- HTTP paths, headers, content types
- Cache keys
- Log message templates (with placeholders like `"[%s.%s] Message here"`)
- Error message templates
- Class names for logging (e.g., `CLASS_NAME_VERSION_SERVICE`)
- Method prefixes (e.g., `METHOD_PREFIX_GET`)
- Default values
- CLI commands and output formats

When adding ANY new string literal, add it to constants first.

### 3. Error Handling with go-cuserr

All errors defined in `version.errors.go` using go-cuserr:

```go
var (
    ErrManifestNotFound = cuserr.New(
        ERR_CATEGORY_MANIFEST,
        "MANIFEST_NOT_FOUND",
        ERR_MSG_MANIFEST_NOT_FOUND,
    )
)

// Use wrapper functions for context
func WrapManifestError(err error, path string) error {
    return cuserr.Wrap(err, ERR_CATEGORY_MANIFEST, "manifest_error").
        WithField("path", path)
}
```

Error categories: `manifest`, `validation`, `buildinfo`, `http`, `cli`

### 4. Thread-Safe Singleton Pattern

```go
var (
    instance     atomic.Value  // stores *Info
    initOnce     sync.Once
    initError    error
    mu           sync.RWMutex
)

// Fast path with atomic.Value for reads
// Slow path with sync.Once for initialization
```

All shared state must use `sync.RWMutex` for thread safety. Test with `-race` flag.

### 5. Interface-First Design

Define interfaces before implementations:
- `Loader` - loads version information
- `Validator` - validates version constraints
- `Repository` - persists/reads manifests
- `HTTPHandler` - HTTP endpoint behavior

### 6. Structured Logging

All logging uses zap with structured fields:

```go
s.logger.Info(
    fmt.Sprintf(LOG_MSG_VERSION_INITIALIZED, CLASS_NAME_VERSION_SERVICE, methodName),
    zap.String("method", methodName),
    zap.String("project", info.Project.Name),
    zap.String("version", info.Project.Version),
)
```

Log message format: `[ClassName.MethodName] Message content`

## Excellence Gates

Every development cycle must pass 7 gates:

1. **CODE EXCELLENCE**: Strategic planning, thread safety, no string literals
2. **TEST EXCELLENCE**: >80% coverage, race detection, benchmarks
3. **BUILD EXCELLENCE**: All binaries build, Docker builds, no vulnerabilities
4. **DOCUMENTATION EXCELLENCE**: All docs updated, code examples tested
5. **VERSION EXCELLENCE**: VERSION file updated, CHANGELOG.md comprehensive
6. **SECURITY EXCELLENCE**: No secrets, sanitized errors, input validation
7. **FUNCTIONAL EXCELLENCE**: End-to-end workflows validated

**Do not skip gates.** Each gate has specific checklists in projectplan.md.

## Key Patterns

### Thread-Safe Accessors

```go
func (i *Info) GetSchemaVersion(name string) (string, bool) {
    i.mu.RLock()
    defer i.mu.RUnlock()
    v, ok := i.Schemas[name]
    return v, ok
}
```

Always use RLock/RUnlock for reads, Lock/Unlock for writes.

### Build-Time Injection

Variables in `version.ldflags.go` can be injected via `-ldflags`:

```bash
-X github.com/itsatony/go-version.GitCommit=$(GIT_COMMIT)
-X github.com/itsatony/go-version.GitTag=$(GIT_TAG)
-X github.com/itsatony/go-version.BuildTime=$(BUILD_TIME)
```

### Version Manifest Format

```yaml
manifest_version: "1.0"
project:
  name: "myapp"
  version: "1.0.0"
schemas:
  postgres_main: "42"
apis:
  rest_v1: "1.0.0"
components:
  service_a: "2.3.1"
custom:
  config_schema: "v5"
```

## Testing Requirements

- **Unit tests**: Use testify/suite, testify/assert, testify/mock
- **Integration tests**: Use testcontainers-go for real containers
- **Race detection**: Always run `make test-race` before commits
- **Coverage**: Minimum 80% for new code, 85% target overall
- **Benchmarks**: Add for any performance-critical paths

Mock generation with mockery:
```bash
mockery --name=Loader --output=testutil --outpkg=testutil
```

## Common Tasks

### Adding a New Version Dimension

1. Update `Manifest` struct in `version.info.go`
2. Update `Info` struct with thread-safe accessor methods
3. Add constants for any new keys/messages
4. Update validation logic if needed
5. Update HTTP response serialization
6. Add tests for new dimension
7. Update documentation and examples

### Adding a New Validator

1. Implement `Validator` interface
2. Add error constants and messages
3. Write unit tests with mocks
4. Document in configuration guide
5. Add example usage

### Adding a New CLI Command

1. Create `goversion.cmd.{command}.go`
2. Add command constants (name, flags, help text)
3. Implement command logic
4. Add tests
5. Update CLI documentation
6. Add shell completion support

## Dependencies

Core library has ZERO external dependencies (stdlib only). Optional dependencies:
- `github.com/itsatony/go-cuserr` - Custom error handling
- `go.uber.org/zap` - Structured logging (interface-based, swappable)
- `gopkg.in/yaml.v3` - YAML parsing (manifest loading)

Test dependencies:
- `github.com/stretchr/testify` - Test framework
- `github.com/testcontainers/testcontainers-go` - Integration tests

## Project Status

See projectplan.md for:
- Complete implementation plan (3-week timeline)
- Detailed architecture decisions
- API examples and usage patterns
- Migration guides for existing projects
- Risk assessment and mitigation strategies

## Code Review Checklist

Before committing:
- [ ] All string literals are constants
- [ ] Errors use go-cuserr with proper categories
- [ ] Thread safety verified (run with -race)
- [ ] Tests added/updated (unit + integration)
- [ ] Documentation updated if API changed
- [ ] File naming follows convention
- [ ] Logging uses structured fields
- [ ] No hardcoded values
- [ ] Interface-first design maintained
