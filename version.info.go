package version

import (
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

// Manifest represents the structure of a versions.yaml file.
// This is the file format that users create to define their version information.
type Manifest struct {
	// ManifestVersion is the version of the manifest format itself
	ManifestVersion string `yaml:"manifest_version" json:"manifest_version"`

	// Project contains the main project version information
	Project ProjectManifest `yaml:"project" json:"project"`

	// Schemas contains database schema versions (e.g., "postgres_main": "47")
	Schemas map[string]string `yaml:"Schemas,omitempty" json:"schemas,omitempty"`

	// APIs contains API version numbers (e.g., "rest_v1": "1.15.0")
	APIs map[string]string `yaml:"APIs,omitempty" json:"apis,omitempty"`

	// Components contains dependency/component versions (e.g., "aigentchat": "3.4.1")
	Components map[string]string `yaml:"Components,omitempty" json:"components,omitempty"`

	// Custom contains any custom version dimensions defined by the user
	Custom map[string]interface{} `yaml:"Custom,omitempty" json:"custom,omitempty"`
}

// ProjectManifest represents the project section of the manifest
type ProjectManifest struct {
	// Name is the project/application name
	Name string `yaml:"name" json:"name"`

	// Version is the semantic version of the project
	Version string `yaml:"version" json:"version"`
}

// Info contains the complete runtime version information for an application.
// It combines information from the manifest file, git metadata, and build-time injection.
//
// Info is IMMUTABLE after creation, making it safe for concurrent access by multiple
// goroutines without requiring locks. All maps are unexported and accessed via defensive
// copy getters to ensure true immutability.
//
// Thread-safe for concurrent use by multiple goroutines.
type Info struct {
	// Project contains the project name and version
	Project ProjectVersion `json:"project"`

	// Git contains git metadata (commit, tag, tree state)
	Git GitInfo `json:"git"`

	// Build contains build-time information
	Build BuildInfo `json:"build"`

	// schemas contains database schema versions (unexported for immutability)
	schemas map[string]string

	// apis contains API version numbers (unexported for immutability)
	apis map[string]string

	// components contains dependency/component versions (unexported for immutability)
	components map[string]string

	// custom contains any custom version dimensions (unexported for immutability)
	custom map[string]interface{}

	// loadedAt is the time this Info was created (internal use)
	loadedAt time.Time
}

// ProjectVersion represents the project name and version
type ProjectVersion struct {
	// Name is the project/application name
	Name string `json:"name"`

	// Version is the semantic version of the project
	Version string `json:"version"`
}

// GitInfo contains git metadata injected at build time or extracted from runtime
type GitInfo struct {
	// Commit is the full git commit hash
	Commit string `json:"commit"`

	// Tag is the git tag (if any) for this commit
	Tag string `json:"tag,omitempty"`

	// TreeState indicates whether there were uncommitted changes ("clean" or "dirty")
	TreeState string `json:"tree_state"`

	// CommitTime is the timestamp of the commit
	CommitTime string `json:"commit_time,omitempty"`
}

// BuildInfo contains information about when and how the binary was built
type BuildInfo struct {
	// Time is when the binary was built (ISO 8601 format)
	Time string `json:"time"`

	// User is who built the binary (username or CI system)
	User string `json:"user,omitempty"`

	// GoVersion is the version of Go used to build the binary
	GoVersion string `json:"go_version"`
}

// GetSchemaVersion returns the version for a named schema.
// Returns the version and true if found, empty string and false if not found.
//
// Thread-safe for concurrent use by multiple goroutines.
func (i *Info) GetSchemaVersion(name string) (string, bool) {
	if i.schemas == nil {
		return "", false
	}
	v, ok := i.schemas[name]
	return v, ok
}

// GetAPIVersion returns the version for a named API.
// Returns the version and true if found, empty string and false if not found.
//
// Thread-safe for concurrent use by multiple goroutines.
func (i *Info) GetAPIVersion(name string) (string, bool) {
	if i.apis == nil {
		return "", false
	}
	v, ok := i.apis[name]
	return v, ok
}

// GetComponentVersion returns the version for a named component.
// Returns the version and true if found, empty string and false if not found.
//
// Thread-safe for concurrent use by multiple goroutines.
func (i *Info) GetComponentVersion(name string) (string, bool) {
	if i.components == nil {
		return "", false
	}
	v, ok := i.components[name]
	return v, ok
}

// GetSchemas returns a defensive copy of all database schema versions.
// Modifying the returned map does not affect the Info instance.
//
// Thread-safe for concurrent use by multiple goroutines.
func (i *Info) GetSchemas() map[string]string {
	if i.schemas == nil {
		return nil
	}
	copy := make(map[string]string, len(i.schemas))
	for k, v := range i.schemas {
		copy[k] = v
	}
	return copy
}

// GetAPIs returns a defensive copy of all API versions.
// Modifying the returned map does not affect the Info instance.
//
// Thread-safe for concurrent use by multiple goroutines.
func (i *Info) GetAPIs() map[string]string {
	if i.apis == nil {
		return nil
	}
	copy := make(map[string]string, len(i.apis))
	for k, v := range i.apis {
		copy[k] = v
	}
	return copy
}

// GetComponents returns a defensive copy of all component versions.
// Modifying the returned map does not affect the Info instance.
//
// Thread-safe for concurrent use by multiple goroutines.
func (i *Info) GetComponents() map[string]string {
	if i.components == nil {
		return nil
	}
	copy := make(map[string]string, len(i.components))
	for k, v := range i.components {
		copy[k] = v
	}
	return copy
}

// GetCustom returns a defensive copy of all custom version dimensions.
// Modifying the returned map does not affect the Info instance.
//
// Thread-safe for concurrent use by multiple goroutines.
func (i *Info) GetCustom() map[string]interface{} {
	if i.custom == nil {
		return nil
	}
	copy := make(map[string]interface{}, len(i.custom))
	for k, v := range i.custom {
		copy[k] = v
	}
	return copy
}

// LoadedAt returns the time when this version info was loaded.
// Useful for diagnostics and cache invalidation.
//
// Thread-safe for concurrent use by multiple goroutines.
func (i *Info) LoadedAt() time.Time {
	return i.loadedAt
}

// LogFields returns structured logging fields for use with zap logger.
// This provides a convenient way to include version info in log entries.
//
// Example:
//
//	logger := zap.NewProduction()
//	logger = logger.With(info.LogFields()...)
//	logger.Info("Application started")
//
// Thread-safe for concurrent use by multiple goroutines.
func (i *Info) LogFields() []zap.Field {
	return []zap.Field{
		zap.String(LogFieldProjectName, i.Project.Name),
		zap.String(LogFieldProjectVersion, i.Project.Version),
		zap.String(LogFieldGitCommit, i.Git.Commit),
		zap.String(LogFieldGitTag, i.Git.Tag),
		zap.String(LogFieldGitTreeState, i.Git.TreeState),
		zap.String(LogFieldBuildTime, i.Build.Time),
		zap.String(LogFieldBuildUser, i.Build.User),
		zap.String(LogFieldGoVersion, i.Build.GoVersion),
	}
}

// String returns a human-readable string representation of the version info.
//
// Thread-safe for concurrent use by multiple goroutines.
func (i *Info) String() string {
	return i.Project.Name + " " + i.Project.Version + " (" + i.Git.Commit + ")"
}

// MarshalJSON implements json.Marshaler to ensure consistent JSON output.
// The loadedAt field is excluded from JSON serialization.
// Unexported map fields are serialized with their JSON names.
func (i *Info) MarshalJSON() ([]byte, error) {
	// Manually construct JSON structure to include unexported fields
	type jsonInfo struct {
		Project    ProjectVersion         `json:"project"`
		Git        GitInfo                `json:"git"`
		Build      BuildInfo              `json:"build"`
		Schemas    map[string]string      `json:"schemas,omitempty"`
		APIs       map[string]string      `json:"apis,omitempty"`
		Components map[string]string      `json:"components,omitempty"`
		Custom     map[string]interface{} `json:"custom,omitempty"`
	}

	return json.Marshal(jsonInfo{
		Project:    i.Project,
		Git:        i.Git,
		Build:      i.Build,
		Schemas:    i.schemas,
		APIs:       i.apis,
		Components: i.components,
		Custom:     i.custom,
	})
}

// UnmarshalJSON implements json.Unmarshaler to populate unexported fields.
func (i *Info) UnmarshalJSON(data []byte) error {
	// Use a temporary struct with exported fields for unmarshaling
	type jsonInfo struct {
		Project    ProjectVersion         `json:"project"`
		Git        GitInfo                `json:"git"`
		Build      BuildInfo              `json:"build"`
		Schemas    map[string]string      `json:"schemas,omitempty"`
		APIs       map[string]string      `json:"apis,omitempty"`
		Components map[string]string      `json:"components,omitempty"`
		Custom     map[string]interface{} `json:"custom,omitempty"`
	}

	var temp jsonInfo
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Populate Info fields
	i.Project = temp.Project
	i.Git = temp.Git
	i.Build = temp.Build
	i.schemas = temp.Schemas
	i.apis = temp.APIs
	i.components = temp.Components
	i.custom = temp.Custom

	return nil
}
