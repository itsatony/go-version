package version

import (
	"fmt"

	"github.com/itsatony/go-cuserr"
)

// Error categories for classification
const (
	// ErrCategoryManifest indicates manifest loading/parsing errors
	ErrCategoryManifest = "MANIFEST"

	// ErrCategoryValidation indicates validation errors
	ErrCategoryValidation = "VALIDATION"

	// ErrCategoryBuildInfo indicates build info extraction errors
	ErrCategoryBuildInfo = "BUILDINFO"

	// ErrCategoryHTTP indicates HTTP handler errors
	ErrCategoryHTTP = "HTTP"

	// ErrCategoryCore indicates core initialization errors
	ErrCategoryCore = "CORE"
)

// Error codes for specific error types
const (
	// ErrCodeManifestNotFound indicates manifest file not found
	ErrCodeManifestNotFound = "MANIFEST_NOT_FOUND"

	// ErrCodeManifestParse indicates manifest parsing failure
	ErrCodeManifestParse = "MANIFEST_PARSE_FAILED"

	// ErrCodeInvalidManifest indicates invalid manifest format
	ErrCodeInvalidManifest = "INVALID_MANIFEST_FORMAT"

	// ErrCodeValidationFailed indicates version validation failure
	ErrCodeValidationFailed = "VALIDATION_FAILED"

	// ErrCodeNotInitialized indicates singleton not initialized
	ErrCodeNotInitialized = "NOT_INITIALIZED"

	// ErrCodeMultipleInitialize indicates Initialize called multiple times
	ErrCodeMultipleInitialize = "MULTIPLE_INITIALIZE"

	// ErrCodeInvalidVersion indicates invalid version format
	ErrCodeInvalidVersion = "INVALID_VERSION_FORMAT"

	// ErrCodeLoadManifest indicates manifest loading failure
	ErrCodeLoadManifest = "LOAD_MANIFEST_FAILED"

	// ErrCodeParseYAML indicates YAML parsing failure
	ErrCodeParseYAML = "PARSE_YAML_FAILED"

	// ErrCodeProjectNameRequired indicates missing project name
	ErrCodeProjectNameRequired = "PROJECT_NAME_REQUIRED"

	// ErrCodeProjectVersionRequired indicates missing project version
	ErrCodeProjectVersionRequired = "PROJECT_VERSION_REQUIRED"
)

// Sentinel errors for common failure conditions.
// These can be checked with cuserr.IsErrorCode() or by comparing error codes.
var (
	// ErrManifestNotFound is returned when the version manifest file cannot be found
	ErrManifestNotFound = cuserr.NewCustomErrorWithCategory(
		cuserr.ErrorCategory(ErrCategoryManifest),
		ErrCodeManifestNotFound,
		ErrMsgManifestNotFound,
	)

	// ErrManifestParse is returned when the manifest file cannot be parsed
	ErrManifestParse = cuserr.NewCustomErrorWithCategory(
		cuserr.ErrorCategory(ErrCategoryManifest),
		ErrCodeManifestParse,
		ErrMsgManifestParse,
	)

	// ErrInvalidManifest is returned when the manifest format is invalid
	ErrInvalidManifest = cuserr.NewCustomErrorWithCategory(
		cuserr.ErrorCategory(ErrCategoryManifest),
		ErrCodeInvalidManifest,
		ErrMsgInvalidManifest,
	)

	// ErrNotInitialized is returned when Get() is called but initialization failed
	ErrNotInitialized = cuserr.NewCustomErrorWithCategory(
		cuserr.ErrorCategory(ErrCategoryCore),
		ErrCodeNotInitialized,
		ErrMsgNotInitialized,
	)

	// ErrValidationFailed is returned when version validation fails
	ErrValidationFailed = cuserr.NewCustomErrorWithCategory(
		cuserr.ErrorCategory(ErrCategoryValidation),
		ErrCodeValidationFailed,
		ErrMsgValidationFailed,
	)

	// ErrInvalidVersion is returned when a version string is not valid semver
	ErrInvalidVersion = cuserr.NewCustomErrorWithCategory(
		cuserr.ErrorCategory(ErrCategoryValidation),
		ErrCodeInvalidVersion,
		ErrMsgInvalidVersion,
	)
)

// ErrorCategory represents the category of an error for internal classification.
// Deprecated: Use string constants ErrCategory* instead.
type ErrorCategory string

const (
	// CategoryManifest indicates manifest loading/parsing errors
	// Deprecated: Use ErrCategoryManifest constant
	CategoryManifest ErrorCategory = "manifest"

	// CategoryValidation indicates validation errors
	// Deprecated: Use ErrCategoryValidation constant
	CategoryValidation ErrorCategory = "validation"

	// CategoryBuildInfo indicates build info extraction errors
	// Deprecated: Use ErrCategoryBuildInfo constant
	CategoryBuildInfo ErrorCategory = "buildinfo"

	// CategoryHTTP indicates HTTP handler errors
	// Deprecated: Use ErrCategoryHTTP constant
	CategoryHTTP ErrorCategory = "http"
)

// wrapError wraps an error with a category and message for internal use.
// The category is included in the error message for better debugging.
func wrapError(err error, category ErrorCategory, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(ErrFmtCategoryWrap, category, msg, err)
}

// newCategoryErrorWithHint creates a new error with a category prefix and actionable hint
func newCategoryErrorWithHint(category ErrorCategory, msg, hint string) error {
	return fmt.Errorf(ErrFmtCategory+"\nHint: %s", category, msg, hint)
}

// wrapErrorWithHint wraps an error with a category, message, and actionable hint
func wrapErrorWithHint(err error, category ErrorCategory, msg, hint string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(ErrFmtCategoryWrap+"\nHint: %s", category, msg, err, hint)
}
