package version

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialize_Success(t *testing.T) {
	Reset()       // Clean state before test
	defer Reset() // Clean up after test

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "test-app"
  version: "1.2.3"
`)

	err := Initialize(WithEmbedded(embeddedData))

	require.NoError(t, err)

	// Verify we can get the info
	info, err := Get()
	require.NoError(t, err)
	assert.Equal(t, "test-app", info.Project.Name)
	assert.Equal(t, "1.2.3", info.Project.Version)
}

func TestInitialize_MultipleCalls(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "test-app"
  version: "1.0.0"
`)

	// First call should succeed
	err := Initialize(WithEmbedded(embeddedData))
	require.NoError(t, err)

	// Second call should error
	err = Initialize(WithEmbedded(embeddedData))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "singleton")
	assert.Contains(t, err.Error(), "Options:")
}

func TestInitialize_WithInvalidManifest(t *testing.T) {
	Reset()
	defer Reset()

	invalidData := []byte(`
manifest_version: "1.0"
project:
  name: "test-app"
  # Missing required version field
`)

	err := Initialize(WithEmbedded(invalidData))
	assert.Error(t, err)
}

func TestGet_AutoInitialize(t *testing.T) {
	Reset()
	defer Reset()

	// Don't call Initialize, just call Get
	info, err := Get()

	// Should auto-initialize with defaults
	require.NoError(t, err)
	require.NotNil(t, info)
	assert.Equal(t, "unknown", info.Project.Name)
	assert.Equal(t, "0.0.0-dev", info.Project.Version)
}

func TestGet_AfterInitialize(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "initialized-app"
  version: "2.0.0"
`)

	err := Initialize(WithEmbedded(embeddedData))
	require.NoError(t, err)

	// Get should return the initialized instance
	info, err := Get()
	require.NoError(t, err)
	assert.Equal(t, "initialized-app", info.Project.Name)
	assert.Equal(t, "2.0.0", info.Project.Version)
}

func TestGet_ReturnsErrorAfterFailedInitialize(t *testing.T) {
	Reset()
	defer Reset()

	invalidData := []byte(`invalid yaml content`)

	// Initialize with invalid data
	err := Initialize(WithEmbedded(invalidData))
	require.Error(t, err)

	// Get should return the same error
	info, err := Get()
	assert.Error(t, err)
	assert.Nil(t, info)
}

func TestMustGet_Success(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "must-app"
  version: "1.0.0"
`)

	err := Initialize(WithEmbedded(embeddedData))
	require.NoError(t, err)

	// Should not panic
	info := MustGet()
	assert.Equal(t, "must-app", info.Project.Name)
}

func TestMustGet_Panic(t *testing.T) {
	Reset()
	defer Reset()

	invalidData := []byte(`invalid yaml`)

	err := Initialize(WithEmbedded(invalidData))
	require.Error(t, err)

	// Should panic
	assert.Panics(t, func() {
		MustGet()
	})
}

func TestNew_IndependentInstances(t *testing.T) {
	Reset()
	defer Reset()

	// Create first instance
	data1 := []byte(`
manifest_version: "1.0"
project:
  name: "app1"
  version: "1.0.0"
`)

	info1, err := New(WithEmbedded(data1))
	require.NoError(t, err)

	// Create second instance (independent of first)
	data2 := []byte(`
manifest_version: "1.0"
project:
  name: "app2"
  version: "2.0.0"
`)

	info2, err := New(WithEmbedded(data2))
	require.NoError(t, err)

	// Verify they're different
	assert.Equal(t, "app1", info1.Project.Name)
	assert.Equal(t, "app2", info2.Project.Name)

	// Verify singleton is unaffected
	singleton, err := Get()
	require.NoError(t, err)
	assert.NotEqual(t, info1.Project.Name, singleton.Project.Name)
	assert.NotEqual(t, info2.Project.Name, singleton.Project.Name)
}

func TestReset(t *testing.T) {
	Reset()
	defer Reset() // Final cleanup

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "test-app"
  version: "1.0.0"
`)

	// Initialize
	err := Initialize(WithEmbedded(embeddedData))
	require.NoError(t, err)

	info1, err := Get()
	require.NoError(t, err)
	assert.Equal(t, "test-app", info1.Project.Name)

	// Reset
	Reset()

	// Should be able to initialize again with different data
	newData := []byte(`
manifest_version: "1.0"
project:
  name: "new-app"
  version: "2.0.0"
`)

	err = Initialize(WithEmbedded(newData))
	require.NoError(t, err)

	info2, err := Get()
	require.NoError(t, err)
	assert.Equal(t, "new-app", info2.Project.Name)
}

func TestGet_ConcurrentAccess(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "concurrent-app"
  version: "1.0.0"
Schemas:
  db1: "10"
  db2: "20"
`)

	err := Initialize(WithEmbedded(embeddedData))
	require.NoError(t, err)

	const numGoroutines = 100
	const numIterations = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch multiple goroutines all calling Get() concurrently
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				info, err := Get()
				assert.NoError(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, "concurrent-app", info.Project.Name)
				assert.Equal(t, "1.0.0", info.Project.Version)

				// Verify immutable access
				v, ok := info.GetSchemaVersion("db1")
				assert.True(t, ok)
				assert.Equal(t, "10", v)
			}
		}()
	}

	wg.Wait()
}

