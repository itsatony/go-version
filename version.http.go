package version

import (
	"encoding/json"
	"net/http"
	"time"
)

// Handler returns an http.Handler that serves version information as JSON.
// The handler responds to GET requests with the current version info.
//
// If the version singleton is not initialized, it will auto-initialize with defaults.
// If initialization fails, returns 500 Internal Server Error.
//
// Response format matches Info.MarshalJSON() output.
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	mux := http.NewServeMux()
//	mux.Handle("/version", version.Handler())
//	http.ListenAndServe(":8080", mux)
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, HTTPErrorMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}

		// Defensive: limit request body size even for GET (defense in depth)
		r.Body = http.MaxBytesReader(w, r.Body, 1024)

		info, err := Get()
		if err != nil {
			http.Error(w, HTTPErrorVersionUnavailable, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", HTTPContentTypeJSON)
		w.Header().Set("Cache-Control", HTTPCacheControl)
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(info); err != nil {
			// If encoding fails mid-stream, log but can't send error to client
			// because headers are already sent
			return
		}
	})
}

// HealthHandler returns an http.Handler that serves a health check endpoint.
// The handler checks if version info is available and responds accordingly.
//
// Returns 200 OK if version info is available, 503 Service Unavailable otherwise.
//
// Response format:
//
//	{
//	  "status": "ok",
//	  "version": "1.2.3",
//	  "timestamp": "2025-01-15T10:30:00Z"
//	}
//
// Or on failure:
//
//	{
//	  "status": "error",
//	  "error": "version not available",
//	  "timestamp": "2025-01-15T10:30:00Z"
//	}
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	mux := http.NewServeMux()
//	mux.Handle("/health", version.HealthHandler())
//	http.ListenAndServe(":8080", mux)
func HealthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, HTTPErrorMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}

		// Defensive: limit request body size even for GET (defense in depth)
		r.Body = http.MaxBytesReader(w, r.Body, 1024)

		type healthResponse struct {
			Status    string    `json:"status"`
			Version   string    `json:"version,omitempty"`
			Error     string    `json:"error,omitempty"`
			Timestamp time.Time `json:"timestamp"`
		}

		info, err := Get()
		timestamp := time.Now().UTC()

		w.Header().Set("Content-Type", HTTPContentTypeJSON)

		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(healthResponse{
				Status:    HTTPStatusError,
				Error:     HTTPHealthErrorMessage,
				Timestamp: timestamp,
			})
			// Note: Cannot send HTTP error if encoding fails after WriteHeader
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(healthResponse{
			Status:    HTTPStatusOK,
			Version:   info.Project.Version,
			Timestamp: timestamp,
		})
		// Note: Cannot send HTTP error if encoding fails after WriteHeader
	})
}

// HandlerFunc is a convenience function that returns an http.HandlerFunc
// instead of http.Handler. It's equivalent to Handler() but returns a function type.
//
// Example:
//
//	http.HandleFunc("/version", version.HandlerFunc())
func HandlerFunc() http.HandlerFunc {
	h := Handler()
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

// HealthHandlerFunc is a convenience function that returns an http.HandlerFunc
// instead of http.Handler. It's equivalent to HealthHandler() but returns a function type.
//
// Example:
//
//	http.HandleFunc("/health", version.HealthHandlerFunc())
func HealthHandlerFunc() http.HandlerFunc {
	h := HealthHandler()
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

// Middleware returns an http middleware that adds version information to response headers.
// The middleware adds the following headers:
//   - X-App-Version: The application version
//   - X-Git-Commit: The git commit hash (if available)
//
// The middleware is non-blocking and will use default values if version info is unavailable.
//
// Thread-safe for concurrent use by multiple goroutines.
//
// Example:
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/api/users", handleUsers)
//
//	// Wrap the entire mux with version middleware
//	http.ListenAndServe(":8080", version.Middleware(mux))
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get version info, but don't block on failure
		if info, err := Get(); err == nil {
			w.Header().Set(HTTPHeaderAppVersion, info.Project.Version)
			if info.Git.Commit != DefaultGitCommit {
				w.Header().Set(HTTPHeaderGitCommit, info.Git.Commit)
			}
		}

		next.ServeHTTP(w, r)
	})
}
