# go-version: Multi-Dimensional Versioning for Go Applications

**Repository:** `github.com/itsatony/go-version`  
**Status:** Design Phase  
**Owner:** vAudience.AI GmbH  
**License:** MIT (open source)  
**Motto:** Excellence. Always.

---

## Executive Summary

`go-version` is a Go library that solves multi-dimensional versioning for complex applications. It provides a unified approach to tracking, injecting, and exposing multiple version dimensions (project, database schemas, API versions, component versions) while remaining simple for basic use cases.

**Design Philosophy:** Zero-config simplicity for basic cases, opt-in complexity for advanced scenarios. No magic, no over-engineering.

---

## Problem Statement

Modern Go applications face versioning challenges:

1. **Multiple Version Dimensions**: Projects track project version, DB schemas, API versions, dependent service versions, configuration versions
2. **Fragmented Approaches**: Teams reimplement version handling in every repo (boilerplate, inconsistency)
3. **Build-time vs Runtime**: Confusion about what should be injected at build time vs loaded at runtime
4. **Observability**: Need to expose version info via HTTP endpoints, logs, and monitoring systems
5. **CI/CD Integration**: Build systems need programmatic access to version information
6. **Agentic Workflows**: Pre-commit version updates by agents require structured, parseable version files

**Current solutions are insufficient:**
- Simple `-ldflags` injection: only handles single string version
- `runtime/debug.BuildInfo`: limited to VCS info, no custom dimensions
- Existing packages: either abandoned, overcomplicated, or too opinionated

---

## Core Requirements

### Functional Requirements

1. **Multi-dimensional versioning** with support for:
   - Project/application version (semantic versioning)
   - Database schema versions (multiple databases)
   - API versions (multiple APIs)
   - Component/dependency versions
   - Custom version dimensions

2. **Multiple version sources:**
   - Git tags (source of truth for project version)
   - Build-time injection via `-ldflags`
   - Embedded version manifest file (YAML/JSON)
   - Runtime `debug.BuildInfo` (fallback/enrichment)

3. **Runtime access:**
   - Simple API to query any version dimension
   - Structured format for logging/observability
   - Type-safe access patterns

4. **HTTP exposure:**
   - Standard `/version` endpoint (JSON)
   - `/health` endpoint with version info
   - Prometheus-compatible metrics
   - Customizable endpoint paths

5. **CI/CD integration:**
   - CLI tool for querying versions from filesystem
   - Parseable output formats (JSON, YAML, env vars)
   - Exit codes for version validation

6. **Agentic workflow support:**
   - Standard version manifest format
   - Validation tools
   - Version bumping utilities
   - Pre-commit hook integration

### Non-Functional Requirements

1. **Zero external dependencies** for core library (stdlib only)
2. **Thread-safe:** All operations safe for concurrent access
3. **Performance:** Negligible runtime overhead (<1ms initialization)
4. **Memory footprint:** <10KB for typical version info
5. **Backwards compatibility:** Version manifest format versioned itself
6. **Testable:** Comprehensive test utilities included

---

## vAudience.AI Coding Standards Integration

### File Naming Convention

All files follow the pattern: `{project}.{type}.{module}.{framework}.go`

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

### Constants Management (NO STRING LITERALS)

**CRITICAL:** Every string literal must be a constant.

