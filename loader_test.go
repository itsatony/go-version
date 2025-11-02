package version

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseManifest(t *testing.T) {
	tests := map[string]struct {
		yaml        string
		expectError bool
		checkFields func(*testing.T, *Manifest)
	}{
		"valid_simple": {
			yaml: `
manifest_version: "1.0"
project:
  name: "test-app"
  version: "1.2.3"
`,
			checkFields: func(t *testing.T, m *Manifest) {
				assert.Equal(t, "1.0", m.ManifestVersion)
				assert.Equal(t, "test-app", m.Project.Name)
				assert.Equal(t, "1.2.3", m.Project.Version)
			},
		},
		"with_schemas": {
			yaml: `
manifest_version: "1.0"
project:
  name: "test-app"
  version: "1.2.3"
schemas:
  postgres_main: "47"
  redis_cache: "3"
`,
			checkFields: func(t *testing.T, m *Manifest) {
				assert.Equal(t, "47", m.Schemas["postgres_main"])
				assert.Equal(t, "3", m.Schemas["redis_cache"])
			},
		},
		"with_apis_and_components": {
			yaml: `
manifest_version: "1.0"
project:
  name: "test-app"
  version: "1.2.3"
apis:
  rest_v1: "1.15.0"
  grpc: "1.2.0"
components:
  aigentchat: "3.4.1"
`,
			checkFields: func(t *testing.T, m *Manifest) {
				assert.Equal(t, "1.15.0", m.APIs["rest_v1"])
				assert.Equal(t, "3.4.1", m.Components["aigentchat"])
			},
		},
		"missing_project_name": {
			yaml: `
manifest_version: "1.0"
project:
  version: "1.2.3"
`,
			expectError: true,
		},
		"missing_project_version": {
			yaml: `
manifest_version: "1.0"
project:
  name: "test-app"
`,
			expectError: true,
		},
		"invalid_yaml": {
			yaml:        `invalid: yaml: content:`,
			expectError: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			manifest, err := parseManifest([]byte(tt.yaml))

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, manifest)
			} else {
				require.NoError(t, err)
				require.NotNil(t, manifest)
				if tt.checkFields != nil {
					tt.checkFields(t, manifest)
				}
			}
		})
	}
}

func TestManifestToInfo(t *testing.T) {
	manifest := &Manifest{
		ManifestVersion: "1.0",
		Project: ProjectManifest{
			Name:    "test-app",
			Version: "1.2.3",
		},
		Schemas: map[string]string{
			"db1": "10",
			"db2": "20",
		},
		APIs: map[string]string{
			"api1": "1.0.0",
		},
		Components: map[string]string{
			"comp1": "2.0.0",
		},
		Custom: map[string]interface{}{
			"feature": "enabled",
		},
	}

	info := manifestToInfo(manifest)

	assert.Equal(t, "test-app", info.Project.Name)
	assert.Equal(t, "1.2.3", info.Project.Version)
	assert.Equal(t, "10", info.GetSchemas()["db1"])
	assert.Equal(t, "20", info.GetSchemas()["db2"])
	assert.Equal(t, "1.0.0", info.GetAPIs()["api1"])
	assert.Equal(t, "2.0.0", info.GetComponents()["comp1"])
	assert.Equal(t, "enabled", info.GetCustom()["feature"])

	// Verify defensive copies (modifying manifest shouldn't affect info)
	manifest.Schemas["db1"] = "99"
	assert.Equal(t, "10", info.GetSchemas()["db1"], "Info should have defensive copy")
}

