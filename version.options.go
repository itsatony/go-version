package version

import "context"

// Validator defines the interface for version validation.
// Custom validators can be implemented to enforce version constraints.
//
// Thread-safe for concurrent use by multiple goroutines.
type Validator interface {
	Validate(ctx context.Context, info *Info) error
}

// LoadOptions contains configuration for loading version information.
// Use the With* functions to configure options.
type LoadOptions struct {
	// manifestPath is the path to the versions.yaml file
	manifestPath string

	// manifestEmbed contains embedded manifest data
	manifestEmbed []byte

	// includeGit enables git information enrichment
	includeGit bool

	// includeBuild enables build information enrichment
	includeBuild bool

	// validators are run after loading to validate the version info
	validators []Validator

	// strictMode enables strict validation and error handling
	strictMode bool

	// ctx is the context for initialization and validation
	// If nil, context.Background() is used
	ctx context.Context
}

// defaultLoadOptions returns the default load options.
// By default, loads from versions.yaml and includes git and build info.
func defaultLoadOptions() *LoadOptions {
	return &LoadOptions{
		manifestPath: ManifestFilenameYAML,
		includeGit:   true,
		includeBuild: true,
	}
}

// Option is a functional option for configuring version loading.
type Option func(*LoadOptions)

// WithManifestPath sets the path to the version manifest file.
// Default is "versions.yaml" in the current directory.
//
// Example:
//
//	info, err := version.New(version.WithManifestPath("./config/versions.yaml"))
func WithManifestPath(path string) Option {
	return func(o *LoadOptions) {
		o.manifestPath = path
	}
}

// WithEmbedded sets embedded manifest data.
// This takes precedence over file-based loading.
//
// Example:
//
//	//go:embed versions.yaml
//	var versionsYAML []byte
//
//	info, err := version.New(version.WithEmbedded(versionsYAML))
func WithEmbedded(data []byte) Option {
	return func(o *LoadOptions) {
		o.manifestEmbed = data
	}
}

// WithGitInfo enables git information enrichment.
// Git info includes commit hash, tag, tree state, and commit time.
// This is enabled by default.
//
// Example:
//
//	info, err := version.New(version.WithGitInfo())
func WithGitInfo() Option {
	return func(o *LoadOptions) {
		o.includeGit = true
	}
}

// WithoutGitInfo disables git information enrichment.
//
// Example:
//
//	info, err := version.New(version.WithoutGitInfo())
func WithoutGitInfo() Option {
	return func(o *LoadOptions) {
		o.includeGit = false
	}
}

// WithBuildInfo enables build information enrichment.
// Build info includes build time, user, and Go version.
// This is enabled by default.
//
// Example:
//
//	info, err := version.New(version.WithBuildInfo())
func WithBuildInfo() Option {
	return func(o *LoadOptions) {
		o.includeBuild = true
	}
}

// WithoutBuildInfo disables build information enrichment.
//
// Example:
//
//	info, err := version.New(version.WithoutBuildInfo())
func WithoutBuildInfo() Option {
	return func(o *LoadOptions) {
		o.includeBuild = false
	}
}

// WithValidators adds custom validators to run after loading.
// Validators are run in order and loading fails if any validator returns an error.
//
// Example:
//
//	info, err := version.New(
//	    version.WithValidators(
//	        version.NewSchemaValidator("postgres_main", "45"),
//	    ),
//	)
func WithValidators(validators ...Validator) Option {
	return func(o *LoadOptions) {
		o.validators = append(o.validators, validators...)
	}
}

// WithStrictMode enables strict validation and error handling.
//
// In strict mode:
//   - Missing manifest files are fatal (no fallback to defaults)
//   - Invalid or unparseable manifests cause immediate failure
//   - All validation errors are treated as fatal
//
// Strict mode is useful for production environments where you want to ensure
// all version information is explicitly defined and valid.
//
// Example:
//
//	err := version.Initialize(
//	    version.WithManifestPath("./versions.yaml"),
//	    version.WithStrictMode(),
//	    version.WithValidators(
//	        version.NewSchemaValidator("postgres_main", "45"),
//	    ),
//	)
//	if err != nil {
//	    log.Fatal("Version initialization failed:", err)
//	}
func WithStrictMode() Option {
	return func(o *LoadOptions) {
		o.strictMode = true
	}
}

// WithContext sets the context for initialization and validation.
//
// The provided context is used for:
//   - Validator execution (allows cancellation, timeouts, tracing)
//   - Future async operations (if added)
//
// If not provided, context.Background() is used.
//
// Example with timeout:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	err := version.Initialize(
//	    version.WithContext(ctx),
//	    version.WithManifestPath("./versions.yaml"),
//	    version.WithValidators(
//	        version.NewSchemaValidator("postgres_main", "45"),
//	    ),
//	)
//
// Example with tracing:
//
//	ctx := trace.ContextWithSpan(context.Background(), span)
//	info, err := version.New(version.WithContext(ctx))
func WithContext(ctx context.Context) Option {
	return func(o *LoadOptions) {
		o.ctx = ctx
	}
}