```go
// version.constants.go
package version

const (
	// Version manifest
	MANIFEST_VERSION             = "1.0"
	MANIFEST_FILENAME_YAML       = "versions.yaml"
	MANIFEST_FILENAME_JSON       = "versions.json"
	MANIFEST_DEFAULT_PROJECT     = "project"
	
	// Build info keys
	BUILD_INFO_KEY_VERSION       = "vcs.revision"
	BUILD_INFO_KEY_TIME          = "vcs.time"
	BUILD_INFO_KEY_MODIFIED      = "vcs.modified"
	
	// HTTP endpoints
	HTTP_PATH_VERSION            = "/version"
	HTTP_PATH_HEALTH             = "/health"
	HTTP_HEADER_CONTENT_TYPE     = "Content-Type"
	HTTP_CONTENT_TYPE_JSON       = "application/json"
	
	// Cache keys
	CACHE_KEY_VERSION_INFO       = "version:info"
	CACHE_KEY_BUILD_INFO         = "version:build"
	
	// Log messages (with placeholders)
	LOG_MSG_MANIFEST_LOADED      = "[%s.%s] Manifest loaded from (%s)"
	LOG_MSG_VERSION_INITIALIZED  = "[%s.%s] Version info initialized: project=%s version=%s"
	LOG_MSG_HTTP_HANDLER_CALLED  = "[%s.%s] Version endpoint called"
	LOG_MSG_VALIDATION_FAILED    = "[%s.%s] Validation failed: %s"
	
	// Error messages (with placeholders)
	ERR_MSG_MANIFEST_NOT_FOUND   = "version manifest not found at path: %s"
	ERR_MSG_MANIFEST_PARSE       = "failed to parse version manifest: %s"
	ERR_MSG_INVALID_FORMAT       = "invalid manifest format: %s"
	ERR_MSG_SCHEMA_VERSION_MISSING = "schema version missing for: %s"
	ERR_MSG_SCHEMA_VERSION_INVALID = "invalid schema version for %s: expected >= %s, got %s"
	ERR_MSG_BUILD_INFO_UNAVAILABLE = "build info unavailable"
	
	// Validation messages
	VALIDATE_MSG_REQUIRED_FIELD  = "required field missing: %s"
	VALIDATE_MSG_INVALID_VERSION = "invalid semantic version: %s"
	
	// Class names for logging
	CLASS_NAME_VERSION_SERVICE   = "VersionService"
	CLASS_NAME_VERSION_LOADER    = "VersionLoader"
	CLASS_NAME_HTTP_HANDLER      = "VersionHTTPHandler"
	CLASS_NAME_CLI               = "VersionCLI"
	
	// Method prefixes
	METHOD_PREFIX_GET            = "Get"
	METHOD_PREFIX_LOAD           = "Load"
	METHOD_PREFIX_VALIDATE       = "Validate"
	METHOD_PREFIX_SERVE          = "Serve"
	
	// Default values
	DEFAULT_GIT_COMMIT           = "dev"
	DEFAULT_GIT_TAG              = ""
	DEFAULT_GIT_TREE_STATE       = "clean"
	DEFAULT_BUILD_TIME           = "unknown"
	DEFAULT_BUILD_USER           = ""
	
	// CLI commands
	CLI_CMD_SHOW                 = "show"
	CLI_CMD_VALIDATE             = "validate"
	CLI_CMD_BUMP                 = "bump"
	CLI_CMD_CHECK                = "check"
	
	// CLI output formats
	CLI_FORMAT_JSON              = "json"
	CLI_FORMAT_YAML              = "yaml"
	CLI_FORMAT_ENV               = "env"
	CLI_FORMAT_HUMAN             = "human"
	
	// Bump types
	BUMP_TYPE_MAJOR              = "major"
	BUMP_TYPE_MINOR              = "minor"
	BUMP_TYPE_PATCH              = "patch"
	BUMP_TYPE_INCREMENT          = "increment"
)
```

### Error Handling with go-cuserr

```go
// version.errors.go
package version

import "github.com/itsatony/go-cuserr"

// Error categories
const (
	ERR_CATEGORY_MANIFEST    = "manifest"
	ERR_CATEGORY_VALIDATION  = "validation"
	ERR_CATEGORY_BUILD_INFO  = "buildinfo"
	ERR_CATEGORY_HTTP        = "http"
	ERR_CATEGORY_CLI         = "cli"
)

// Sentinel errors
var (
	ErrManifestNotFound = cuserr.New(
		ERR_CATEGORY_MANIFEST,
		"MANIFEST_NOT_FOUND",
		ERR_MSG_MANIFEST_NOT_FOUND,
	)
	
	ErrManifestParse = cuserr.New(
		ERR_CATEGORY_MANIFEST,
		"MANIFEST_PARSE",
		ERR_MSG_MANIFEST_PARSE,
	)
	
	ErrInvalidFormat = cuserr.New(
		ERR_CATEGORY_MANIFEST,
		"INVALID_FORMAT",
		ERR_MSG_INVALID_FORMAT,
	)
	
	ErrSchemaVersionMissing = cuserr.New(
		ERR_CATEGORY_VALIDATION,
		"SCHEMA_VERSION_MISSING",
		ERR_MSG_SCHEMA_VERSION_MISSING,
	)
	
	ErrSchemaVersionInvalid = cuserr.New(
		ERR_CATEGORY_VALIDATION,
		"SCHEMA_VERSION_INVALID",
		ERR_MSG_SCHEMA_VERSION_INVALID,
	)
	
	ErrBuildInfoUnavailable = cuserr.New(
		ERR_CATEGORY_BUILD_INFO,
		"BUILD_INFO_UNAVAILABLE",
		ERR_MSG_BUILD_INFO_UNAVAILABLE,
	)
)

// Error wrapping helpers
func WrapManifestError(err error, path string) error {
	return cuserr.Wrap(err, ERR_CATEGORY_MANIFEST, "manifest_error").
		WithField("path", path)
}

func WrapValidationError(err error, field string) error {
	return cuserr.Wrap(err, ERR_CATEGORY_VALIDATION, "validation_error").
		WithField("field", field)
}
```

### Thread-Safe Singleton Pattern