func TestGet_ConcurrentAutoInitialize(t *testing.T) {
	Reset()
	defer Reset()

	const numGoroutines = 50

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch multiple goroutines all trying to auto-initialize simultaneously
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			info, err := Get()
			assert.NoError(t, err)
			assert.NotNil(t, info)
			// All should get the same singleton instance
			assert.Equal(t, "unknown", info.Project.Name)
		}()
	}

	wg.Wait()

	// Verify only one instance was created
	info1, _ := Get()
	info2, _ := Get()
	assert.Same(t, info1, info2, "Should return same instance")
}

func TestInitialize_WithValidators(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "validated-app"
  version: "1.0.0"
Schemas:
  postgres_main: "50"
  redis_cache: "5"
`)

	err := Initialize(
		WithEmbedded(embeddedData),
		WithValidators(
			NewSchemaValidator("postgres_main", "45"),
			NewSchemaValidator("redis_cache", "3"),
		),
	)

	require.NoError(t, err)

	info, err := Get()
	require.NoError(t, err)
	assert.Equal(t, "validated-app", info.Project.Name)
}

func TestInitialize_ValidatorFailure(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "old-schema-app"
  version: "1.0.0"
Schemas:
  postgres_main: "40"
`)

	err := Initialize(
		WithEmbedded(embeddedData),
		WithValidators(
			NewSchemaValidator("postgres_main", "45"), // Require version 45
		),
	)

	// Should fail validation
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "less than required minimum")
}

func TestInitialize_MissingSchemaValidation(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "missing-schema-app"
  version: "1.0.0"
`)

	err := Initialize(
		WithEmbedded(embeddedData),
		WithValidators(
			NewSchemaValidator("nonexistent_db", "1"),
		),
	)

	// Should fail because schema doesn't exist
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestNew_WithValidators(t *testing.T) {
	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "new-validated"
  version: "1.0.0"
Schemas:
  db: "10"
`)

	info, err := New(
		WithEmbedded(embeddedData),
		WithValidators(
			NewSchemaValidator("db", "5"),
		),
	)

	require.NoError(t, err)
	assert.Equal(t, "new-validated", info.Project.Name)
}

func TestInitialize_WithMultipleOptions(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "multi-opt-app"
  version: "1.0.0"
Schemas:
  db: "10"
APIs:
  rest: "1.0.0"
Components:
  chat: "2.0.0"
`)

	err := Initialize(
		WithEmbedded(embeddedData),
		WithGitInfo(),
		WithBuildInfo(),
		WithValidators(
			NewSchemaValidator("db", "5"),
			NewAPIValidator("rest", "1.0.0"),
			NewComponentValidator("chat", "1.5.0"),
		),
	)

	require.NoError(t, err)

	info, err := Get()
	require.NoError(t, err)
	assert.Equal(t, "multi-opt-app", info.Project.Name)
	assert.NotEmpty(t, info.Build.GoVersion)
}

// Benchmark tests
func BenchmarkGet(b *testing.B) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "bench-app"
  version: "1.0.0"
`)

	_ = Initialize(WithEmbedded(embeddedData))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Get()
	}
}

func BenchmarkGet_Concurrent(b *testing.B) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "bench-app"
  version: "1.0.0"
`)

	_ = Initialize(WithEmbedded(embeddedData))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = Get()
		}
	})
}

func BenchmarkNew(b *testing.B) {
	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "bench-app"
  version: "1.0.0"
`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = New(WithEmbedded(embeddedData))
	}
}

func BenchmarkMustGet(b *testing.B) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "bench-app"
  version: "1.0.0"
`)

	_ = Initialize(WithEmbedded(embeddedData))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MustGet()
	}
}

func BenchmarkInfo_GetSchemas(b *testing.B) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "bench-app"
  version: "1.0.0"
Schemas:
  db: "45"
  cache: "12"
  messaging: "3"
  analytics: "8"
  search: "2"
`)

	_ = Initialize(WithEmbedded(embeddedData))
	info := MustGet()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = info.GetSchemas()
	}
}

func BenchmarkInfo_GetAPIs(b *testing.B) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "bench-app"
  version: "1.0.0"
