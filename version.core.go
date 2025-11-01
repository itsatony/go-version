// Package version provides multi-dimensional versioning for Go applications.
//
// It tracks project versions, database schemas, API versions, and component versions,
// combining information from manifest files, git metadata, and build-time injection.
//
// # Basic Usage
//
// Zero-config usage with auto-initialization:
//
//	v := version.MustGet()
//	log.Printf("Starting %s v%s", v.Project.Name, v.Project.Version)
//
// With custom configuration:
//
//	err := version.Initialize(
//	    version.WithManifestPath("./config/versions.yaml"),
//	    version.WithValidators(
//	        version.NewSchemaValidator("postgres_main", "45"),
//	    ),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	v := version.MustGet()
//
// # Thread Safety
//
// All functions and methods in this package are safe for concurrent use by
// multiple goroutines. The Info struct is immutable after creation.
//
// # Version Sources
//
// Version information is loaded from multiple sources with this precedence:
//  1. Embedded manifest (via WithEmbedded option)
//  2. File manifest (versions.yaml or custom path)
//  3. Defaults (if no manifest found)
//
// Git and build information is enriched from:
//   - Build-time ldflags injection
//   - runtime/debug.BuildInfo
//   - Git command execution (fallback)
package version

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

var (
	// instance holds the singleton Info instance using atomic for fast reads
	instance atomic.Value // stores *Info

	// initOnce ensures initialization happens exactly once
	initOnce sync.Once

	// initError stores any error from initialization
	initError error

	// resetMu protects Reset() operations from concurrent access
	resetMu sync.Mutex
)

// Initialize configures and initializes the version singleton.
//
// ⚠️  SINGLETON BEHAVIOR: This function can only be called ONCE per process.
// Subsequent calls will return an error, even with different options.
//
// If Initialize is not called, Get() and MustGet() will auto-initialize with
// default options on first access. Once auto-initialized, calling Initialize()
// will return an error.
//
// Thread-safe for concurrent use by multiple goroutines.
//
// When to use Initialize():
//   - Call Initialize() at application startup if you need custom configuration
//   - Call it BEFORE any Get()/MustGet() calls to ensure your config is used
//   - In tests, use version.Reset() (test-only) to allow re-initialization
//
// When to use New():
//   - If you need multiple version info instances with different configs
//   - When building libraries that should avoid global state
//   - For loading version info of other applications/services
//
// Example:
//
//	func main() {
//	    // Initialize with custom config at startup
//	    err := version.Initialize(
//	        version.WithManifestPath("./versions.yaml"),
//	        version.WithGitInfo(),
//	        version.WithBuildInfo(),
//	        version.WithValidators(
//	            version.NewSchemaValidator("postgres_main", "45"),
//	        ),
//	    )
//	    if err != nil {
//	        log.Fatal("Failed to initialize version:", err)
//	    }
//
//	    // Now use Get() or MustGet() anywhere
//	    info := version.MustGet()
//	    log.Printf("Starting %s v%s", info.Project.Name, info.Project.Version)
//	}
func Initialize(opts ...Option) error {
	didRun := false

	initOnce.Do(func() {
		didRun = true
		var info *Info
		info, initError = loadVersionInfo(opts...)
		if initError == nil {
			instance.Store(info)
		}
	})

	// If Do() didn't run, it means Initialize was already called
	if !didRun {
		return fmt.Errorf("%s\nHint: %s", ErrMsgInitializeMultiple, ErrHintInitializeMultiple)
	}

	return initError
}

// Get returns the singleton version info instance.
//
// If Initialize was not called, Get auto-initializes with default options:
//   - Loads from versions.yaml in current directory (or uses defaults if not found)
//   - Includes git information
//   - Includes build information
//
// Returns ErrNotInitialized only if Initialize was called and failed.
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	info, err := version.Get()
//	if err != nil {
//	    log.Printf("Version unavailable: %v", err)
//	    return
//	}
//
//	log.Printf("Running version %s", info.Project.Version)
func Get() (*Info, error) {
	// Fast path: already initialized
	if v := instance.Load(); v != nil {
		if info, ok := v.(*Info); ok && info != nil {
			return info, nil
		}
	}

	// Slow path: auto-initialize with defaults
	initOnce.Do(func() {
		var info *Info
		info, initError = loadVersionInfo(
			WithGitInfo(),
			WithBuildInfo(),
		)
		if initError == nil {
			instance.Store(info)
		}
	})

	if initError != nil {
		return nil, initError
	}

	v := instance.Load()
	if v == nil {
		// This should never happen, but handle it gracefully
		return nil, ErrNotInitialized
	}

	info, ok := v.(*Info)
	if !ok || info == nil {
		return nil, ErrNotInitialized
	}

	return info, nil
}

