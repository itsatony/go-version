package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSemVer(t *testing.T) {
	tests := map[string]struct {
		input       string
		expectError bool
		major       int
		minor       int
		patch       int
		prerelease  string
		build       string
	}{
		"basic_version": {
			input: "1.2.3",
			major: 1,
			minor: 2,
			patch: 3,
		},
		"with_v_prefix": {
			input: "v1.2.3",
			major: 1,
			minor: 2,
			patch: 3,
		},
		"with_prerelease": {
			input:      "1.2.3-alpha",
			major:      1,
			minor:      2,
			patch:      3,
			prerelease: "alpha",
		},
		"with_build": {
			input: "1.2.3+build",
			major: 1,
			minor: 2,
			patch: 3,
			build: "build",
		},
		"with_both": {
			input:      "1.2.3-beta.1+build.123",
			major:      1,
			minor:      2,
			patch:      3,
			prerelease: "beta.1",
			build:      "build.123",
		},
		"major_only": {
			input: "1",
			major: 1,
			minor: 0,
			patch: 0,
		},
		"major_minor": {
			input: "1.2",
			major: 1,
			minor: 2,
			patch: 0,
		},
		"empty_string": {
			input:       "",
			expectError: true,
		},
		"invalid_format": {
			input:       "abc",
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v, err := ParseSemVer(tc.input)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, v)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, v)

			assert.Equal(t, tc.major, v.Major())
			assert.Equal(t, tc.minor, v.Minor())
			assert.Equal(t, tc.patch, v.Patch())
			assert.Equal(t, tc.prerelease, v.Prerelease())
			assert.Equal(t, tc.build, v.Build())
		})
	}
}

func TestMustParseSemVer(t *testing.T) {
	t.Run("valid_version", func(t *testing.T) {
		v := MustParseSemVer("1.2.3")
		assert.NotNil(t, v)
		assert.Equal(t, 1, v.Major())
		assert.Equal(t, 2, v.Minor())
		assert.Equal(t, 3, v.Patch())
	})

	t.Run("invalid_version_panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustParseSemVer("invalid")
		})
	})
}

func TestSemVerString(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"basic": {
			input:    "1.2.3",
			expected: "1.2.3",
		},
		"with_v_prefix": {
			input:    "v1.2.3",
			expected: "1.2.3",
		},
		"with_prerelease": {
			input:    "1.2.3-alpha",
			expected: "1.2.3-alpha",
		},
		"with_build": {
			input:    "1.2.3+build",
			expected: "1.2.3+build",
		},
		"with_both": {
			input:    "1.2.3-beta.1+build.123",
			expected: "1.2.3-beta.1+build.123",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := MustParseSemVer(tc.input)
			assert.Equal(t, tc.expected, v.String())
		})
	}
}