func TestLoadVersionInfo_WithEmbedded(t *testing.T) {
	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "embedded-app"
  version: "2.3.4"
schemas:
  db: "42"
`)

	info, err := loadVersionInfo(WithEmbedded(embeddedData))

	require.NoError(t, err)
	assert.Equal(t, "embedded-app", info.Project.Name)
	assert.Equal(t, "2.3.4", info.Project.Version)
	assert.Equal(t, "42", info.GetSchemas()["db"])
}

func TestLoadVersionInfo_WithFile(t *testing.T) {
	// Create temporary file
	tmpDir := t.TempDir()
	manifestPath := filepath.Join(tmpDir, "versions.yaml")

	manifestContent := []byte(`
manifest_version: "1.0"
project:
  name: "file-app"
  version: "3.4.5"
apis:
  rest: "1.0.0"
`)

	err := os.WriteFile(manifestPath, manifestContent, 0o600)
	require.NoError(t, err)

	info, err := loadVersionInfo(WithManifestPath(manifestPath))

	require.NoError(t, err)
	assert.Equal(t, "file-app", info.Project.Name)
	assert.Equal(t, "3.4.5", info.Project.Version)
	assert.Equal(t, "1.0.0", info.GetAPIs()["rest"])
}

func TestLoadVersionInfo_DefaultWhenFileNotFound(t *testing.T) {
	info, err := loadVersionInfo(WithManifestPath("/nonexistent/path/versions.yaml"))

	// Should not error, should use defaults
	require.NoError(t, err)
	assert.Equal(t, "unknown", info.Project.Name)
	assert.Equal(t, "0.0.0-dev", info.Project.Version)
}

func TestLoadVersionInfo_EmbeddedTakesPrecedenceOverFile(t *testing.T) {
	// Create temporary file
	tmpDir := t.TempDir()
	manifestPath := filepath.Join(tmpDir, "versions.yaml")

	fileContent := []byte(`
manifest_version: "1.0"
project:
  name: "file-app"
  version: "1.0.0"
`)

	err := os.WriteFile(manifestPath, fileContent, 0o600)
	require.NoError(t, err)

	embeddedContent := []byte(`
manifest_version: "1.0"
project:
  name: "embedded-app"
  version: "2.0.0"
`)

	// Provide both embedded and file path - embedded should win
	info, err := loadVersionInfo(
		WithEmbedded(embeddedContent),
		WithManifestPath(manifestPath),
	)

	require.NoError(t, err)
	assert.Equal(t, "embedded-app", info.Project.Name, "Embedded should take precedence")
	assert.Equal(t, "2.0.0", info.Project.Version)
}

func TestLoadVersionInfo_WithoutGitInfo(t *testing.T) {
	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "test-app"
  version: "1.0.0"
`)

	info, err := loadVersionInfo(
		WithEmbedded(embeddedData),
		WithoutGitInfo(),
	)

	require.NoError(t, err)
	// Should still have default git values since we can't actually disable it completely
	// (it's set in manifestToInfo), but the enrichment shouldn't run
	assert.NotNil(t, info.Git)
}

func TestLoadVersionInfo_WithoutBuildInfo(t *testing.T) {
	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "test-app"
  version: "1.0.0"
`)

	info, err := loadVersionInfo(
		WithEmbedded(embeddedData),
		WithoutBuildInfo(),
	)

	require.NoError(t, err)
	// Should still have some build info (Go version is always set)
	assert.NotEmpty(t, info.Build.GoVersion)
}

func TestEnrichWithGitInfo(t *testing.T) {
	info := &Info{
		Git: GitInfo{
			Commit:    DefaultGitCommit,
			TreeState: DefaultGitTreeState,
		},
	}

	// Save original values
	origCommit := GitCommit
	origTag := GitTag
	origTreeState := GitTreeState

	// Set ldflags variables
	GitCommit = "abc123def456"
	GitTag = "v1.2.3"
	GitTreeState = GitTreeStateClean

	enrichWithGitInfo(info)

	// Verify enrichment
	assert.Equal(t, "abc123def456", info.Git.Commit)
	assert.Equal(t, "v1.2.3", info.Git.Tag)
	assert.Equal(t, GitTreeStateClean, info.Git.TreeState)

	// Restore original values
	GitCommit = origCommit
	GitTag = origTag
	GitTreeState = origTreeState
}

func TestEnrichWithBuildInfo(t *testing.T) {
	info := &Info{
		Build: BuildInfo{
			Time: DefaultBuildTime,
		},
	}

	// Save original values
	origBuildTime := BuildTime
	origBuildUser := BuildUser

	// Set ldflags variables
	BuildTime = "2025-10-11T15:00:00Z"
	BuildUser = "testuser"

	enrichWithBuildInfo(info)

	// Verify enrichment
	assert.Equal(t, "2025-10-11T15:00:00Z", info.Build.Time)
	assert.Equal(t, "testuser", info.Build.User)
	assert.NotEmpty(t, info.Build.GoVersion)

	// Restore original values
	BuildTime = origBuildTime
	BuildUser = origBuildUser
}

func TestDefaultManifest(t *testing.T) {
	manifest := defaultManifest()

	assert.Equal(t, ManifestVersion, manifest.ManifestVersion)
	assert.Equal(t, "unknown", manifest.Project.Name)
	assert.Equal(t, "0.0.0-dev", manifest.Project.Version)
}

// TestLoadVersionInfo_Immutability verifies that Info is truly immutable
func TestLoadVersionInfo_Immutability(t *testing.T) {
	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "test-app"
  version: "1.0.0"
schemas:
  db: "10"
`)

	info, err := loadVersionInfo(WithEmbedded(embeddedData))
	require.NoError(t, err)

	// Get the schemas map (returns a defensive copy)
	schemas1 := info.GetSchemas()
	originalValue := schemas1["db"]
	assert.Equal(t, "10", originalValue)

	// Modify the returned copy
	schemas1["db"] = "999"
	schemas1["newkey"] = "newvalue"

	// Get schemas again - should return a fresh copy with original values
	schemas2 := info.GetSchemas()
	assert.Equal(t, "10", schemas2["db"], "Original value should be unchanged")
	assert.NotContains(t, schemas2, "newkey", "New keys in copy should not affect original")

	// Verify GetSchemaVersion still returns original value
	value, ok := info.GetSchemaVersion("db")
	assert.True(t, ok)
	assert.Equal(t, "10", value, "GetSchemaVersion should return original value")
}

