package version

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInfo_GetSchemaVersion(t *testing.T) {
	tests := map[string]struct {
		info     *Info
		name     string
		expected string
		found    bool
	}{
		"existing_schema": {
			info: &Info{
				schemas: map[string]string{
					"postgres_main": "47",
					"redis_cache":   "3",
				},
			},
			name:     "postgres_main",
			expected: "47",
			found:    true,
		},
		"non_existing_schema": {
			info: &Info{
				schemas: map[string]string{
					"postgres_main": "47",
				},
			},
			name:     "postgres_analytics",
			expected: "",
			found:    false,
		},
		"nil_schemas_map": {
			info:     &Info{},
			name:     "postgres_main",
			expected: "",
			found:    false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			version, found := tt.info.GetSchemaVersion(tt.name)
			assert.Equal(t, tt.expected, version)
			assert.Equal(t, tt.found, found)
		})
	}
}

func TestInfo_GetAPIVersion(t *testing.T) {
	tests := map[string]struct {
		info     *Info
		name     string
		expected string
		found    bool
	}{
		"existing_api": {
			info: &Info{
				apis: map[string]string{
					"rest_v1": "1.15.0",
					"grpc":    "1.2.0",
				},
			},
			name:     "rest_v1",
			expected: "1.15.0",
			found:    true,
		},
		"non_existing_api": {
			info: &Info{
				apis: map[string]string{
					"rest_v1": "1.15.0",
				},
			},
			name:     "rest_v2",
			expected: "",
			found:    false,
		},
		"nil_apis_map": {
			info:     &Info{},
			name:     "rest_v1",
			expected: "",
			found:    false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			version, found := tt.info.GetAPIVersion(tt.name)
			assert.Equal(t, tt.expected, version)
			assert.Equal(t, tt.found, found)
		})
	}
}

func TestInfo_GetComponentVersion(t *testing.T) {
	tests := map[string]struct {
		info     *Info
		name     string
		expected string
		found    bool
	}{
		"existing_component": {
			info: &Info{
				components: map[string]string{
					"aigentchat": "3.4.1",
					"hyperrag":   "2.0.5",
				},
			},
			name:     "aigentchat",
			expected: "3.4.1",
			found:    true,
		},
		"non_existing_component": {
			info: &Info{
				components: map[string]string{
					"aigentchat": "3.4.1",
				},
			},
			name:     "aigentflow",
			expected: "",
			found:    false,
		},
		"nil_components_map": {
			info:     &Info{},
			name:     "aigentchat",
			expected: "",
			found:    false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			version, found := tt.info.GetComponentVersion(tt.name)
			assert.Equal(t, tt.expected, version)
			assert.Equal(t, tt.found, found)
		})
	}
}

func TestInfo_LoadedAt(t *testing.T) {
	now := time.Now()
	info := &Info{
		loadedAt: now,
	}

	assert.Equal(t, now, info.LoadedAt())
}

func TestInfo_LogFields(t *testing.T) {
	info := &Info{
		Project: ProjectVersion{
			Name:    "test-app",
			Version: "1.2.3",
		},
		Git: GitInfo{
			Commit:     "abc123",
			Tag:        "v1.2.3",
			TreeState:  "clean",
			CommitTime: "2025-10-11T10:00:00Z",
		},
		Build: BuildInfo{
			Time:      "2025-10-11T12:00:00Z",
			User:      "ci",
			GoVersion: "go1.21.0",
		},
	}

	fields := info.LogFields()

	require.Len(t, fields, 8)

	// Verify all fields are present (zap.Field doesn't expose values easily, so we just check length)
	assert.NotNil(t, fields)
}

func TestInfo_String(t *testing.T) {
	info := &Info{
		Project: ProjectVersion{
			Name:    "test-app",
			Version: "1.2.3",
		},
		Git: GitInfo{
			Commit: "abc123",
		},
	}

	expected := "test-app 1.2.3 (abc123)"
	assert.Equal(t, expected, info.String())
}