```go
// version.go
package version

import (
	"sync"
	"sync/atomic"
)

var (
	instance     atomic.Value // stores *Info
	initOnce     sync.Once
	initError    error
	mu           sync.RWMutex
)

// Get returns the singleton version info instance
// Thread-safe for concurrent access
func Get() (*Info, error) {
	// Fast path: already initialized
	if v := instance.Load(); v != nil {
		return v.(*Info), nil
	}
	
	// Slow path: initialize
	initOnce.Do(func() {
		var info *Info
		info, initError = loadVersionInfo()
		if initError == nil {
			instance.Store(info)
		}
	})
	
	if initError != nil {
		return nil, initError
	}
	
	return instance.Load().(*Info), nil
}

// MustGet returns version info or panics
// Use only in initialization code where failure is fatal
func MustGet() *Info {
	info, err := Get()
	if err != nil {
		panic(err)
	}
	return info
}
```

### Interface-First Design

```go
// version.go
package version

import "context"

// Loader defines how version information is loaded
type Loader interface {
	Load(ctx context.Context) (*Info, error)
}

// Validator defines version validation logic
type Validator interface {
	Validate(ctx context.Context, info *Info) error
}

// Repository defines version persistence operations
type Repository interface {
	Read(ctx context.Context, path string) (*Manifest, error)
	Write(ctx context.Context, path string, manifest *Manifest) error
}

// HTTPHandler defines HTTP endpoint behavior
type HTTPHandler interface {
	ServeVersion(w http.ResponseWriter, r *http.Request)
	ServeHealth(w http.ResponseWriter, r *http.Request)
}
```

### Structured Logging Integration

```go
// version.service.go
package version

import "go.uber.org/zap"

type VersionService struct {
	loader     Loader
	validators []Validator
	logger     *zap.Logger
	mu         sync.RWMutex
	cache      *Info
}

func (s *VersionService) GetInfo(ctx context.Context) (*Info, error) {
	methodName := METHOD_PREFIX_GET + "Info"
	
	s.mu.RLock()
	if s.cache != nil {
		s.mu.RUnlock()
		return s.cache, nil
	}
	s.mu.RUnlock()
	
	s.logger.Info(
		fmt.Sprintf(LOG_MSG_VERSION_INITIALIZED, CLASS_NAME_VERSION_SERVICE, methodName),
		zap.String("method", methodName),
	)
	
	info, err := s.loader.Load(ctx)
	if err != nil {
		s.logger.Error(
			fmt.Sprintf(LOG_MSG_VALIDATION_FAILED, CLASS_NAME_VERSION_SERVICE, methodName),
			zap.Error(err),
			zap.String("method", methodName),
		)
		return nil, WrapManifestError(err, "load")
	}
	
	// Validate
	for _, validator := range s.validators {
		if err := validator.Validate(ctx, info); err != nil {
			return nil, WrapValidationError(err, "validation")
		}
	}
	
	s.mu.Lock()
	s.cache = info
	s.mu.Unlock()
	
	s.logger.Info(
		fmt.Sprintf(LOG_MSG_VERSION_INITIALIZED, CLASS_NAME_VERSION_SERVICE, methodName),
		zap.String("project", info.Project.Name),
		zap.String("version", info.Project.Version),
		zap.String("method", methodName),
	)
	
	return info, nil
}
```

---

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

### Core Types

```go
// version.info.go
package version

import "time"

// Manifest represents the versions.yaml structure
type Manifest struct {
	ManifestVersion string                 `yaml:"manifest_version" json:"manifest_version"`
	Project         ProjectManifest        `yaml:"project" json:"project"`
	Schemas         map[string]string      `yaml:"schemas,omitempty" json:"schemas,omitempty"`
	APIs            map[string]string      `yaml:"apis,omitempty" json:"apis,omitempty"`
	Components      map[string]string      `yaml:"components,omitempty" json:"components,omitempty"`
	Custom          map[string]interface{} `yaml:"custom,omitempty" json:"custom,omitempty"`
}

type ProjectManifest struct {
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
}

// Info is the complete runtime version information
type Info struct {
	Project    ProjectVersion         `json:"project"`
	Git        GitInfo                `json:"git"`
	Build      BuildInfo              `json:"build"`
	Schemas    map[string]string      `json:"schemas,omitempty"`
	APIs       map[string]string      `json:"apis,omitempty"`
	Components map[string]string      `json:"components,omitempty"`
	Custom     map[string]interface{} `json:"custom,omitempty"`
	
	// Internal
	loadedAt time.Time
	mu       sync.RWMutex
}

type ProjectVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type GitInfo struct {
	Commit       string `json:"commit"`
	Tag          string `json:"tag,omitempty"`
	TreeState    string `json:"tree_state"` // clean, dirty
	CommitTime   string `json:"commit_time,omitempty"`
}

type BuildInfo struct {
	Time      string `json:"time"`
	User      string `json:"user,omitempty"`
	GoVersion string `json:"go_version"`
}

// Thread-safe accessors
func (i *Info) GetSchemaVersion(name string) (string, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	v, ok := i.Schemas[name]
	return v, ok
}

func (i *Info) GetAPIVersion(name string) (string, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	v, ok := i.APIs[name]
	return v, ok
}

func (i *Info) GetComponentVersion(name string) (string, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	v, ok := i.Components[name]
	return v, ok
}

// LogFields returns structured logging fields
func (i *Info) LogFields() []zap.Field {
	i.mu.RLock()
	defer i.mu.RUnlock()
	
	return []zap.Field{
		zap.String("project_name", i.Project.Name),
		zap.String("project_version", i.Project.Version),
		zap.String("git_commit", i.Git.Commit),
		zap.String("git_tag", i.Git.Tag),
		zap.String("build_time", i.Build.Time),
		zap.String("go_version", i.Build.GoVersion),
	}
}
```