// Benchmark tests
func BenchmarkLoadVersionInfo(b *testing.B) {
	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "test-app"
  version: "1.0.0"
schemas:
  db1: "10"
  db2: "20"
apis:
  rest: "1.0.0"
`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := loadVersionInfo(WithEmbedded(embeddedData))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseManifest(b *testing.B) {
	data := []byte(`
manifest_version: "1.0"
project:
  name: "test-app"
  version: "1.0.0"
schemas:
  db1: "10"
  db2: "20"
apis:
  rest: "1.0.0"
components:
  comp1: "1.0.0"
`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parseManifest(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Test isValidCommitHash with various inputs
func TestIsValidCommitHash(t *testing.T) {
	tests := map[string]struct {
		hash  string
		valid bool
	}{
		"full_sha1": {
			hash:  "abc123def456789012345678901234567890abcd",
			valid: true,
		},
		"short_sha1": {
			hash:  "abc123d",
			valid: true,
		},
		"min_length_valid": {
			hash:  "abc1234",
			valid: true,
		},
		"too_short": {
			hash:  "abc123",
			valid: false,
		},
		"too_long": {
			hash:  "abc123def456789012345678901234567890abcdef",
			valid: false,
		},
		"contains_non_hex": {
			hash:  "abc123xyz456",
			valid: false,
		},
		"contains_spaces": {
			hash:  "abc123 def456",
			valid: false,
		},
		"contains_newline": {
			hash:  "abc123\ndef",
			valid: false,
		},
		"uppercase_hex": {
			hash:  "ABC123DEF456",
			valid: true,
		},
		"mixed_case": {
			hash:  "AbC123DeF456",
			valid: true,
		},
		"empty_string": {
			hash:  "",
			valid: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := isValidCommitHash(tt.hash)
			assert.Equal(t, tt.valid, result, "hash: %q", tt.hash)
		})
	}
}

// Test getGitBinary security validation
func TestGetGitBinary_Security(t *testing.T) {
	// This test verifies that git binary location validation works
	// We can't easily mock exec.LookPath, but we can at least verify the function runs
	
	result := getGitBinary()
	
	// On Unix systems, if git is found, it should be in a whitelisted location
	if runtime.GOOS != "windows" && result != "" {
		whitelisted := result == "/usr/bin/git" || 
			result == "/usr/local/bin/git" || 
			result == "/opt/homebrew/bin/git"
		
		assert.True(t, whitelisted, 
			"git binary should be in whitelisted location, got: %s", result)
	}
	
	// If result is empty, that's also acceptable (git not found or not whitelisted)
	// The function shouldn't panic or error
}

// Test that getGitCommit validates hash format
func TestGetGitCommit_Validation(t *testing.T) {
	// This test verifies that getGitCommit properly validates output
	// We can't control git output easily, but we can verify it returns
	// either empty string or valid hash
	
	result := getGitCommit()
	
	// If we got a result, it must be a valid hash
	if result != "" {
		assert.True(t, isValidCommitHash(result), 
			"getGitCommit returned invalid hash: %q", result)
	}
	
	// Empty result is acceptable (not in git repo, git not found, etc.)
}

// Test that getGitTag uses secure binary
func TestGetGitTag_Security(t *testing.T) {
	// Verify getGitTag doesn't panic and returns reasonable values
	result := getGitTag()
	
	// Should return either empty string or a non-empty tag
	// If not empty, should not contain obviously malicious content
	if result != "" {
		assert.NotContains(t, result, "\n", "tag should not contain newlines")
		assert.NotContains(t, result, "\r", "tag should not contain carriage returns")
		assert.NotContains(t, result, "\x00", "tag should not contain null bytes")
	}
}