func TestInfo_MarshalJSON(t *testing.T) {
	now := time.Now()
	info := &Info{
		Project: ProjectVersion{
			Name:    "test-app",
			Version: "1.2.3",
		},
		Git: GitInfo{
			Commit:    "abc123",
			Tag:       "v1.2.3",
			TreeState: "clean",
		},
		Build: BuildInfo{
			Time:      "2025-10-11T12:00:00Z",
			GoVersion: "go1.21.0",
		},
		schemas: map[string]string{
			"db": "47",
		},
		loadedAt: now,
	}

	data, err := json.Marshal(info)
	require.NoError(t, err)

	// Parse back to map to verify structure
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// Verify main fields are present
	assert.Contains(t, result, "project")
	assert.Contains(t, result, "git")
	assert.Contains(t, result, "build")
	assert.Contains(t, result, "schemas")

	// Verify loadedAt is NOT in JSON (internal field)
	assert.NotContains(t, result, "loadedAt")
	assert.NotContains(t, result, "LoadedAt")
}

func TestInfo_JSONRoundTrip(t *testing.T) {
	original := &Info{
		Project: ProjectVersion{
			Name:    "test-app",
			Version: "1.2.3",
		},
		Git: GitInfo{
			Commit:     "abc123def456",
			Tag:        "v1.2.3",
			TreeState:  "clean",
			CommitTime: "2025-10-11T10:00:00Z",
		},
		Build: BuildInfo{
			Time:      "2025-10-11T12:00:00Z",
			User:      "jenkins",
			GoVersion: "go1.21.0",
		},
		schemas: map[string]string{
			"postgres_main": "47",
			"redis_cache":   "3",
		},
		apis: map[string]string{
			"rest_v1": "1.15.0",
		},
		components: map[string]string{
			"aigentchat": "3.4.1",
		},
		custom: map[string]interface{}{
			"feature_flags": "2024.10",
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	require.NoError(t, err)

	// Unmarshal back
	var decoded Info
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	// Compare (excluding loadedAt which is not serialized)
	assert.Equal(t, original.Project, decoded.Project)
	assert.Equal(t, original.Git, decoded.Git)
	assert.Equal(t, original.Build, decoded.Build)
	assert.Equal(t, original.GetSchemas(), decoded.GetSchemas())
	assert.Equal(t, original.GetAPIs(), decoded.GetAPIs())
	assert.Equal(t, original.GetComponents(), decoded.GetComponents())
}

// TestInfo_ConcurrentReads tests that multiple goroutines can safely read from Info
func TestInfo_ConcurrentReads(t *testing.T) {
	info := &Info{
		Project: ProjectVersion{
			Name:    "test-app",
			Version: "1.2.3",
		},
		schemas: map[string]string{
			"db1": "1",
			"db2": "2",
			"db3": "3",
		},
		apis: map[string]string{
			"api1": "1.0.0",
			"api2": "2.0.0",
		},
		components: map[string]string{
			"comp1": "1.0.0",
			"comp2": "2.0.0",
		},
	}

	const numGoroutines = 100
	const numIterations = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Start multiple goroutines reading concurrently
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				// Read schema
				v, ok := info.GetSchemaVersion("db1")
				assert.True(t, ok)
				assert.Equal(t, "1", v)

				// Read API
				v, ok = info.GetAPIVersion("api1")
				assert.True(t, ok)
				assert.Equal(t, "1.0.0", v)

				// Read component
				v, ok = info.GetComponentVersion("comp1")
				assert.True(t, ok)
				assert.Equal(t, "1.0.0", v)

				// Call other methods
				_ = info.String()
				_ = info.LoadedAt()
				_ = info.LogFields()
			}
		}()
	}

	wg.Wait()
}

// TestInfo_ConcurrentReadsWithJSON tests concurrent JSON marshaling
func TestInfo_ConcurrentReadsWithJSON(t *testing.T) {
	info := &Info{
		Project: ProjectVersion{
			Name:    "test-app",
			Version: "1.2.3",
		},
		Git: GitInfo{
			Commit: "abc123",
		},
		schemas: map[string]string{
			"db": "47",
		},
	}

	const numGoroutines = 50

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				data, err := json.Marshal(info)
				assert.NoError(t, err)
				assert.NotNil(t, data)
			}
		}()
	}

	wg.Wait()
}

// Example test for documentation
func ExampleInfo_GetSchemaVersion() {
	info := &Info{
		schemas: map[string]string{
			"postgres_main": "47",
		},
	}

	version, found := info.GetSchemaVersion("postgres_main")
	if found {
		// Use fmt.Println instead of println for example tests
		fmt.Println("Schema version:", version)
	}
	// Output: Schema version: 47
}
