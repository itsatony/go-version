package version

// Manifest constants
const (
	// ManifestVersion is the current version of the manifest format
	ManifestVersion = "1.0"

	// ManifestFilenameYAML is the default YAML manifest filename
	ManifestFilenameYAML = "versions.yaml"

	// ManifestFilenameJSON is the default JSON manifest filename (future support)
	ManifestFilenameJSON = "versions.json"
)

// Default values for version information
const (
	// DefaultGitCommit is used when git info is unavailable
	DefaultGitCommit = "dev"

	// DefaultGitTreeState is used when git status cannot be determined
	DefaultGitTreeState = "clean"

	// DefaultBuildTime is used when build time is not injected
	DefaultBuildTime = "unknown"

	// DefaultBuildUser is used when build user is not injected
	DefaultBuildUser = ""

	// DefaultGoVersion is used as fallback
	DefaultGoVersion = "unknown"

	// DefaultProjectName is used when project name is not specified
	DefaultProjectName = "unknown"

	// DefaultProjectVersion is used when project version is not specified
	DefaultProjectVersion = "0.0.0-dev"
)

// HTTP constants
const (
	// HTTPPathVersion is the default path for version endpoint
	HTTPPathVersion = "/version"

	// HTTPPathHealth is the default path for health endpoint
	HTTPPathHealth = "/health"

	// HTTPContentTypeJSON is the content type for JSON responses
	HTTPContentTypeJSON = "application/json"

	// HTTPCacheControl is the cache control header value for version endpoint
	HTTPCacheControl = "public, max-age=300"

	// HTTPHeaderAppVersion is the header name for application version
	HTTPHeaderAppVersion = "X-App-Version"

	// HTTPHeaderGitCommit is the header name for git commit
	HTTPHeaderGitCommit = "X-Git-Commit"

	// HTTPStatusOK is the status string for successful health checks
	HTTPStatusOK = "ok"

	// HTTPStatusError is the status string for failed health checks
	HTTPStatusError = "error"

	// HTTPErrorMethodNotAllowed is the error message for invalid HTTP methods
	HTTPErrorMethodNotAllowed = "Method not allowed"

	// HTTPErrorVersionUnavailable is the error message when version info cannot be retrieved
	HTTPErrorVersionUnavailable = "Failed to get version info"

	// HTTPHealthErrorMessage is the error message in health check responses
	HTTPHealthErrorMessage = "version not available"
)

// Git tree states
const (
	// GitTreeStateClean indicates no uncommitted changes
	GitTreeStateClean = "clean"

	// GitTreeStateDirty indicates uncommitted changes present
	GitTreeStateDirty = "dirty"
)

// Git command components
const (
	// GitCmdName is the git command name
	GitCmdName = "git"

	// GitCmdRevParse is the git rev-parse subcommand
	GitCmdRevParse = "rev-parse"

	// GitCmdDescribe is the git describe subcommand
	GitCmdDescribe = "describe"

	// GitCmdStatus is the git status subcommand
	GitCmdStatus = "status"

	// GitArgHead is the HEAD argument
	GitArgHead = "HEAD"

	// GitArgTags is the --tags argument
	GitArgTags = "--tags"

	// GitArgExactMatch is the --exact-match argument
	GitArgExactMatch = "--exact-match"

	// GitArgPorcelain is the --porcelain argument
	GitArgPorcelain = "--porcelain"
)

// VCS build info keys (from runtime/debug.BuildInfo)
const (
	// VCSKeyRevision is the key for VCS revision in build info
	VCSKeyRevision = "vcs.revision"

	// VCSKeyTime is the key for VCS commit time in build info
	VCSKeyTime = "vcs.time"

	// VCSKeyModified is the key for VCS modified status in build info
	VCSKeyModified = "vcs.modified"

	// VCSValueTrue is the string value "true" for VCS modified flag
	VCSValueTrue = "true"
)

// Error messages
const (
	// ErrMsgManifestNotFound is returned when manifest file cannot be found
	ErrMsgManifestNotFound = "version manifest not found"

	// ErrMsgManifestParse is returned when manifest cannot be parsed
	ErrMsgManifestParse = "failed to parse version manifest"

	// ErrMsgInvalidManifest is returned when manifest format is invalid
	ErrMsgInvalidManifest = "invalid manifest format"

	// ErrMsgNotInitialized is returned when Get() is called but initialization failed
	ErrMsgNotInitialized = "version not initialized"

	// ErrMsgValidationFailed is returned when version validation fails
	ErrMsgValidationFailed = "version validation failed"

	// ErrMsgInvalidVersion is returned when version string is not valid semver
	ErrMsgInvalidVersion = "invalid version format"

	// ErrMsgInitializeMultiple is returned when Initialize is called more than once
	ErrMsgInitializeMultiple = "version.Initialize() was already called (singleton enforces single initialization)"

	// ErrMsgLoadManifest is returned when manifest loading fails
	ErrMsgLoadManifest = "failed to load manifest"

	// ErrMsgParseYAML is returned when YAML parsing fails
	ErrMsgParseYAML = "failed to parse YAML"

	// ErrMsgProjectNameRequired is returned when project name is missing from manifest
	ErrMsgProjectNameRequired = "project name is required in manifest"

	// ErrMsgProjectVersionRequired is returned when project version is missing from manifest
	ErrMsgProjectVersionRequired = "project version is required in manifest"

	// ErrMsgValidationFailedWrap is returned when validation fails during loading
	ErrMsgValidationFailedWrap = "validation failed"

	// ErrMsgStrictModeManifestRequired is returned in strict mode when manifest is missing
	ErrMsgStrictModeManifestRequired = "strict mode: manifest file is required but not found"
)