### Configuration

```go
// version.config.go
package version

type Config struct {
	// Manifest loading
	ManifestPath     string
	DisableEmbedded  bool
	
	// Build info
	DisableGitInfo   bool
	DisableBuildInfo bool
	
	// Validation
	Validators       []Validator
	StrictMode       bool // Fail on any validation error
	
	// HTTP
	EnableHTTP       bool
	HTTPPath         string
	HealthPath       string
	HealthCheck      func() error
	
	// Logging
	Logger           *zap.Logger
	
	// Custom loaders
	CustomLoaders    []Loader
}

func DefaultConfig() *Config {
	return &Config{
		ManifestPath:     MANIFEST_FILENAME_YAML,
		DisableEmbedded:  false,
		DisableGitInfo:   false,
		DisableBuildInfo: false,
		StrictMode:       false,
		EnableHTTP:       true,
		HTTPPath:         HTTP_PATH_VERSION,
		HealthPath:       HTTP_PATH_HEALTH,
		Logger:           zap.NewNop(),
	}
}
```

---

## File Structure

```
github.com/itsatony/go-version/
├── README.md
├── LICENSE (MIT)
├── VERSION                           # Single source of truth
├── CHANGELOG.md
├── implementation_plan.md            # Detailed implementation plan
├── adrs.md                          # Architecture Decision Records
├── go.mod
├── go.sum
├── Makefile                         # Build, test, lint targets
├── .golangci.yml                    # Linter configuration
├── .env.example                     # Example environment variables
│
├── version.go                       # Public API, singleton
├── version.info.go                  # Info struct
├── version.config.go                # Configuration
├── version.constants.go             # ALL string constants
├── version.errors.go                # Sentinel errors
├── version.service.go               # Service layer
├── version.manifest.go              # Manifest parsing
├── version.buildinfo.go             # runtime/debug integration
├── version.ldflags.go               # Build-time injection points
├── version.http.go                  # HTTP handlers
├── version.http.middleware.go       # HTTP middleware
├── version.validate.go              # Validation logic
├── version.repository.file.go       # File repository
├── version.repository.embed.go      # Embedded repository
│
├── cmd/
│   └── go-version/                  # CLI tool
│       ├── main.go
│       ├── goversion.constants.go
│       ├── goversion.errors.go
│       ├── goversion.cmd.show.go
│       ├── goversion.cmd.validate.go
│       ├── goversion.cmd.bump.go
│       └── goversion.cmd.check.go
│
├── internal/
│   ├── semver/                      # Semantic version parsing
│   │   ├── semver.go
│   │   ├── semver.constants.go
│   │   ├── semver.errors.go
│   │   └── semver_test.go
│   └── git/                         # Git integration helpers
│       ├── git.go
│       ├── git.constants.go
│       ├── git.errors.go
│       └── git_test.go
│
├── testutil/                        # Testing utilities
│   ├── testutil.go
│   ├── testutil.mock.go
│   └── testutil.fixtures.go
│
├── examples/
│   ├── simple/                      # Minimal example
│   │   ├── main.go
│   │   └── versions.yaml
│   ├── complex/                     # Multi-dimensional example
│   │   ├── main.go
│   │   └── versions.yaml
│   └── microservice/                # Full-featured service
│       ├── main.go
│       ├── versions.yaml
│       └── docker-compose.yml
│
├── _templates/                      # Templates for users
│   ├── versions.yaml.tmpl
│   ├── pre-commit-hook.sh
│   └── Makefile.tmpl
│
└── docs/
    ├── getting-started.md
    ├── configuration.md
    ├── http-endpoints.md
    ├── cli-usage.md
    ├── migration-guide.md
    └── best-practices.md
```

---

## Excellence Gates

Every development cycle passes through seven gates:

### 🚀 GATE 1: CODE EXCELLENCE

**Pre-Development Planning:**
- Review existing code to avoid duplication/conflicts
- Design comprehensive tests covering edge cases
- Plan for thread safety and concurrent access
- Architecture validation against patterns
- Security considerations (no data leaks in errors)

