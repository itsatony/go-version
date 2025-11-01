package version

import (
	"context"
	"fmt"

	"github.com/itsatony/go-version/internal/semver"
)

// genericVersionValidator handles version validation for any dimension (schema, API, or component).
// It uses a getter function to retrieve the version and error format strings for messages.
type genericVersionValidator struct {
	dimensionType  string // "schema", "API", or "component"
	name           string // specific item name (e.g., "postgres_main")
	minVersion     string // minimum required version
	getterFunc     func(*Info, string) (string, bool)
	errNotFoundFmt string
	errInvalidFmt  string
	errTooOldFmt   string
}

// Validate checks if the version meets the minimum requirement.
func (v *genericVersionValidator) Validate(ctx context.Context, info *Info) error {
	actual, ok := v.getterFunc(info, v.name)
	if !ok {
		// Add appropriate hint based on dimension type
		var hint string
		switch v.dimensionType {
		case "schema":
			hint = ErrHintSchemaNotFound
		case "API":
			hint = ErrHintAPINotFound
		case "component":
			hint = ErrHintComponentNotFound
		}
		return fmt.Errorf(v.errNotFoundFmt+"\nHint: %s", v.name, hint)
	}

	actualVer, err := semver.Parse(actual)
	if err != nil {
		return fmt.Errorf(v.errInvalidFmt, actual, v.name, err)
	}

	minVer, err := semver.Parse(v.minVersion)
	if err != nil {
		return fmt.Errorf(ErrFmtInvalidMinVersion, v.minVersion, err)
	}

	if actualVer.LessThan(minVer) {
		return fmt.Errorf(v.errTooOldFmt+"\nHint: %s", v.name, actual, v.minVersion, ErrHintVersionTooOld)
	}

	return nil
}

// SchemaValidator validates that a database schema meets a minimum version requirement.
// It compares the actual schema version from the Info against the required minimum version.
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	err := version.Initialize(
//	    version.WithValidators(
//	        version.NewSchemaValidator("postgres_main", "45"),
//	        version.NewSchemaValidator("redis_cache", "3"),
//	    ),
//	)
//	if err != nil {
//	    log.Fatal("Schema version too old:", err)
//	}
type SchemaValidator = genericVersionValidator

// NewSchemaValidator creates a validator that enforces a minimum schema version.
// The schemaName must match a key in the Info.Schemas map.
// The minVersion should be a semantic version string (e.g., "1.2.3" or "45").
//
// Returns an error during validation if:
//   - The schema is not found in the manifest
//   - The actual version is less than the minimum version
//   - Version parsing fails
func NewSchemaValidator(schemaName, minVersion string) *SchemaValidator {
	return &genericVersionValidator{
		dimensionType:  "schema",
		name:           schemaName,
		minVersion:     minVersion,
		getterFunc:     (*Info).GetSchemaVersion,
		errNotFoundFmt: ErrFmtSchemaNotFound,
		errInvalidFmt:  ErrFmtInvalidSchemaVersion,
		errTooOldFmt:   ErrFmtSchemaTooOld,
	}
}

// APIValidator validates that an API meets a minimum version requirement.
// It compares the actual API version from the Info against the required minimum version.
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	err := version.Initialize(
//	    version.WithValidators(
//	        version.NewAPIValidator("rest_v1", "1.15.0"),
//	        version.NewAPIValidator("grpc", "1.2.0"),
//	    ),
//	)
type APIValidator = genericVersionValidator

// NewAPIValidator creates a validator that enforces a minimum API version.
// The apiName must match a key in the Info.APIs map.
// The minVersion should be a semantic version string (e.g., "1.2.3").
func NewAPIValidator(apiName, minVersion string) *APIValidator {
	return &genericVersionValidator{
		dimensionType:  "API",
		name:           apiName,
		minVersion:     minVersion,
		getterFunc:     (*Info).GetAPIVersion,
		errNotFoundFmt: ErrFmtAPINotFound,
		errInvalidFmt:  ErrFmtInvalidAPIVersion,
		errTooOldFmt:   ErrFmtAPITooOld,
	}
}

// ComponentValidator validates that a component meets a minimum version requirement.
// It compares the actual component version from the Info against the required minimum version.
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	err := version.Initialize(
//	    version.WithValidators(
//	        version.NewComponentValidator("aigentchat", "3.4.0"),
//	    ),
//	)
type ComponentValidator = genericVersionValidator

// NewComponentValidator creates a validator that enforces a minimum component version.
// The componentName must match a key in the Info.Components map.
// The minVersion should be a semantic version string (e.g., "1.2.3").
func NewComponentValidator(componentName, minVersion string) *ComponentValidator {
	return &genericVersionValidator{
		dimensionType:  "component",
		name:           componentName,
		minVersion:     minVersion,
		getterFunc:     (*Info).GetComponentVersion,
		errNotFoundFmt: ErrFmtComponentNotFound,
		errInvalidFmt:  ErrFmtInvalidComponentVersion,
		errTooOldFmt:   ErrFmtComponentTooOld,
	}
}

// ValidatorFunc is a function adapter that allows using functions as Validators.
// This enables inline validator creation without defining new types.
//
// Example:
//
//	customValidator := version.ValidatorFunc(func(ctx context.Context, info *Info) error {
//	    if info.Project.Version == "0.0.0-dev" {
//	        return fmt.Errorf("development version not allowed in production")
//	    }
//	    return nil
//	})
//
//	err := version.Initialize(
//	    version.WithValidators(customValidator),
//	)
type ValidatorFunc func(ctx context.Context, info *Info) error

// Validate calls the function.
func (f ValidatorFunc) Validate(ctx context.Context, info *Info) error {
	return f(ctx, info)
}
