package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := map[string]struct {
		input       string
		expected    *Version
		expectError bool
	}{
		"simple_version": {
			input: "1.2.3",
			expected: &Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
		},
		"with_v_prefix": {
			input: "v1.2.3",
			expected: &Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
		},
		"with_prerelease": {
			input: "1.2.3-alpha",
			expected: &Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				Prerelease: "alpha",
			},
		},
		"with_build": {
			input: "1.2.3+build123",
			expected: &Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: "build123",
			},
		},
		"with_prerelease_and_build": {
			input: "1.2.3-beta.1+build456",
			expected: &Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				Prerelease: "beta.1",
				Build:      "build456",
			},
		},
		"major_only": {
			input: "2",
			expected: &Version{
				Major: 2,
				Minor: 0,
				Patch: 0,
			},
		},
		"major_minor": {
			input: "2.1",
			expected: &Version{
				Major: 2,
				Minor: 1,
				Patch: 0,
			},
		},
		"empty_string": {
			input:       "",
			expectError: true,
		},
		"invalid_major": {
			input:       "abc.2.3",
			expectError: true,
		},
		"invalid_minor": {
			input:       "1.abc.3",
			expectError: true,
		},
		"invalid_patch": {
			input:       "1.2.abc",
			expectError: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := Parse(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.expected.Major, result.Major)
				assert.Equal(t, tt.expected.Minor, result.Minor)
				assert.Equal(t, tt.expected.Patch, result.Patch)
				assert.Equal(t, tt.expected.Prerelease, result.Prerelease)
				assert.Equal(t, tt.expected.Build, result.Build)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	tests := map[string]struct {
		version  *Version
		expected string
	}{
		"simple": {
			version:  &Version{Major: 1, Minor: 2, Patch: 3},
			expected: "1.2.3",
		},
		"with_prerelease": {
			version:  &Version{Major: 1, Minor: 2, Patch: 3, Prerelease: "alpha"},
			expected: "1.2.3-alpha",
		},
		"with_build": {
			version:  &Version{Major: 1, Minor: 2, Patch: 3, Build: "build123"},
			expected: "1.2.3+build123",
		},
		"with_both": {
			version:  &Version{Major: 1, Minor: 2, Patch: 3, Prerelease: "beta", Build: "build456"},
			expected: "1.2.3-beta+build456",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.version.String())
		})
	}
}

func TestVersion_Compare(t *testing.T) {
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
		"major_greater": {
			v1:       "2.0.0",
			v2:       "1.9.9",
			expected: 1,
		},
		"major_less": {
			v1:       "1.0.0",
			v2:       "2.0.0",
			expected: -1,
		},
		"minor_greater": {
			v1:       "1.2.0",
			v2:       "1.1.9",
			expected: 1,
		},
		"minor_less": {
			v1:       "1.1.0",
			v2:       "1.2.0",
			expected: -1,
		},
		"patch_greater": {
			v1:       "1.2.4",
			v2:       "1.2.3",
			expected: 1,
		},
		"patch_less": {
			v1:       "1.2.3",
			v2:       "1.2.4",
			expected: -1,
		},
		"prerelease_vs_release": {
			v1:       "1.2.3",
			v2:       "1.2.3-alpha",
			expected: 1, // Release > prerelease
		},
		"prerelease_comparison": {
			v1:       "1.2.3-alpha",
			v2:       "1.2.3-beta",
			expected: -1,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			v1, err := Parse(tt.v1)
			require.NoError(t, err)

			v2, err := Parse(tt.v2)
			require.NoError(t, err)

			result := v1.Compare(v2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVersion_ComparisonMethods(t *testing.T) {
	v1, _ := Parse("1.2.3")
	v2, _ := Parse("1.2.4")
	v3, _ := Parse("1.2.3")

	// LessThan
	assert.True(t, v1.LessThan(v2))
	assert.False(t, v2.LessThan(v1))
	assert.False(t, v1.LessThan(v3))

	// GreaterThan
	assert.True(t, v2.GreaterThan(v1))
	assert.False(t, v1.GreaterThan(v2))
	assert.False(t, v1.GreaterThan(v3))

	// Equal
	assert.True(t, v1.Equal(v3))
	assert.False(t, v1.Equal(v2))

	// GreaterThanOrEqual
	assert.True(t, v2.GreaterThanOrEqual(v1))
	assert.True(t, v1.GreaterThanOrEqual(v3))
	assert.False(t, v1.GreaterThanOrEqual(v2))

	// LessThanOrEqual
	assert.True(t, v1.LessThanOrEqual(v2))
	assert.True(t, v1.LessThanOrEqual(v3))
	assert.False(t, v2.LessThanOrEqual(v1))
}

func TestVersion_ParseAndString_RoundTrip(t *testing.T) {
	inputs := []string{
		"1.2.3",
		"1.2.3-alpha",
		"1.2.3+build",
		"1.2.3-beta.1+build456",
		"2.0.0",
		"0.0.1",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			v, err := Parse(input)
			require.NoError(t, err)

			output := v.String()
			// Remove v prefix if present in input for comparison
			cleanInput := input
			if input != "" && input[0] == 'v' {
				cleanInput = input[1:]
			}
			assert.Equal(t, cleanInput, output)
		})
	}
}