**Checklist:**
- [ ] Strategic analysis complete
- [ ] Thread safety strategy defined
- [ ] Error handling patterns established
- [ ] Constants file created (NO string literals)
- [ ] Interface definitions finalized

### 🧪 GATE 2: TEST EXCELLENCE

**No shortcuts. Full functionality validation.**

```makefile
# Makefile
.PHONY: test test-unit test-integration test-race test-coverage

test: test-unit test-integration test-race

test-unit:
	@echo "Running unit tests..."
	go test -v -count=1 ./... -short

test-integration:
	@echo "Running integration tests..."
	go test -v -count=1 ./... -run Integration

test-race:
	@echo "Running race detection..."
	go test -race -count=1 ./...

test-coverage:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage: $$(go tool cover -func=coverage.out | grep total | awk '{print $$3}')"

test-bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...
```

**Test Requirements:**
- Unit tests with edge cases (>80% coverage)
- Integration tests validating end-to-end functionality
- Race condition detection: `go test -race`
- Benchmark tests with regression checks
- Mock tests for external dependencies (using testify/mock)

**Checklist:**
- [ ] All tests pass (no skipped tests)
- [ ] Race detector shows no issues
- [ ] Coverage meets minimum threshold
- [ ] Integration tests validate real workflows
- [ ] Benchmarks show acceptable performance

### 🔧 GATE 3: BUILD EXCELLENCE

```makefile
.PHONY: build build-cli build-all clean docker-build

VERSION ?= $(shell cat VERSION)
GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_TAG := $(shell git describe --tags --always)
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

LDFLAGS := -X github.com/itsatony/go-version.GitCommit=$(GIT_COMMIT)
LDFLAGS += -X github.com/itsatony/go-version.GitTag=$(GIT_TAG)
LDFLAGS += -X github.com/itsatony/go-version.BuildTime=$(BUILD_TIME)

build:
	go build -ldflags "$(LDFLAGS)" -o bin/example ./examples/simple

build-cli:
	go build -ldflags "$(LDFLAGS)" -o bin/go-version ./cmd/go-version

build-all: build build-cli

docker-build:
	docker build -t go-version:$(VERSION) .
	docker build -t go-version:latest .

clean:
	rm -rf bin/ coverage.out coverage.html
```

**Checklist:**
- [ ] All binaries build successfully
- [ ] Docker builds complete without errors
- [ ] No binaries committed to repository
- [ ] All files tracked in git (no orphaned files)
- [ ] Vulnerability scan passes (go list -m all | nancy sleuth)

### 📚 GATE 4: DOCUMENTATION EXCELLENCE

**Files to Update Each Cycle:**
- `README.md` - Quick start, features, installation
- `CHANGELOG.md` - What changed in this version
- `implementation_plan.md` - Current implementation status
- `adrs.md` - Architecture decisions made
- `docs/getting-started.md` - Tutorial walkthrough
- `docs/configuration.md` - All configuration options
- `docs/http-endpoints.md` - API documentation
- `docs/cli-usage.md` - CLI command reference
- `docs/migration-guide.md` - How to migrate from other solutions
- `docs/best-practices.md` - Usage recommendations

**Standards:**
- No broken links
- All code examples tested and functional
- API docs synchronized with actual code
- Markdown-lint compliance
- No claims about "world-class" or "production-ready" - just facts

**Checklist:**
- [ ] All documentation files reviewed
- [ ] Code examples tested
- [ ] Links verified
- [ ] API changes documented
- [ ] Migration guide updated

### 🔖 GATE 5: VERSION EXCELLENCE

**VERSION File (Single Source of Truth):**
```
2.1.3
```

**Semantic Versioning:**
- MAJOR: Breaking API changes
- MINOR: New features, backwards compatible
- PATCH: Bug fixes, backwards compatible

**CHANGELOG.md Format:**
```markdown
# Changelog

## [2.1.3] - 2025-10-11

### Added
- New validation rules for schema versions
- CLI command for version bumping

### Changed
- Improved error messages with more context
- Updated HTTP handler to include cache headers

### Fixed
- Race condition in singleton initialization
- Memory leak in file watcher

### Security
- Sanitized error messages to prevent data leaks
```

**Checklist:**
- [ ] VERSION file updated
- [ ] CHANGELOG.md comprehensive
- [ ] Git tag created (vX.Y.Z)
- [ ] Semantic versioning rules followed
- [ ] Breaking changes clearly documented

### 🛡️ GATE 6: SECURITY EXCELLENCE

**Security Requirements:**
- No secrets in code (API keys, passwords, tokens)
- No sensitive data in error messages
- Input validation on all external inputs
- No arbitrary code execution paths
- Thread-safe concurrent access