// MustGet returns the singleton version info instance or panics if unavailable.
//
// Use this in initialization code where version failure should crash the application.
// For runtime code where graceful degradation is preferred, use Get() instead.
//
// If Initialize was not called, MustGet auto-initializes with default options.
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	func main() {
//	    v := version.MustGet()
//	    logger := zap.NewProduction()
//	    logger = logger.With(v.LogFields()...)
//	    logger.Info("Application started")
//	}
func MustGet() *Info {
	info, err := Get()
	if err != nil {
		panic(fmt.Sprintf(ErrFmtMustGetPanic, err))
	}
	return info
}

// New creates a new Info instance without using the singleton.
//
// This is useful for:
//   - Testing with multiple independent instances
//   - Loading version info for other applications/services
//   - Avoiding global state in libraries
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	// Load version from a specific file
//	info, err := version.New(version.WithManifestPath("./app-versions.yaml"))
//	if err != nil {
//	    return err
//	}
//
//	fmt.Printf("Application version: %s\n", info.Project.Version)
func New(opts ...Option) (*Info, error) {
	return loadVersionInfo(opts...)
}

// IsInitialized checks if the version singleton has been initialized.
//
// This function returns true if Initialize() or Get()/MustGet() have been called
// and successfully loaded version information. It returns false if:
//   - No initialization has occurred yet
//   - Initialize() was called but failed
//
// Unlike Get(), this function does NOT trigger auto-initialization, making it
// useful for:
//   - Checking initialization state before making decisions
//   - Avoiding side effects of auto-initialization
//   - Testing initialization workflows
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	if version.IsInitialized() {
//	    info := version.MustGet()
//	    log.Printf("Running version %s", info.Project.Version)
//	} else {
//	    log.Println("Version information not available")
//	}
func IsInitialized() bool {
	v := instance.Load()
	if v == nil {
		return false
	}
	info, ok := v.(*Info)
	return ok && info != nil
}

// Reset clears the singleton instance. USE WITH EXTREME CAUTION.
//
// ⚠️  SECURITY WARNING: This function is ONLY for testing and will PANIC in production.
//
// Reset() is designed exclusively for test scenarios where you need to reinitialize
// the version singleton between tests. It should NEVER be used in production code.
//
// Production Safety:
//   - Panics if not called from within `go test` (uses testing.Testing())
//   - Uses mutex to serialize reset operations
//   - Only call during test setup/teardown with no concurrent goroutines
//
// Thread Safety Limitations:
//   - sync.Once cannot be atomically reset, creating race conditions
//   - Other goroutines may observe partial reset state
//   - Only safe when called with exclusive access (no concurrent version usage)
//
// Best Practices:
//   - Call Reset() at the start or end of each test that needs it
//   - Never call from production code (it will panic)
//   - Ensure no other goroutines are accessing version info during reset
//   - Use defer to ensure cleanup even if test fails
//
// Example (testing only):
//
//	func TestMyFunction(t *testing.T) {
//	    version.Reset() // Clean slate at start
//	    defer version.Reset() // Clean up after test
//
//	    err := version.Initialize(version.WithManifestPath("./testdata/versions.yaml"))
//	    require.NoError(t, err)
//	    // ... test code ...
//	}
func Reset() {
	// SECURITY: Prevent accidental production use
	// Only allow Reset() when running under `go test`
	if !testing.Testing() {
		panic("version.Reset() is only allowed in test environment (must be called from `go test`). " +
			"This function is unsafe for production use due to race conditions. " +
			"See https://pkg.go.dev/testing#Testing for details.")
	}

	// Use mutex to serialize Reset operations
	resetMu.Lock()
	defer resetMu.Unlock()

	// Clear the singleton instance
	var nilInfo *Info
	instance.Store(nilInfo)

	// Reset initialization state
	// Note: sync.Once cannot be truly reset atomically, but this is sufficient
	// for single-threaded test scenarios where Reset is intended to be used
	initOnce = sync.Once{}
	initError = nil
}
