package semver

import (
	"fmt"
	"strconv"
	"strings"
)

// Version represents a semantic version
type Version struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
	Build      string
}

// Parse parses a semantic version string.
// Supports formats: "1.2.3", "v1.2.3", "1.2.3-alpha", "1.2.3+build"
func Parse(s string) (*Version, error) {
	// Remove leading 'v' if present
	s = strings.TrimPrefix(s, "v")

	if s == "" {
		return nil, fmt.Errorf("empty version string")
	}

	// Split on '+' for build metadata
	parts := strings.SplitN(s, "+", 2)
	versionPart := parts[0]
	buildPart := ""
	if len(parts) > 1 {
		buildPart = parts[1]
	}

	// Split on '-' for prerelease
	parts = strings.SplitN(versionPart, "-", 2)
	corePart := parts[0]
	prereleasePart := ""
	if len(parts) > 1 {
		prereleasePart = parts[1]
	}

	// Parse core version (major.minor.patch)
	coreParts := strings.Split(corePart, ".")
	if len(coreParts) < 1 || len(coreParts) > 3 {
		return nil, fmt.Errorf("invalid version format: %s", s)
	}

	v := &Version{
		Prerelease: prereleasePart,
		Build:      buildPart,
	}

	var err error

	// Parse major
	v.Major, err = strconv.Atoi(coreParts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", coreParts[0])
	}

	// Parse minor (default 0)
	if len(coreParts) > 1 {
		v.Minor, err = strconv.Atoi(coreParts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid minor version: %s", coreParts[1])
		}
	}

	// Parse patch (default 0)
	if len(coreParts) > 2 {
		v.Patch, err = strconv.Atoi(coreParts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid patch version: %s", coreParts[2])
		}
	}

	return v, nil
}

// String returns the string representation of the version
func (v *Version) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)

	if v.Prerelease != "" {
		s += "-" + v.Prerelease
	}

	if v.Build != "" {
		s += "+" + v.Build
	}

	return s
}

// Compare compares two versions.
// Returns -1 if v < other, 0 if v == other, 1 if v > other.
// Prerelease versions have lower precedence than normal versions.
func (v *Version) Compare(other *Version) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}

	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}

	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}

	// Handle prerelease comparison
	// No prerelease > prerelease
	if v.Prerelease == "" && other.Prerelease != "" {
		return 1
	}
	if v.Prerelease != "" && other.Prerelease == "" {
		return -1
	}

	// Both have prerelease or both don't
	if v.Prerelease != other.Prerelease {
		if v.Prerelease < other.Prerelease {
			return -1
		}
		return 1
	}

	return 0
}

// LessThan returns true if v < other
func (v *Version) LessThan(other *Version) bool {
	return v.Compare(other) < 0
}

// GreaterThan returns true if v > other
func (v *Version) GreaterThan(other *Version) bool {
	return v.Compare(other) > 0
}

// Equal returns true if v == other
func (v *Version) Equal(other *Version) bool {
	return v.Compare(other) == 0
}

// GreaterThanOrEqual returns true if v >= other
func (v *Version) GreaterThanOrEqual(other *Version) bool {
	return v.Compare(other) >= 0
}

// LessThanOrEqual returns true if v <= other
func (v *Version) LessThanOrEqual(other *Version) bool {
	return v.Compare(other) <= 0
}