APIs:
  rest: "2.0"
  graphql: "1.5"
  grpc: "1.0"
`)

	_ = Initialize(WithEmbedded(embeddedData))
	info := MustGet()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = info.GetAPIs()
	}
}

func BenchmarkInfo_MarshalJSON(b *testing.B) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "bench-app"
  version: "1.2.3"
Schemas:
  db: "45"
  cache: "12"
APIs:
  rest: "2.0"
Components:
  cache: "1.5"
`)

	_ = Initialize(WithEmbedded(embeddedData))
	info := MustGet()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = info.MarshalJSON()
	}
}

func BenchmarkInfo_String(b *testing.B) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "bench-app"
  version: "1.2.3"
`)

	_ = Initialize(WithEmbedded(embeddedData))
	info := MustGet()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = info.String()
	}
}

func BenchmarkInfo_LogFields(b *testing.B) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "bench-app"
  version: "1.2.3"
`)

	_ = Initialize(WithEmbedded(embeddedData))
	info := MustGet()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = info.LogFields()
	}
}

func BenchmarkParseSemVer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseSemVer("1.2.3-alpha+build.123")
	}
}

func BenchmarkSemVer_Compare(b *testing.B) {
	v1 := MustParseSemVer("1.2.3")
	v2 := MustParseSemVer("2.0.0")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v1.Compare(v2)
	}
}

func BenchmarkCompareVersions(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = CompareVersions("1.2.3", "2.0.0")
	}
}

func BenchmarkIsNewerVersion(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = IsNewerVersion("2.0.0", "1.2.3")
	}
}

func TestIsInitialized_BeforeInitialize(t *testing.T) {
	Reset()
	defer Reset()

	// Before any initialization
	assert.False(t, IsInitialized())
}

func TestIsInitialized_AfterSuccessfulInitialize(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "test-app"
  version: "1.0.0"
`)

	err := Initialize(WithEmbedded(embeddedData))
	require.NoError(t, err)

	// After successful initialization
	assert.True(t, IsInitialized())
}

func TestIsInitialized_AfterFailedInitialize(t *testing.T) {
	Reset()
	defer Reset()

	invalidData := []byte(`invalid yaml`)

	err := Initialize(WithEmbedded(invalidData))
	require.Error(t, err)

	// After failed initialization, should still be false
	assert.False(t, IsInitialized())
}

func TestIsInitialized_AfterAutoInitialize(t *testing.T) {
	Reset()
	defer Reset()

	// Trigger auto-initialization via Get()
	_, err := Get()
	require.NoError(t, err)

	// Should be initialized now
	assert.True(t, IsInitialized())
}

func TestIsInitialized_AfterReset(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "test-app"
  version: "1.0.0"
`)

	err := Initialize(WithEmbedded(embeddedData))
	require.NoError(t, err)
	assert.True(t, IsInitialized())

	// Reset should clear initialization
	Reset()
	assert.False(t, IsInitialized())
}

func TestIsInitialized_DoesNotTriggerAutoInit(t *testing.T) {
	Reset()
	defer Reset()

	// Call IsInitialized multiple times
	assert.False(t, IsInitialized())
	assert.False(t, IsInitialized())
	assert.False(t, IsInitialized())

	// Should still not be initialized (no auto-init)
	assert.False(t, IsInitialized())

	// Now manually initialize
	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "test-app"
  version: "1.0.0"
`)

	err := Initialize(WithEmbedded(embeddedData))
	require.NoError(t, err)

	// Now it should be true
	assert.True(t, IsInitialized())
}

func TestIsInitialized_ConcurrentAccess(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "concurrent-test"
  version: "1.0.0"
`)

	err := Initialize(WithEmbedded(embeddedData))
	require.NoError(t, err)

	const numGoroutines = 100
	const numIterations = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch multiple goroutines all calling IsInitialized() concurrently
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				assert.True(t, IsInitialized())
			}
		}()
	}

	wg.Wait()
}

func TestWithStrictMode_MissingManifest(t *testing.T) {
	Reset()
	defer Reset()

	// In strict mode, missing manifest should be fatal
	err := Initialize(
		WithManifestPath("/nonexistent/path/versions.yaml"),
		WithStrictMode(),
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "strict mode")
	assert.Contains(t, err.Error(), "manifest file is required")
}

func TestWithStrictMode_InvalidManifest(t *testing.T) {
	Reset()
	defer Reset()

	invalidData := []byte(`
manifest_version: "1.0"
project:
  name: "test-app"
  # Missing required version field
`)

	// In strict mode, invalid manifest should be fatal
	err := Initialize(
		WithEmbedded(invalidData),
		WithStrictMode(),
	)

	assert.Error(t, err)
}