**Checklist:**
- [ ] No hardcoded secrets
- [ ] Error messages sanitized
- [ ] Input validation comprehensive
- [ ] Dependencies scanned for vulnerabilities
- [ ] Thread safety verified

### 🚀 GATE 7: FUNCTIONAL EXCELLENCE

**End-to-End Validation:**

```bash
# Test 1: Basic functionality
./bin/go-version show
# Expected: JSON output with version info

# Test 2: HTTP endpoints
cd examples/microservice
docker-compose up -d
curl http://localhost:8080/version | jq .
curl http://localhost:8080/health | jq .
docker-compose down

# Test 3: CLI validation
./bin/go-version validate examples/simple/versions.yaml
echo $? # Expected: 0

# Test 4: Version bumping
./bin/go-version bump --type patch --file examples/simple/versions.yaml
# Expected: versions.yaml updated

# Test 5: Integration test
cd examples/complex
go run main.go
# Expected: Application starts, logs version info, serves HTTP
```

**Checklist:**
- [ ] End-to-end workflow validated
- [ ] Container deployment tested
- [ ] API endpoints respond correctly
- [ ] CLI commands work as expected
- [ ] User journey completed successfully

---

## Testing Strategy

### Unit Tests (testify + testify/mock)

```go
// version_test.go
package version

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type VersionTestSuite struct {
	suite.Suite
}

func (s *VersionTestSuite) TestGet_ReturnsSameInstance() {
	info1, err1 := Get()
	require.NoError(s.T(), err1)
	require.NotNil(s.T(), info1)
	
	info2, err2 := Get()
	require.NoError(s.T(), err2)
	require.NotNil(s.T(), info2)
	
	// Verify singleton behavior
	assert.Same(s.T(), info1, info2)
}

func (s *VersionTestSuite) TestGetSchemaVersion_ThreadSafe() {
	info := &Info{
		Schemas: map[string]string{
			"db": "42",
		},
	}
	
	// Concurrent reads
	done := make(chan bool)
	for i := 0; i < 100; i++ {
		go func() {
			version, ok := info.GetSchemaVersion("db")
			assert.True(s.T(), ok)
			assert.Equal(s.T(), "42", version)
			done <- true
		}()
	}
	
	for i := 0; i < 100; i++ {
		<-done
	}
}

func TestVersionTestSuite(t *testing.T) {
	suite.Run(t, new(VersionTestSuite))
}
```

### Integration Tests (testcontainers-go)

```go
// version_integration_test.go
// +build integration

package version

import (
	"context"
	"net/http"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

func TestHTTPEndpoints_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	
	ctx := context.Background()
	
	// Start test server in container
	req := testcontainers.ContainerRequest{
		Image:        "go-version-example:latest",
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor:   wait.ForHTTP("/health").WithPort("8080"),
	}
	
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer container.Terminate(ctx)
	
	// Get container endpoint
	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "8080")
	require.NoError(t, err)
	
	baseURL := fmt.Sprintf("http://%s:%s", host, port.Port())
	
	// Test /version endpoint
	resp, err := http.Get(baseURL + "/version")
	require.NoError(t, err)
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, HTTP_CONTENT_TYPE_JSON, resp.Header.Get(HTTP_HEADER_CONTENT_TYPE))
	
	// Parse response
	var info Info
	err = json.NewDecoder(resp.Body).Decode(&info)
	require.NoError(t, err)
	
	assert.NotEmpty(t, info.Project.Version)
	assert.NotEmpty(t, info.Git.Commit)
}
```

### Mock Generation (mockery)

```bash
# Generate mocks
mockery --name=Loader --output=testutil --outpkg=testutil
mockery --name=Validator --output=testutil --outpkg=testutil
mockery --name=Repository --output=testutil --outpkg=testutil
```

```go
// Using generated mocks
func TestVersionService_GetInfo(t *testing.T) {
	mockLoader := new(testutil.MockLoader)
	mockValidator := new(testutil.MockValidator)
	
	expectedInfo := &Info{
		Project: ProjectVersion{Name: "test", Version: "1.0.0"},
	}
	
	mockLoader.On("Load", mock.Anything).Return(expectedInfo, nil)
	mockValidator.On("Validate", mock.Anything, expectedInfo).Return(nil)
	
	service := &VersionService{
		loader:     mockLoader,
		validators: []Validator{mockValidator},
		logger:     zap.NewNop(),
	}
	
	info, err := service.GetInfo(context.Background())
	
	assert.NoError(t, err)
	assert.Equal(t, expectedInfo, info)
	mockLoader.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
}
```

---

## Implementation Strategy

### Phase 1: Foundation (Week 1)

**Day 1-2: Core Types & Constants**
- [ ] Create all files with proper naming
- [ ] Define all constants (version.constants.go)
- [ ] Define all errors (version.errors.go)
- [ ] Define core types (Info, Manifest, etc.)
- [ ] Write unit tests for type methods

