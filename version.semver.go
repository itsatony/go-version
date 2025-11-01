package version

import (
	"fmt"

	"github.com/itsatony/go-version/internal/semver"
)

// SemVer represents a semantic version following the semver.org specification.
//
// It can parse and compare versions in formats:
//   - "1.2.3"
//   - "v1.2.3"
//   - "1.2.3-alpha"
//   - "1.2.3+build"
//   - "1.2.3-beta.1+build.123"
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	v1, err := version.ParseSemVer("1.2.3")
//	v2, err := version.ParseSemVer("2.0.0")
//	if v1.LessThan(v2) {
//	    fmt.Println("v1 is older")
//	}
type SemVer struct {
	internal *semver.Version
}

// ParseSemVer parses a semantic version string.
//
// Supported formats:
//   - Basic: "1.2.3", "v1.2.3"
//   - Prerelease: "1.2.3-alpha", "1.2.3-beta.1"
//   - Build metadata: "1.2.3+build", "1.2.3+20230101"
//   - Combined: "1.2.3-alpha+build"
//
// Returns an error if the version string is not valid semver format.
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	v, err := version.ParseSemVer("1.2.3-alpha+build")
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Version: %s\n", v)
func ParseSemVer(s string) (*SemVer, error) {
	internal, err := semver.Parse(s)
	if err != nil {
		return nil, err
	}
	return &SemVer{internal: internal}, nil
}

// MustParseSemVer parses a semantic version string or panics on error.
//
// Use this in initialization code where invalid versions should crash the application.
// For runtime parsing where graceful error handling is needed, use ParseSemVer instead.
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	v := version.MustParseSemVer("1.2.3")
//	fmt.Printf("Version: %s\n", v)
func MustParseSemVer(s string) *SemVer {
	v, err := ParseSemVer(s)
	if err != nil {
		panic(fmt.Sprintf("version.MustParseSemVer: %v", err))
	}
	return v
}

// String returns the string representation of the semantic version.
//
// Format: "major.minor.patch[-prerelease][+build]"
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	v := version.MustParseSemVer("1.2.3-alpha+build")
//	fmt.Println(v.String()) // "1.2.3-alpha+build"
func (v *SemVer) String() string {
	return v.internal.String()
}

// Major returns the major version number.
//
// Thread-safe for concurrent use by multiple goroutines.
func (v *SemVer) Major() int {
	return v.internal.Major
}

// Minor returns the minor version number.
//
// Thread-safe for concurrent use by multiple goroutines.
func (v *SemVer) Minor() int {
	return v.internal.Minor
}

// Patch returns the patch version number.
//
// Thread-safe for concurrent use by multiple goroutines.
func (v *SemVer) Patch() int {
	return v.internal.Patch
}

// Prerelease returns the prerelease identifier (empty string if none).
//
// Thread-safe for concurrent use by multiple goroutines.
func (v *SemVer) Prerelease() string {
	return v.internal.Prerelease
}

// Build returns the build metadata (empty string if none).
//
// Thread-safe for concurrent use by multiple goroutines.
func (v *SemVer) Build() string {
	return v.internal.Build
}

// Compare compares two semantic versions.
//
// Returns:
//   - -1 if v < other
//   - 0 if v == other
//   - 1 if v > other
//
// Comparison follows semver.org specification:
//   - Major, minor, patch are compared numerically
//   - Prerelease versions have lower precedence than release versions
//   - Build metadata is ignored in comparison
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	v1 := version.MustParseSemVer("1.2.3")
//	v2 := version.MustParseSemVer("2.0.0")
//	result := v1.Compare(v2) // -1
func (v *SemVer) Compare(other *SemVer) int {
	return v.internal.Compare(other.internal)
}

// LessThan returns true if v < other.
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	v1 := version.MustParseSemVer("1.2.3")
//	v2 := version.MustParseSemVer("2.0.0")
//	if v1.LessThan(v2) {
//	    fmt.Println("v1 is older")
//	}
func (v *SemVer) LessThan(other *SemVer) bool {
	return v.internal.LessThan(other.internal)
}

// GreaterThan returns true if v > other.
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	v1 := version.MustParseSemVer("2.0.0")
//	v2 := version.MustParseSemVer("1.2.3")
//	if v1.GreaterThan(v2) {
//	    fmt.Println("v1 is newer")
//	}
func (v *SemVer) GreaterThan(other *SemVer) bool {
	return v.internal.GreaterThan(other.internal)
}

// Equal returns true if v == other.
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	v1 := version.MustParseSemVer("1.2.3")
//	v2 := version.MustParseSemVer("1.2.3")
//	if v1.Equal(v2) {
//	    fmt.Println("versions match")
//	}
func (v *SemVer) Equal(other *SemVer) bool {
	return v.internal.Equal(other.internal)
}

// GreaterThanOrEqual returns true if v >= other.
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	v1 := version.MustParseSemVer("2.0.0")
//	v2 := version.MustParseSemVer("1.2.3")
//	if v1.GreaterThanOrEqual(v2) {
//	    fmt.Println("v1 is newer or same")
//	}
func (v *SemVer) GreaterThanOrEqual(other *SemVer) bool {
	return v.internal.GreaterThanOrEqual(other.internal)
}

// LessThanOrEqual returns true if v <= other.
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	v1 := version.MustParseSemVer("1.2.3")
//	v2 := version.MustParseSemVer("2.0.0")
//	if v1.LessThanOrEqual(v2) {
//	    fmt.Println("v1 is older or same")
//	}
func (v *SemVer) LessThanOrEqual(other *SemVer) bool {
	return v.internal.LessThanOrEqual(other.internal)
}

// CompareVersions compares two version strings.
//
// This is a convenience function that parses both versions and compares them.
//
// Returns:
//   - -1 if v1 < v2
//   - 0 if v1 == v2
//   - 1 if v1 > v2
//   - error if either version is invalid
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	result, err := version.CompareVersions("1.2.3", "2.0.0")
//	if err != nil {
//	    return err
//	}
//	if result < 0 {
//	    fmt.Println("v1 is older")
//	}
func CompareVersions(v1, v2 string) (int, error) {
	ver1, err := ParseSemVer(v1)
	if err != nil {
		return 0, fmt.Errorf("invalid version v1: %w", err)
	}

	ver2, err := ParseSemVer(v2)
	if err != nil {
		return 0, fmt.Errorf("invalid version v2: %w", err)
	}

	return ver1.Compare(ver2), nil
}

// IsNewerVersion checks if version v1 is newer than v2.
//
// This is a convenience function for common version comparison use cases.
//
// Returns:
//   - true if v1 > v2
//   - false if v1 <= v2
//   - error if either version is invalid
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	isNewer, err := version.IsNewerVersion("2.0.0", "1.2.3")
//	if err != nil {
//	    return err
//	}
//	if isNewer {
//	    fmt.Println("Upgrade available")
//	}
func IsNewerVersion(v1, v2 string) (bool, error) {
	ver1, err := ParseSemVer(v1)
	if err != nil {
		return false, fmt.Errorf("invalid version v1: %w", err)
	}

	ver2, err := ParseSemVer(v2)
	if err != nil {
		return false, fmt.Errorf("invalid version v2: %w", err)
	}

	return ver1.GreaterThan(ver2), nil
}