// Error hints - actionable suggestions for fixing errors
const (
	// ErrHintManifestNotFound provides guidance when manifest is missing
	ErrHintManifestNotFound = "Create a versions.yaml file or use WithEmbedded() option. Example:\n" +
		"  err := version.Initialize(version.WithManifestPath(\"./versions.yaml\"))"

	// ErrHintStrictMode provides guidance for strict mode errors
	ErrHintStrictMode = "Either create the required manifest file or remove WithStrictMode() option"

	// ErrHintProjectNameRequired provides example for missing project name
	ErrHintProjectNameRequired = "Add project name to your manifest:\n" +
		"  project:\n" +
		"    name: \"your-app-name\"\n" +
		"    version: \"1.0.0\""

	// ErrHintProjectVersionRequired provides example for missing project version
	ErrHintProjectVersionRequired = "Add project version to your manifest:\n" +
		"  project:\n" +
		"    name: \"your-app-name\"\n" +
		"    version: \"1.0.0\""

	// ErrHintParseYAML provides guidance for YAML parsing errors
	ErrHintParseYAML = "Check YAML syntax at https://www.yamllint.com/ or validate with: yamllint versions.yaml"

	// ErrHintInitializeMultiple provides guidance for multiple initialization
	ErrHintInitializeMultiple = "The singleton can only be initialized once. Options:\n" +
		"  1. Call Initialize() at application startup before any Get()/MustGet() calls\n" +
		"  2. Use version.Get() to retrieve the already-initialized singleton\n" +
		"  3. Use version.New(opts...) to create a new non-singleton instance\n" +
		"  4. In tests only: call version.Reset() before re-initializing"

	// ErrHintSchemaNotFound provides guidance when schema is missing
	ErrHintSchemaNotFound = "Add the schema to your versions.yaml:\n" +
		"  schemas:\n" +
		"    schema_name: \"version\""

	// ErrHintAPINotFound provides guidance when API is missing
	ErrHintAPINotFound = "Add the API to your versions.yaml:\n" +
		"  apis:\n" +
		"    api_name: \"version\""

	// ErrHintComponentNotFound provides guidance when component is missing
	ErrHintComponentNotFound = "Add the component to your versions.yaml:\n" +
		"  components:\n" +
		"    component_name: \"version\""

	// ErrHintVersionTooOld provides guidance when version doesn't meet requirements
	ErrHintVersionTooOld = "Update the version in your manifest to meet the minimum requirement"
)

// Validation error message formats
const (
	// ErrFmtSchemaNotFound is the format string for schema not found errors
	ErrFmtSchemaNotFound = "schema '%s' not found in manifest"

	// ErrFmtInvalidSchemaVersion is the format string for invalid schema version errors
	ErrFmtInvalidSchemaVersion = "invalid schema version '%s' for '%s': %w"

	// ErrFmtInvalidMinVersion is the format string for invalid minimum version errors
	ErrFmtInvalidMinVersion = "invalid minimum version '%s' for validator: %w"

	// ErrFmtSchemaTooOld is the format string for schema version too old errors
	ErrFmtSchemaTooOld = "schema '%s' version %s is less than required minimum %s"

	// ErrFmtAPINotFound is the format string for API not found errors
	ErrFmtAPINotFound = "API '%s' not found in manifest"

	// ErrFmtInvalidAPIVersion is the format string for invalid API version errors
	ErrFmtInvalidAPIVersion = "invalid API version '%s' for '%s': %w"

	// ErrFmtAPITooOld is the format string for API version too old errors
	ErrFmtAPITooOld = "API '%s' version %s is less than required minimum %s"

	// ErrFmtComponentNotFound is the format string for component not found errors
	ErrFmtComponentNotFound = "component '%s' not found in manifest"

	// ErrFmtInvalidComponentVersion is the format string for invalid component version errors
	ErrFmtInvalidComponentVersion = "invalid component version '%s' for '%s': %w"

	// ErrFmtComponentTooOld is the format string for component version too old errors
	ErrFmtComponentTooOld = "component '%s' version %s is less than required minimum %s"
)

// Error wrapping format strings
const (
	// ErrFmtCategoryWrap is the format string for wrapping errors with category
	ErrFmtCategoryWrap = "[%s] %s: %w"

	// ErrFmtCategory is the format string for creating errors with category
	ErrFmtCategory = "[%s] %s"

	// ErrFmtMustGetPanic is the format string for MustGet panic messages
	ErrFmtMustGetPanic = "version.MustGet: %v"
)

// Logging field names for structured logging (zap, logrus, etc.)
const (
	// LogFieldProjectName is the field name for project name
	LogFieldProjectName = "project_name"

	// LogFieldProjectVersion is the field name for project version
	LogFieldProjectVersion = "project_version"

	// LogFieldGitCommit is the field name for git commit
	LogFieldGitCommit = "git_commit"

	// LogFieldGitTag is the field name for git tag
	LogFieldGitTag = "git_tag"

	// LogFieldGitTreeState is the field name for git tree state
	LogFieldGitTreeState = "git_tree_state"

	// LogFieldBuildTime is the field name for build time
	LogFieldBuildTime = "build_time"

	// LogFieldBuildUser is the field name for build user
	LogFieldBuildUser = "build_user"

	// LogFieldGoVersion is the field name for Go version
	LogFieldGoVersion = "go_version"
)

// String formatting
const (
	// StringFormatSeparator is the separator used in String() method
	StringFormatSeparator = " "

	// StringFormatGitPrefix is the prefix for git info in String() method
	StringFormatGitPrefix = " ("

	// StringFormatGitSuffix is the suffix for git info in String() method
	StringFormatGitSuffix = ")"
)