**Day 3-4: Repository Layer**
- [ ] Implement file repository
- [ ] Implement embedded repository
- [ ] Implement buildinfo integration
- [ ] Write unit tests with mocks
- [ ] Test concurrent access

**Day 5: Service Layer**
- [ ] Implement VersionService
- [ ] Implement loader composition
- [ ] Implement validation
- [ ] Write comprehensive unit tests
- [ ] Pass GATE 2 (Test Excellence)

### Phase 2: HTTP & CLI (Week 2)

**Day 1-2: HTTP Layer**
- [ ] Implement HTTP handlers
- [ ] Implement middleware (caching, compression)
- [ ] Add Prometheus metrics support
- [ ] Write integration tests
- [ ] Test with real HTTP server

**Day 3-5: CLI Tool**
- [ ] Implement `show` command
- [ ] Implement `validate` command
- [ ] Implement `bump` command
- [ ] Implement `check` command
- [ ] Write CLI integration tests
- [ ] Add shell completion

### Phase 3: Polish & Release (Week 3)

**Day 1-2: Documentation**
- [ ] Write comprehensive README
- [ ] Create all docs/ files
- [ ] Write examples
- [ ] Create templates
- [ ] Pass GATE 4 (Documentation Excellence)

**Day 3-4: Final Testing**
- [ ] Run all excellence gates
- [ ] Fix any issues found
- [ ] Performance benchmarking
- [ ] Security audit
- [ ] Pass all gates

**Day 5: Release**
- [ ] Tag v1.0.0
- [ ] Publish to GitHub
- [ ] Announce internally
- [ ] Create migration guides
- [ ] Deploy to first vAI project

---

## API Examples

### Simple Usage

```go
package main

import (
	"log"
	"github.com/itsatony/go-version"
)

func main() {
	// Zero-config usage
	v := version.MustGet()
	
	log.Printf("Starting %s version %s (commit: %s)", 
		v.Project.Name, 
		v.Project.Version, 
		v.Git.Commit)
}
```

### HTTP Integration

```go
package main

import (
	"net/http"
	"github.com/itsatony/go-version"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	
	// Initialize with config
	config := version.DefaultConfig()
	config.Logger = logger
	config.EnableHTTP = true
	config.HealthCheck = func() error {
		// Your health checks
		return nil
	}
	
	err := version.Initialize(config)
	if err != nil {
		logger.Fatal("Failed to initialize version", zap.Error(err))
	}
	
	// Register handlers
	http.HandleFunc("/version", version.HTTPHandler())
	http.HandleFunc("/health", version.HealthHandler())
	
	logger.Info("Starting server", zap.String("addr", ":8080"))
	http.ListenAndServe(":8080", nil)
}
```

### Advanced Usage with Validation

```go
package main

import (
	"context"
	"log"
	"github.com/itsatony/go-version"
)

func main() {
	// Custom configuration
	config := &version.Config{
		ManifestPath: "./versions.yaml",
		StrictMode:   true,
		Validators: []version.Validator{
			version.NewSchemaValidator("postgres_nexus", "45"), // Require minimum version
			version.NewAPIValidator("rest_v2", "2.0.0"),
		},
	}
	
	err := version.Initialize(config)
	if err != nil {
		log.Fatalf("Version validation failed: %v", err)
	}
	
	v := version.MustGet()
	
	// Check specific versions
	if schemaVer, ok := v.GetSchemaVersion("postgres_nexus"); ok {
		log.Printf("Database schema version: %s", schemaVer)
	}
	
	// Use structured logging fields
	logger := zap.NewProduction()
	logger = logger.With(v.LogFields()...)
	logger.Info("Application started with version info")
}
```

---

## Build Integration

### Makefile

```makefile
.PHONY: all build test lint docker-build release

VERSION := $(shell cat VERSION)
GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_TAG := $(shell git describe --tags --always)
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

LDFLAGS := -s -w
LDFLAGS += -X github.com/itsatony/go-version.GitCommit=$(GIT_COMMIT)
LDFLAGS += -X github.com/itsatony/go-version.GitTag=$(GIT_TAG)
LDFLAGS += -X github.com/itsatony/go-version.BuildTime=$(BUILD_TIME)

all: lint test build

build:
	@echo "Building version $(VERSION)..."
	go build -ldflags "$(LDFLAGS)" -o bin/go-version ./cmd/go-version

test:
	@echo "Running all tests..."
	go test -v -race -count=1 ./...

test-coverage:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage: $$(go tool cover -func=coverage.out | grep total | awk '{print $$3}')"

lint:
	@echo "Running linters..."
	golangci-lint run ./...

docker-build:
	@echo "Building Docker image..."
	docker build -t go-version:$(VERSION) .
	docker tag go-version:$(VERSION) go-version:latest

release: all docker-build
	@echo "Creating release $(VERSION)..."
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)

clean:
	rm -rf bin/ coverage.out coverage.html
```