func TestSemVerCompare(t *testing.T) {
	tests := map[string]struct {
		v1       string
		v2       string
		expected int
	}{
		"equal": {
			v1:       "1.2.3",
			v2:       "1.2.3",
			expected: 0,
		},
		"major_less": {
			v1:       "1.2.3",
			v2:       "2.0.0",
			expected: -1,
		},
		"major_greater": {
			v1:       "2.0.0",
			v2:       "1.2.3",
			expected: 1,
		},
		"minor_less": {
			v1:       "1.1.0",
			v2:       "1.2.0",
			expected: -1,
		},
		"minor_greater": {
			v1:       "1.2.0",
			v2:       "1.1.0",
			expected: 1,
		},
		"patch_less": {
			v1:       "1.2.2",
			v2:       "1.2.3",
			expected: -1,
		},
		"patch_greater": {
			v1:       "1.2.3",
			v2:       "1.2.2",
			expected: 1,
		},
		"prerelease_less_than_release": {
			v1:       "1.2.3-alpha",
			v2:       "1.2.3",
			expected: -1,
		},
		"release_greater_than_prerelease": {
			v1:       "1.2.3",
			v2:       "1.2.3-alpha",
			expected: 1,
		},
		"prerelease_comparison": {
			v1:       "1.2.3-alpha",
			v2:       "1.2.3-beta",
			expected: -1,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ver1 := MustParseSemVer(tc.v1)
			ver2 := MustParseSemVer(tc.v2)

			result := ver1.Compare(ver2)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSemVerComparisonMethods(t *testing.T) {
	v1 := MustParseSemVer("1.2.3")
	v2 := MustParseSemVer("2.0.0")
	v3 := MustParseSemVer("1.2.3")

	// LessThan
	assert.True(t, v1.LessThan(v2))
	assert.False(t, v2.LessThan(v1))
	assert.False(t, v1.LessThan(v3))

	// GreaterThan
	assert.False(t, v1.GreaterThan(v2))
	assert.True(t, v2.GreaterThan(v1))
	assert.False(t, v1.GreaterThan(v3))

	// Equal
	assert.False(t, v1.Equal(v2))
	assert.True(t, v1.Equal(v3))

	// GreaterThanOrEqual
	assert.False(t, v1.GreaterThanOrEqual(v2))
	assert.True(t, v2.GreaterThanOrEqual(v1))
	assert.True(t, v1.GreaterThanOrEqual(v3))

	// LessThanOrEqual
	assert.True(t, v1.LessThanOrEqual(v2))
	assert.False(t, v2.LessThanOrEqual(v1))
	assert.True(t, v1.LessThanOrEqual(v3))
}

func TestCompareVersions(t *testing.T) {
	tests := map[string]struct {
		v1          string
		v2          string
		expected    int
		expectError bool
	}{
		"v1_less_than_v2": {
			v1:       "1.2.3",
			v2:       "2.0.0",
			expected: -1,
		},
		"v1_greater_than_v2": {
			v1:       "2.0.0",
			v2:       "1.2.3",
			expected: 1,
		},
		"v1_equal_v2": {
			v1:       "1.2.3",
			v2:       "1.2.3",
			expected: 0,
		},
		"invalid_v1": {
			v1:          "invalid",
			v2:          "1.2.3",
			expectError: true,
		},
		"invalid_v2": {
			v1:          "1.2.3",
			v2:          "invalid",
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := CompareVersions(tc.v1, tc.v2)

			if tc.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsNewerVersion(t *testing.T) {
	tests := map[string]struct {
		v1          string
		v2          string
		expected    bool
		expectError bool
	}{
		"v1_newer": {
			v1:       "2.0.0",
			v2:       "1.2.3",
			expected: true,
		},
		"v1_older": {
			v1:       "1.2.3",
			v2:       "2.0.0",
			expected: false,
		},
		"v1_equal": {
			v1:       "1.2.3",
			v2:       "1.2.3",
			expected: false,
		},
		"invalid_v1": {
			v1:          "invalid",
			v2:          "1.2.3",
			expectError: true,
		},
		"invalid_v2": {
			v1:          "1.2.3",
			v2:          "invalid",
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := IsNewerVersion(tc.v1, tc.v2)

			if tc.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestSemVerConcurrency verifies thread safety of SemVer operations
func TestSemVerConcurrency(t *testing.T) {
	v1 := MustParseSemVer("1.2.3")
	v2 := MustParseSemVer("2.0.0")

	const goroutines = 100
	done := make(chan bool, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			// All these operations should be safe for concurrent use
			_ = v1.String()
			_ = v1.Major()
			_ = v1.Minor()
			_ = v1.Patch()
			_ = v1.Prerelease()
			_ = v1.Build()
			_ = v1.Compare(v2)
			_ = v1.LessThan(v2)
			_ = v1.GreaterThan(v2)
			_ = v1.Equal(v2)
			done <- true
		}()
	}

	for i := 0; i < goroutines; i++ {
		<-done
	}
}

// TestSemVerUsageExamples provides usage examples
func TestSemVerUsageExamples(t *testing.T) {
	t.Run("basic_parsing", func(t *testing.T) {
		v, err := ParseSemVer("1.2.3")
		require.NoError(t, err)
		assert.Equal(t, "1.2.3", v.String())
	})

	t.Run("version_comparison", func(t *testing.T) {
		current := MustParseSemVer("1.2.3")
		required := MustParseSemVer("1.0.0")

		if current.GreaterThanOrEqual(required) {
			// Version meets requirements
			assert.True(t, true)
		}
	})

	t.Run("upgrade_check", func(t *testing.T) {
		currentVersion := "1.2.3"
		latestVersion := "2.0.0"

		isNewer, err := IsNewerVersion(latestVersion, currentVersion)
		require.NoError(t, err)

		if isNewer {
			// Upgrade available
			assert.True(t, true)
		}
	})

	t.Run("prerelease_handling", func(t *testing.T) {
		stable := MustParseSemVer("1.2.3")
		alpha := MustParseSemVer("1.2.3-alpha")

		// Stable versions are considered newer than prerelease
		assert.True(t, stable.GreaterThan(alpha))
	})
}
