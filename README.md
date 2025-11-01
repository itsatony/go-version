# go-version

Multi-dimensional versioning for Go applications. Track project versions, database schemas, API versions, and component versions in one unified system.

[![Go Reference](https://pkg.go.dev/badge/github.com/itsatony/go-version.svg)](https://pkg.go.dev/github.com/itsatony/go-version)
[![Go Report Card](https://goreportcard.com/badge/github.com/itsatony/go-version)](https://goreportcard.com/report/github.com/itsatony/go-version)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Coverage](https://img.shields.io/badge/coverage-86.3%25-brightgreen.svg)](https://github.com/itsatony/go-version)

## Features

- 🚀 **Zero-config usage** - Works out of the box with sensible defaults
- 📦 **Multi-dimensional versioning** - Track project, schemas, APIs, and components
- 🔒 **Thread-safe** - Concurrent access without locks (immutable design with defensive copies)
- ✅ **Validation** - Enforce minimum version requirements with context support
- 🌐 **HTTP endpoints** - Ready-to-use JSON endpoints and health checks
- 🔧 **Build-time injection** - Inject git and build metadata via ldflags
- 🖥️ **CLI tool** - Command-line interface for version queries and CI/CD
- 📊 **Structured logging** - First-class support for zap and similar loggers
- 🔢 **Semantic versioning** - Built-in semver parsing and comparison utilities
- 🔐 **Security hardened** - Git binary validation, command injection protection
- 🎯 **Production-ready** - 86%+ test coverage with race detection
- 📝 **Well-documented** - Comprehensive examples and API documentation

## Installation

```bash
go get github.com/itsatony/go-version
```

## Quick Start

### Zero-Config Usage

```go
package main

import (
    "log"
    "github.com/itsatony/go-version"
)

func main() {
    // Get version info with zero configuration
    info := version.MustGet()

    log.Printf("Starting %s v%s", info.Project.Name, info.Project.Version)
    log.Printf("Git commit: %s", info.Git.Commit)
}
```

### With Configuration File

Create `versions.yaml` in your project root:

```yaml
manifest_version: "1.0"
project:
  name: "my-app"
  version: "1.2.3"
schemas:
  postgres_main: "45"
  redis_cache: "3"
apis:
  rest_v1: "1.15.0"
  grpc: "1.2.0"
components:
  auth_service: "2.1.0"
```

The library will auto-discover and load this file.

### With Validation

```go
package main

import (
    "log"
    "github.com/itsatony/go-version"
)

func main() {
    // Initialize with version validation
    err := version.Initialize(
        version.WithManifestPath("versions.yaml"),
        version.WithValidators(
            version.NewSchemaValidator("postgres_main", "45"),
            version.NewAPIValidator("rest_v1", "1.10.0"),
        ),
    )
    if err != nil {
        log.Fatal("Version requirements not met:", err)
    }

    info := version.MustGet()
    log.Printf("Starting %s v%s", info.Project.Name, info.Project.Version)
}
```

## Usage Examples

### HTTP Endpoints

```go
package main

import (
    "net/http"
    "github.com/itsatony/go-version"
)

func main() {
    mux := http.NewServeMux()

    // Version endpoint (returns JSON)
    mux.Handle("/version", version.Handler())

    // Health check endpoint
    mux.Handle("/health", version.HealthHandler())

    // Your API endpoints
    mux.HandleFunc("/api/users", handleUsers)

    // Wrap with middleware to add version headers
    handler := version.Middleware(mux)

    http.ListenAndServe(":8080", handler)
}
```

Test the endpoints:
```bash
curl http://localhost:8080/version
curl http://localhost:8080/health
curl -I http://localhost:8080/api/users  # Check X-App-Version header
```

### Embedded Manifest

```go
package main

import (
    _ "embed"
    "github.com/itsatony/go-version"
)

//go:embed versions.yaml
var versionsYAML []byte

func main() {
    // Initialize with embedded manifest
    err := version.Initialize(
        version.WithEmbedded(versionsYAML),
        version.WithGitInfo(),
        version.WithBuildInfo(),
    )
    if err != nil {
        panic(err)
    }

    info := version.MustGet()
    // ...
}
```

### Version Checking

```go
package main

import (
    "log"
    "github.com/itsatony/go-version"
)

func main() {
    info := version.MustGet()

    // Check if schema version is available
    schemaVer, ok := info.GetSchemaVersion("postgres_main")
    if !ok {
        log.Fatal("PostgreSQL schema version not defined")
    }
    log.Printf("Database schema version: %s", schemaVer)

    // Check API version
    apiVer, ok := info.GetAPIVersion("rest_v1")
    if ok {
        log.Printf("REST API version: %s", apiVer)
    }

    // Access custom metadata
    if env, ok := info.Custom["environment"].(string); ok {
        log.Printf("Environment: %s", env)
    }
}
```

### Structured Logging (Zap)

```go
package main

import (
    "go.uber.org/zap"
    "github.com/itsatony/go-version"
)

func main() {
    info := version.MustGet()

    // Create logger with version fields
    logger := zap.NewProduction()
    logger = logger.With(info.LogFields()...)

    // All log entries now include version information
    logger.Info("Application started")
    // Output includes: version, git_commit, build_time, etc.
}
```

### Custom Validators

```go
package main

import (
    "context"
    "fmt"
    "github.com/itsatony/go-version"
)

func main() {
    err := version.Initialize(
        version.WithManifestPath("versions.yaml"),
        version.WithValidators(
            // Built-in validators
            version.NewSchemaValidator("postgres_main", "45"),

            // Custom validator using ValidatorFunc
            version.ValidatorFunc(func(ctx context.Context, info *version.Info) error {
                // Ensure production builds have a git tag
                if env, ok := info.Custom["environment"].(string); ok && env == "production" {
                    if info.Git.Tag == "" {
                        return fmt.Errorf("production builds must have a git tag")
                    }
                }
                return nil
            }),
        ),
    )
    if err != nil {
        panic(err)
    }
}
```

### Semantic Version Comparison

```go
package main

import (
    "fmt"
    "log"
    "github.com/itsatony/go-version"
)

func main() {
    // Parse semantic versions
    v1, err := version.ParseSemVer("1.2.3")
    if err != nil {
        log.Fatal(err)
    }

    v2 := version.MustParseSemVer("2.0.0")

    // Compare versions
    if v1.LessThan(v2) {
        fmt.Println("v1 is older than v2")
    }

    // Convenience functions for string comparison
    isNewer, err := version.IsNewerVersion("2.1.0", "2.0.0")
    if err != nil {
        log.Fatal(err)
    }
    if isNewer {
        fmt.Println("Upgrade available!")
    }

    // Compare with current project version
    info := version.MustGet()
    current := version.MustParseSemVer(info.Project.Version)
    required := version.MustParseSemVer("1.0.0")

    if current.GreaterThanOrEqual(required) {
        fmt.Println("Version requirements met")
    }
}
```

### Context Propagation

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/itsatony/go-version"
)

func main() {
    // Create context with timeout for validation
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err := version.Initialize(
        version.WithContext(ctx),
        version.WithManifestPath("versions.yaml"),
        version.WithValidators(
            version.ValidatorFunc(func(ctx context.Context, info *version.Info) error {
                // Check for context cancellation
                select {
                case <-ctx.Done():
                    return ctx.Err()
                default:
                }

                // Perform validation with context
                // Can use ctx for tracing, logging, etc.
                if info.Project.Version == "" {
                    return fmt.Errorf("version required")
                }
                return nil
            }),
        ),
    )
    if err != nil {
        panic(err)
    }
}
```

### Non-Singleton Usage

```go
package main

import (
    "log"
    "github.com/itsatony/go-version"
)

func main() {
    // Create independent version info (not using singleton)
    info, err := version.New(
        version.WithManifestPath("service-a-versions.yaml"),
    )
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Service A version: %s", info.Project.Version)

    // Singleton is unaffected
    globalInfo := version.MustGet()
    log.Printf("Global version: %s", globalInfo.Project.Version)
}
```

## Command-Line Tool

A CLI tool is available for displaying version information from the terminal.

### Installation

```bash
go install github.com/itsatony/go-version/cmd/go-version@latest
```

### Usage

```bash
# Show all version information
go-version

# JSON output
go-version -json

# Compact format
go-version -compact

# Custom manifest
go-version -manifest ./config/versions.yaml

# Show only schemas
go-version -schemas

# Show only git info
go-version -git
```

### Examples

Show version in CI/CD:
```bash
VERSION=$(go-version -compact)
echo "Deploying: $VERSION"
```

Extract specific fields:
```bash
go-version -json | jq -r '.project.version'
go-version -json | jq -r '.git.commit'
```

See [cmd/go-version/README.md](cmd/go-version/README.md) for complete CLI documentation.

## Build-Time Injection

Inject git and build metadata at compile time using ldflags:

### Basic Injection

```bash
go build -ldflags="\
  -X github.com/itsatony/go-version.GitCommit=$(git rev-parse HEAD) \
  -X github.com/itsatony/go-version.GitTag=$(git describe --tags --always) \
  -X github.com/itsatony/go-version.BuildTime=$(date -u '+%Y-%m-%dT%H:%M:%SZ') \
  -X github.com/itsatony/go-version.BuildUser=$(whoami)"
```

### Makefile Example

```makefile
VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
BUILD_USER := $(shell whoami)

LDFLAGS := -X github.com/itsatony/go-version.GitCommit=$(COMMIT)
LDFLAGS += -X github.com/itsatony/go-version.GitTag=$(VERSION)
LDFLAGS += -X github.com/itsatony/go-version.BuildTime=$(BUILD_TIME)
LDFLAGS += -X github.com/itsatony/go-version.BuildUser=$(BUILD_USER)

build:
	go build -ldflags="$(LDFLAGS)" -o myapp ./cmd/myapp

.PHONY: build
```

### GitHub Actions Example

```yaml
- name: Build with version info
  run: |
    go build -ldflags="\
      -X github.com/itsatony/go-version.GitCommit=${{ github.sha }} \
      -X github.com/itsatony/go-version.GitTag=${{ github.ref_name }} \
      -X github.com/itsatony/go-version.BuildTime=$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
```

## API Overview

### Core Functions

- `Initialize(opts ...Option) error` - Configure and initialize version singleton
- `Get() (*Info, error)` - Get version info (auto-initializes with defaults)
- `MustGet() *Info` - Get version info or panic
- `New(opts ...Option) (*Info, error)` - Create non-singleton instance
- `IsInitialized() bool` - Check if singleton is initialized (no auto-init)

### Options

- `WithManifestPath(path string)` - Load from custom file path
- `WithEmbedded(data []byte)` - Use embedded manifest data
- `WithGitInfo()` - Include git information (default: true)
- `WithoutGitInfo()` - Disable git information
- `WithBuildInfo()` - Include build information (default: true)
- `WithoutBuildInfo()` - Disable build information
- `WithValidators(validators ...Validator)` - Add version validators
- `WithContext(ctx context.Context)` - Set context for validation (supports cancellation/tracing)
- `WithStrictMode()` - Require manifest file and strict validation

### Validators

- `NewSchemaValidator(name, minVersion string)` - Validate schema version
- `NewAPIValidator(name, minVersion string)` - Validate API version
- `NewComponentValidator(name, minVersion string)` - Validate component version
- `ValidatorFunc` - Create custom validator from function

### Semantic Versioning

- `ParseSemVer(s string) (*SemVer, error)` - Parse semantic version string
- `MustParseSemVer(s string) *SemVer` - Parse or panic
- `CompareVersions(v1, v2 string) (int, error)` - Compare two version strings
- `IsNewerVersion(v1, v2 string) (bool, error)` - Check if v1 > v2

### SemVer Methods

- `String() string` - Get version string (e.g., "1.2.3-alpha+build")
- `Compare(other *SemVer) int` - Returns -1, 0, or 1
- `LessThan(other *SemVer) bool` - Check if v < other
- `GreaterThan(other *SemVer) bool` - Check if v > other
- `Equal(other *SemVer) bool` - Check if v == other
- `GreaterThanOrEqual(other *SemVer) bool` - Check if v >= other
- `LessThanOrEqual(other *SemVer) bool` - Check if v <= other
- `Major() int`, `Minor() int`, `Patch() int` - Get version components
- `Prerelease() string`, `Build() string` - Get metadata

### HTTP Handlers

- `Handler() http.Handler` - Version info endpoint (JSON)
- `HealthHandler() http.Handler` - Health check endpoint
- `HandlerFunc() http.HandlerFunc` - Version info as HandlerFunc
- `HealthHandlerFunc() http.HandlerFunc` - Health check as HandlerFunc
- `Middleware(next http.Handler) http.Handler` - Add version headers to responses

### Info Methods

- `GetSchemas() map[string]string` - Get all schemas (defensive copy)
- `GetAPIs() map[string]string` - Get all APIs (defensive copy)
- `GetComponents() map[string]string` - Get all components (defensive copy)
- `GetCustom() map[string]interface{}` - Get custom metadata (defensive copy)
- `GetSchemaVersion(name string) (string, bool)` - Get schema version
- `GetAPIVersion(name string) (string, bool)` - Get API version
- `GetComponentVersion(name string) (string, bool)` - Get component version
- `LogFields() []zap.Field` - Get zap log fields
- `LoadedAt() time.Time` - Get time version info was loaded
- `String() string` - Get compact string representation
- `MarshalJSON() ([]byte, error)` - Custom JSON serialization

## Manifest File Format

See [_templates/versions.yaml.tmpl](_templates/versions.yaml.tmpl) for a comprehensive template with detailed comments.

### Minimal Example

```yaml
manifest_version: "1.0"
project:
  name: "my-app"
  version: "1.0.0"
```

### Complete Example

```yaml
manifest_version: "1.0"

project:
  name: "my-app"
  version: "1.2.3"

schemas:
  postgres_main: "45"
  redis_cache: "3"

apis:
  rest_v1: "1.15.0"
  grpc: "1.2.0"

components:
  auth_service: "2.1.0"
  notification_service: "1.5.0"

custom:
  environment: "production"
  region: "us-east-1"
  license: "MIT"
```

## Thread Safety

All functions and methods are safe for concurrent use by multiple goroutines. The `Info` struct is immutable after creation, providing lock-free reads with zero overhead.

### Immutability Guarantees

- **Singleton access**: Lock-free atomic reads using `atomic.Value`
- **Field immutability**: All Info struct fields are immutable after initialization
- **Defensive copies**: Map getters (`GetSchemas()`, `GetAPIs()`, etc.) return defensive copies to prevent external mutation
- **Zero contention**: Multiple goroutines can safely read concurrently

```go
// Safe to call from multiple goroutines
go func() {
    info := version.MustGet()
    log.Println(info.Project.Version)

    // GetSchemas() returns a defensive copy - safe to modify
    schemas := info.GetSchemas()
    schemas["new_key"] = "value" // Does not affect Info
}()

go func() {
    info := version.MustGet()
    log.Println(info.Git.Commit)

    // Each goroutine gets its own copy of maps
    apis := info.GetAPIs()
    // Safe to modify without affecting other goroutines
}()
```

## Security

The library follows security best practices to ensure safe operation in production environments:

### Command Execution Safety

Git commands are protected against injection and PATH attacks:
- **PATH validation**: Only allows git from trusted system locations (`/usr/bin`, `/usr/local/bin`, `/opt/homebrew/bin`)
- **Command injection protection**: All git arguments are constants, preventing command injection
- **Output validation**: Commit hashes validated as hexadecimal before use
- **Timeout protection**: 5-second timeout prevents hanging if git is unresponsive
- **Graceful degradation**: Falls back to defaults if git is unavailable

```go
// Security: Git binary locations are whitelisted
// Only these paths are allowed on Unix systems:
safeLocations := []string{
    "/usr/bin/git",
    "/usr/local/bin/git",
    "/opt/homebrew/bin/git",
}
```

### HTTP Request Protection

HTTP handlers enforce request size limits using `MaxBytesReader`:
- Limits request bodies to 1KB (defense in depth)
- Prevents resource exhaustion attacks
- Applied even to GET requests as an additional safety layer

### Production Safeguards

The `Reset()` function includes production protection:
- **Only allowed in test environment**: Uses `testing.Testing()` to detect `go test`
- Panics if called outside test environment
- Prevents accidental state corruption in production

**Warning**: Never use `Reset()` in production code - it's exclusively for testing.

### Dependency Management

- Uses Go 1.24.6+ to ensure all stdlib CVEs are patched
- Regular dependency updates via `go get -u` and `go mod tidy`
- Zero external runtime dependencies (only dev/test dependencies)

### Best Practices

1. **Always validate Go version**: Use Go 1.24.6 or later
2. **Run with race detector**: Test with `go test -race` to catch concurrency issues
3. **Test with latest security patches**: Keep dependencies updated
4. **Monitor timeouts**: If git operations timeout frequently, check git installation

```bash
# Example: Running tests with race detection
go test -race ./...

# Example: Checking for vulnerabilities
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

## Testing

The library includes `Reset()` function for testing. It automatically detects when running under `go test` and panics if called in production:

```go
func TestMyFunction(t *testing.T) {
    defer version.Reset() // Clean up after test

    err := version.Initialize(
        version.WithManifestPath("testdata/versions.yaml"),
    )
    require.NoError(t, err)

    // Test code...
}
```

⚠️ **Warning**: `Reset()` is ONLY for testing and uses `testing.Testing()` to prevent accidental production use.

```bash
# Run tests normally - Reset() automatically allowed
go test ./...

# Run with race detector (recommended)
go test -race ./...

# Run with coverage
go test -cover ./...
```

## Performance

- **Singleton access**: Lock-free atomic reads (~1ns per call)
- **Concurrent safety**: Zero contention (immutable design)
- **Memory overhead**: Single shared instance
- **Initialization**: One-time cost on first access

Benchmark results:
```
BenchmarkGet-8                  1000000000    0.5 ns/op
BenchmarkGet_Concurrent-8       1000000000    0.6 ns/op
```

## Examples

Complete working examples are available in the [examples](examples/) directory:

- **[examples/simple](examples/simple)** - Zero-config usage with basic features
- **[examples/advanced](examples/advanced)** - Production-ready HTTP server with validation

Run examples:
```bash
cd examples/simple && go run main.go
cd examples/advanced && go run main.go
```

## Best Practices

### 1. Initialize Early

```go
func main() {
    // Initialize version before other setup
    if err := version.Initialize(/* options */); err != nil {
        log.Fatal(err)
    }

    // Rest of application initialization
}
```

### 2. Use Validation for Critical Versions

```go
err := version.Initialize(
    version.WithValidators(
        // Ensure database schema is compatible
        version.NewSchemaValidator("postgres_main", "45"),
    ),
)
```

### 3. Embed Manifest in Production

```go
//go:embed versions.yaml
var versionsYAML []byte

version.Initialize(version.WithEmbedded(versionsYAML))
```

### 4. Add Version Headers to HTTP Responses

```go
handler := version.Middleware(mux)
http.ListenAndServe(":8080", handler)
```

### 5. Include in Structured Logs

```go
logger = logger.With(info.LogFields()...)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/itsatony/go-version.git
cd go-version

# Run tests
go test -race -cover ./...

# Run examples
go run examples/simple/main.go
```

### Running Tests

```bash
# All tests with race detection
go test -race ./...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...

# Specific test
go test -run TestInitialize_Success
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

Built with best practices from the Go community and inspired by version management needs in real-world microservices architectures.

## Support

- 📚 [Documentation](https://pkg.go.dev/github.com/itsatony/go-version)
- 💬 [Issues](https://github.com/itsatony/go-version/issues)
- 📧 Contact: [your-email@example.com](mailto:your-email@example.com)

---

Made with ❤️ by [itsatony](https://github.com/itsatony)