### GitHub Actions CI/CD

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
      
      - name: Run tests
        run: make test
      
      - name: Run race detector
        run: go test -race ./...
      
      - name: Generate coverage
        run: make test-coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
  
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: golangci/golangci-lint-action@v3
        with:
          version: latest
  
  build:
    runs-on: ubuntu-latest
    needs: [test, lint]
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Build
        run: make build
      
      - name: Docker build
        run: make docker-build
```

---

## Migration Path for vAI Projects

### Step 1: Add Dependency

```bash
go get github.com/itsatony/go-version@latest
```

### Step 2: Create versions.yaml

```yaml
# versions.yaml
manifest_version: "1.0"

project:
  name: "nexus"
  version: "2.1.3"

schemas:
  postgres_nexus: "47"
  postgres_analytics: "12"
  redis_cache: "3"

apis:
  rest_v1: "1.15.0"
  rest_v2: "2.3.0"
  grpc: "1.2.0"

components:
  aigentchat: "3.4.1"
  hyperrag: "2.0.5"
  aigentflow: "1.8.2"

custom:
  config_schema: "v5"
```

### Step 3: Embed in Binary

```go
package main

import (
	_ "embed"
	"github.com/itsatony/go-version"
)

//go:embed versions.yaml
var versionsYAML []byte

func init() {
	// Set embedded content before Initialize
	version.SetEmbeddedManifest(versionsYAML)
}
```

### Step 4: Replace Custom Version Code

**Before:**
```go
var (
	Version   = "dev"
	GitCommit = "unknown"
)

func main() {
	log.Printf("Starting version %s", Version)
}
```

**After:**
```go
import "github.com/itsatony/go-version"

func main() {
	v := version.MustGet()
	log.Printf("Starting %s version %s (commit: %s)", 
		v.Project.Name, v.Project.Version, v.Git.Commit)
}
```

### Step 5: Update Makefile

```makefile
LDFLAGS := -X github.com/itsatony/go-version.GitCommit=$(GIT_COMMIT)
LDFLAGS += -X github.com/itsatony/go-version.GitTag=$(GIT_TAG)
LDFLAGS += -X github.com/itsatony/go-version.BuildTime=$(BUILD_TIME)

build:
	go build -ldflags "$(LDFLAGS)" -o bin/nexus ./cmd/nexus
```

### Step 6: Update HTTP Endpoints

**Before:**
```go
http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"version": Version})
})
```

**After:**
```go
import "github.com/itsatony/go-version"

http.HandleFunc("/version", version.HTTPHandler())
http.HandleFunc("/health", version.HealthHandler())
```

**Estimated effort per project: 2-4 hours**

---

## Success Criteria

1. **Adoption:** Used in all vAI Go projects within 3 months
2. **Reduced boilerplate:** 80% less custom version code in new projects
3. **Consistency:** Uniform /version endpoints across all services
4. **CI/CD integration:** Version checks in all pipelines
5. **Test coverage:** >85% for core library
6. **Zero production bugs:** No critical bugs in first 6 months
7. **Community:** 100+ GitHub stars within 6 months (if open-sourced)

---

## Open Questions

1. **Manifest format:** YAML only or both YAML and JSON?
   - **Recommendation:** Both, detect by extension
   
2. **CLI binary name:** `go-version` or `gover`?
   - **Recommendation:** `go-version` for clarity
   
3. **Prometheus metrics:** Opt-in or always included?
   - **Recommendation:** Opt-in to avoid dependency
   
4. **Minimum Go version:** 1.21+ or 1.18+?
   - **Recommendation:** 1.21+ for latest stdlib features
   
5. **Default logger:** zap, slog, or interface?
   - **Recommendation:** Interface with zap example

---

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|-----------|
| Over-engineering | Medium | High | Strict scope, phased releases |
| Adoption resistance | Low | Medium | Excellent docs, easy migration |
| Breaking changes needed | Low | High | Careful API design, versioned manifest |
| Performance issues | Very Low | Low | Benchmark suite, caching |
| Thread safety bugs | Low | High | Comprehensive race testing |
| Test maintenance burden | Medium | Medium | Good test utilities, clear patterns |

---

## Conclusion

`go-version` solves multi-dimensional versioning for vAudience.AI and the broader Go community. It follows vAI's strict coding standards, passes all excellence gates, and integrates seamlessly with existing workflows.

**Development Estimate:** 3 weeks (1 senior engineer)  
**Maintenance Burden:** Low (stable API, minimal dependencies)  
**Strategic Value:** High (used across all vAI projects, potential open-source visibility)

**Recommendation:** Proceed with implementation following this specification.

---

**Excellence. Always.**