func TestWithStrictMode_ValidManifest(t *testing.T) {
	Reset()
	defer Reset()

	validData := []byte(`
manifest_version: "1.0"
project:
  name: "strict-app"
  version: "1.0.0"
`)

	err := Initialize(
		WithEmbedded(validData),
		WithStrictMode(),
	)

	require.NoError(t, err)

	info, err := Get()
	require.NoError(t, err)
	assert.Equal(t, "strict-app", info.Project.Name)
}

func TestWithStrictMode_WithValidators(t *testing.T) {
	Reset()
	defer Reset()

	validData := []byte(`
manifest_version: "1.0"
project:
  name: "strict-validated-app"
  version: "1.0.0"
Schemas:
  postgres_main: "50"
`)

	err := Initialize(
		WithEmbedded(validData),
		WithStrictMode(),
		WithValidators(
			NewSchemaValidator("postgres_main", "45"),
		),
	)

	require.NoError(t, err)

	info, err := Get()
	require.NoError(t, err)
	assert.Equal(t, "strict-validated-app", info.Project.Name)
}

func TestWithStrictMode_ValidatorFailure(t *testing.T) {
	Reset()
	defer Reset()

	validData := []byte(`
manifest_version: "1.0"
project:
  name: "strict-fail-app"
  version: "1.0.0"
Schemas:
  postgres_main: "40"
`)

	err := Initialize(
		WithEmbedded(validData),
		WithStrictMode(),
		WithValidators(
			NewSchemaValidator("postgres_main", "45"),
		),
	)

	// Should fail validation even in strict mode
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "less than required minimum")
}

func TestWithoutStrictMode_MissingManifest(t *testing.T) {
	Reset()
	defer Reset()

	// Without strict mode, missing manifest should use defaults
	err := Initialize(
		WithManifestPath("/nonexistent/path/versions.yaml"),
	)

	// Should succeed with defaults
	require.NoError(t, err)

	info, err := Get()
	require.NoError(t, err)
	assert.Equal(t, "unknown", info.Project.Name)
	assert.Equal(t, "0.0.0-dev", info.Project.Version)
}

func TestNew_WithStrictMode(t *testing.T) {
	validData := []byte(`
manifest_version: "1.0"
project:
  name: "new-strict-app"
  version: "2.0.0"
`)

	info, err := New(
		WithEmbedded(validData),
		WithStrictMode(),
	)

	require.NoError(t, err)
	assert.Equal(t, "new-strict-app", info.Project.Name)
}

func TestNew_WithStrictMode_Missing(t *testing.T) {
	info, err := New(
		WithManifestPath("/nonexistent/strict.yaml"),
		WithStrictMode(),
	)

	assert.Error(t, err)
	assert.Nil(t, info)
	assert.Contains(t, err.Error(), "strict mode")
}

// TestInitialize_WithContext verifies context propagation to validators
func TestInitialize_WithContext(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "test-app"
  version: "1.0.0"
Schemas:
  db: "10"
`)

	// Create a validator that checks for context
	contextReceived := false
	validator := ValidatorFunc(func(ctx context.Context, info *Info) error {
		contextReceived = true
		// Verify context is not nil
		assert.NotNil(t, ctx)
		// Check if we can read values from context
		if val := ctx.Value("test_key"); val != nil {
			assert.Equal(t, "test_value", val.(string))
		}
		return nil
	})

	// Create context with value
	ctx := context.WithValue(context.Background(), "test_key", "test_value")

	err := Initialize(
		WithContext(ctx),
		WithEmbedded(embeddedData),
		WithValidators(validator),
	)

	require.NoError(t, err)
	assert.True(t, contextReceived, "Validator should have been called with context")
}

// TestInitialize_WithContextCancellation verifies context cancellation works
func TestInitialize_WithContextCancellation(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "test-app"
  version: "1.0.0"
`)

	// Create a validator that checks for cancellation
	validator := ValidatorFunc(func(ctx context.Context, info *Info) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	})

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := Initialize(
		WithContext(ctx),
		WithEmbedded(embeddedData),
		WithValidators(validator),
	)

	// Should get cancellation error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}

// TestNew_WithContext verifies context works with non-singleton New()
func TestNew_WithContext(t *testing.T) {
	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "test-app"
  version: "1.0.0"
`)

	contextReceived := false
	validator := ValidatorFunc(func(ctx context.Context, info *Info) error {
		contextReceived = true
		assert.NotNil(t, ctx)
		return nil
	})

	ctx := context.Background()
	info, err := New(
		WithContext(ctx),
		WithEmbedded(embeddedData),
		WithValidators(validator),
	)

	require.NoError(t, err)
	assert.NotNil(t, info)
	assert.True(t, contextReceived)
}
