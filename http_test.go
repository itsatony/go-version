package version

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_Success(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "http-test-app"
  version: "1.2.3"
schemas:
  db: "10"
`)

	err := Initialize(WithEmbedded(embeddedData))
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/version", http.NoBody)
	w := httptest.NewRecorder()

	handler := Handler()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, HTTPContentTypeJSON, w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Cache-Control"), "max-age")

	var info Info
	err = json.Unmarshal(w.Body.Bytes(), &info)
	require.NoError(t, err)
	assert.Equal(t, "http-test-app", info.Project.Name)
	assert.Equal(t, "1.2.3", info.Project.Version)
	assert.Equal(t, "10", info.GetSchemas()["db"])
}

func TestHandler_AutoInitialize(t *testing.T) {
	Reset()
	defer Reset()

	req := httptest.NewRequest(http.MethodGet, "/version", http.NoBody)
	w := httptest.NewRecorder()

	handler := Handler()
	handler.ServeHTTP(w, req)

	// Should auto-initialize and succeed
	assert.Equal(t, http.StatusOK, w.Code)

	var info Info
	err := json.Unmarshal(w.Body.Bytes(), &info)
	require.NoError(t, err)
	assert.Equal(t, "unknown", info.Project.Name)
	assert.Equal(t, "0.0.0-dev", info.Project.Version)
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	Reset()
	defer Reset()

	req := httptest.NewRequest(http.MethodPost, "/version", http.NoBody)
	w := httptest.NewRecorder()

	handler := Handler()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestHandler_InitializationError(t *testing.T) {
	Reset()
	defer Reset()

	// Initialize with invalid data
	invalidData := []byte(`invalid yaml`)
	err := Initialize(WithEmbedded(invalidData))
	require.Error(t, err)

	req := httptest.NewRequest(http.MethodGet, "/version", http.NoBody)
	w := httptest.NewRecorder()

	handler := Handler()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHealthHandler_Success(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "health-test-app"
  version: "2.3.4"
`)

	err := Initialize(WithEmbedded(embeddedData))
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
	w := httptest.NewRecorder()

	handler := HealthHandler()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, HTTPContentTypeJSON, w.Header().Get("Content-Type"))

	var response struct {
		Status    string `json:"status"`
		Version   string `json:"version"`
		Timestamp string `json:"timestamp"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "ok", response.Status)
	assert.Equal(t, "2.3.4", response.Version)
	assert.NotEmpty(t, response.Timestamp)
}

func TestHealthHandler_Error(t *testing.T) {
	Reset()
	defer Reset()

	// Initialize with invalid data
	invalidData := []byte(`invalid yaml`)
	err := Initialize(WithEmbedded(invalidData))
	require.Error(t, err)

	req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
	w := httptest.NewRecorder()

	handler := HealthHandler()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "error", response.Status)
	assert.Contains(t, response.Error, "not available")
}

func TestHealthHandler_MethodNotAllowed(t *testing.T) {
	Reset()
	defer Reset()

	req := httptest.NewRequest(http.MethodPost, "/health", http.NoBody)
	w := httptest.NewRecorder()

	handler := HealthHandler()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestHandlerFunc(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "func-test"
  version: "1.0.0"
`)

	err := Initialize(WithEmbedded(embeddedData))
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/version", http.NoBody)
	w := httptest.NewRecorder()

	// Use HandlerFunc instead of Handler
	handlerFunc := HandlerFunc()
	handlerFunc(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var info Info
	err = json.Unmarshal(w.Body.Bytes(), &info)
	require.NoError(t, err)
	assert.Equal(t, "func-test", info.Project.Name)
}

func TestHealthHandlerFunc(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "health-func-test"
  version: "1.0.0"
`)

	err := Initialize(WithEmbedded(embeddedData))
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
	w := httptest.NewRecorder()

	// Use HealthHandlerFunc instead of HealthHandler
	handlerFunc := HealthHandlerFunc()
	handlerFunc(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  string `json:"status"`
		Version string `json:"version"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "ok", response.Status)
	assert.Equal(t, "1.0.0", response.Version)
}

func TestMiddleware(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "middleware-test"
  version: "3.2.1"
`)

	err := Initialize(WithEmbedded(embeddedData))
	require.NoError(t, err)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test response"))
	})

	// Wrap with middleware
	wrappedHandler := Middleware(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	w := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "3.2.1", w.Header().Get("X-App-Version"))
	// Git commit header should be present (will be "dev" or actual commit)
	assert.NotEmpty(t, w.Header().Get("X-Git-Commit"))
	assert.Equal(t, "test response", w.Body.String())
}

func TestMiddleware_NoGitCommit(t *testing.T) {
	Reset()
	defer Reset()

	embeddedData := []byte(`
manifest_version: "1.0"
project:
  name: "middleware-test"
  version: "1.0.0"
`)

	// Initialize without git info
	err := Initialize(WithEmbedded(embeddedData), WithoutGitInfo())
	require.NoError(t, err)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := Middleware(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	w := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "1.0.0", w.Header().Get("X-App-Version"))
	// Since git commit is "dev" (default), it should still be set
	// The middleware checks if it's not DefaultGitCommit before setting
	if w.Header().Get("X-Git-Commit") != "" {
		assert.NotEqual(t, DefaultGitCommit, w.Header().Get("X-Git-Commit"))
	}
}

func TestMiddleware_InitializationFailure(t *testing.T) {
	Reset()
	defer Reset()

	// Initialize with invalid data
	invalidData := []byte(`invalid yaml`)
	err := Initialize(WithEmbedded(invalidData))
	require.Error(t, err)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test response"))
	})

	wrappedHandler := Middleware(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	w := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, req)

	// Middleware should not block request on failure
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test response", w.Body.String())
	// Headers should not be set
	assert.Empty(t, w.Header().Get("X-App-Version"))
}

func TestHandler_ConcurrentRequests(t *testing.T) {
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

	handler := Handler()

	// Make 50 concurrent requests
	const numRequests = 50
	done := make(chan bool, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, "/version", http.NoBody)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
			done <- true
		}()
	}

	// Wait for all requests to complete
	for i := 0; i < numRequests; i++ {
		<-done
	}
}

// Example demonstrating HTTP handler usage
func ExampleHandler() {
	// Initialize with custom version info
	_ = Initialize(WithManifestPath("versions.yaml"))

	// Create a mux and register the version handler
	mux := http.NewServeMux()
	mux.Handle("/version", Handler())

	// Start server (example only, not actually starting)
	// http.ListenAndServe(":8080", mux)
}

// Example demonstrating health check handler
func ExampleHealthHandler() {
	mux := http.NewServeMux()
	mux.Handle("/health", HealthHandler())
	mux.Handle("/version", Handler())

	// Wrap with middleware to add version headers to all responses
	_ = Middleware(mux)

	// Start server (example only, not actually starting)
	// http.ListenAndServe(":8080", Middleware(mux))
}
