package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/itsatony/go-version"
)

// Advanced example demonstrating comprehensive usage of go-version.
// This example shows:
// - Custom initialization with embedded manifest
// - Version validation
// - HTTP endpoints for version info
// - Middleware integration
// - Graceful shutdown
//
// Run with:
//
//	go run main.go
func main() {
	fmt.Println("=== Advanced go-version Example ===")
	fmt.Println()

	// Initialize with custom configuration
	if err := initializeVersion(); err != nil {
		log.Fatalf("Failed to initialize version: %v", err)
	}

	// Display version information
	displayVersionInfo()

	// Start HTTP server with version endpoints
	srv := startHTTPServer()

	// Wait for shutdown signal
	waitForShutdown(srv)
}

// initializeVersion demonstrates advanced initialization with:
// - Embedded manifest data
// - Custom validators
// - Git and build info enrichment
func initializeVersion() error {
	// In a real application, you might use go:embed to embed versions.yaml
	// For this example, we'll use inline YAML
	manifestData := []byte(`
manifest_version: "1.0"
project:
  name: "advanced-example-app"
  version: "2.5.1"
Schemas:
  postgres_main: "50"
  postgres_analytics: "15"
  redis_cache: "5"
APIs:
  rest_v1: "1.15.0"
  rest_v2: "2.1.0"
  grpc: "1.3.0"
Components:
  aigentchat: "3.4.1"
  notification_service: "2.1.0"
Custom:
  environment: "production"
  region: "us-east-1"
  deployment_id: "dep-20250111-001"
`)

	// Initialize with:
	// 1. Embedded manifest data
	// 2. Git information
	// 3. Build information
	// 4. Validators to enforce minimum versions
	err := version.Initialize(
		version.WithEmbedded(manifestData),
		version.WithGitInfo(),
		version.WithBuildInfo(),
		version.WithValidators(
			// Enforce minimum schema versions
			version.NewSchemaValidator("postgres_main", "45"),
			version.NewSchemaValidator("redis_cache", "3"),
			// Enforce minimum API versions
			version.NewAPIValidator("rest_v1", "1.10.0"),
			// Enforce minimum component versions
			version.NewComponentValidator("aigentchat", "3.0.0"),
			// Custom validator using ValidatorFunc
			version.ValidatorFunc(func(ctx context.Context, info *version.Info) error {
				// Example: Ensure production builds have a valid git tag
				if env, ok := info.GetCustom()["environment"].(string); ok && env == "production" {
					if info.Git.Tag == "" {
						return fmt.Errorf("production builds must have a git tag")
					}
				}
				return nil
			}),
		),
	)

	if err != nil {
		return fmt.Errorf("version initialization failed: %w", err)
	}

	return nil
}

// displayVersionInfo shows comprehensive version information
func displayVersionInfo() {
	info := version.MustGet()

	fmt.Println("Project Information:")
	fmt.Printf("  Name:    %s\n", info.Project.Name)
	fmt.Printf("  Version: %s\n", info.Project.Version)

	fmt.Println("\nGit Information:")
	fmt.Printf("  Commit:     %s\n", info.Git.Commit)
	fmt.Printf("  Tag:        %s\n", info.Git.Tag)
	fmt.Printf("  Tree State: %s\n", info.Git.TreeState)
	if info.Git.CommitTime != "" {
		fmt.Printf("  Commit Time: %s\n", info.Git.CommitTime)
	}

	fmt.Println("\nBuild Information:")
	fmt.Printf("  Time:       %s\n", info.Build.Time)
	fmt.Printf("  User:       %s\n", info.Build.User)
	fmt.Printf("  Go Version: %s\n", info.Build.GoVersion)

	if len(info.GetSchemas()) > 0 {
		fmt.Println("\nDatabase Schemas:")
		for name, ver := range info.GetSchemas() {
			fmt.Printf("  %-25s %s\n", name+":", ver)
		}
	}

	if len(info.GetAPIs()) > 0 {
		fmt.Println("\nAPI Versions:")
		for name, ver := range info.GetAPIs() {
			fmt.Printf("  %-25s %s\n", name+":", ver)
		}
	}

	if len(info.GetComponents()) > 0 {
		fmt.Println("\nComponent Versions:")
		for name, ver := range info.GetComponents() {
			fmt.Printf("  %-25s %s\n", name+":", ver)
		}
	}

	if len(info.GetCustom()) > 0 {
		fmt.Println("\nCustom Metadata:")
		for key, val := range info.GetCustom() {
			fmt.Printf("  %-25s %v\n", key+":", val)
		}
	}

	fmt.Printf("\nLoaded at: %s\n", info.LoadedAt().Format(time.RFC3339))
	fmt.Println()
}

// startHTTPServer sets up an HTTP server with version endpoints
func startHTTPServer() *http.Server {
	mux := http.NewServeMux()

	// Version endpoint - returns full version info as JSON
	mux.Handle("/version", version.Handler())

	// Health check endpoint
	mux.Handle("/health", version.HealthHandler())

	// Example API endpoint
	mux.HandleFunc("/api/users", handleUsers)

	// Example admin endpoint
	mux.HandleFunc("/admin/status", handleStatus)

	// Wrap the entire mux with version middleware
	// This adds X-App-Version and X-Git-Commit headers to all responses
	handler := version.Middleware(mux)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		fmt.Printf("Starting HTTP server on %s\n", srv.Addr)
		fmt.Println("\nAvailable endpoints:")
		fmt.Println("  GET  /version       - Full version information (JSON)")
		fmt.Println("  GET  /health        - Health check")
		fmt.Println("  GET  /api/users     - Example API endpoint")
		fmt.Println("  GET  /admin/status  - Example admin endpoint")
		fmt.Println("\nTry:")
		fmt.Println("  curl http://localhost:8080/version")
		fmt.Println("  curl http://localhost:8080/health")
		fmt.Println("  curl -I http://localhost:8080/api/users")
		fmt.Println()
		fmt.Println("Press Ctrl+C to shutdown")
		fmt.Println()

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	return srv
}

// handleUsers is an example API endpoint
func handleUsers(w http.ResponseWriter, r *http.Request) {
	// The version middleware automatically adds X-App-Version and X-Git-Commit headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"users": [], "message": "Check response headers for version info"}`))
}

// handleStatus is an example admin endpoint showing version info
func handleStatus(w http.ResponseWriter, r *http.Request) {
	info := version.MustGet()

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Application Status\n")
	fmt.Fprintf(w, "==================\n\n")
	fmt.Fprintf(w, "Application: %s v%s\n", info.Project.Name, info.Project.Version)
	fmt.Fprintf(w, "Git Commit:  %s\n", info.Git.Commit)
	fmt.Fprintf(w, "Build Time:  %s\n", info.Build.Time)
	fmt.Fprintf(w, "\nUptime: %s\n", time.Since(info.LoadedAt()).Round(time.Second))

	if len(info.GetSchemas()) > 0 {
		fmt.Fprintf(w, "\nDatabase Schemas:\n")
		for name, ver := range info.GetSchemas() {
			fmt.Fprintf(w, "  %s: %s\n", name, ver)
		}
	}
}

// waitForShutdown waits for interrupt signal and gracefully shuts down the server
func waitForShutdown(srv *http.Server) {
	// Create channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until signal received
	<-quit
	fmt.Println("\nShutdown signal received, stopping server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server stopped")
}

// Example demonstrating version comparison in application logic
func exampleVersionComparison() {
	info := version.MustGet()

	// Check if API version meets minimum requirement
	if apiVer, ok := info.GetAPIVersion("rest_v1"); ok {
		fmt.Printf("REST API v1 is at version %s\n", apiVer)

		// You could use internal/semver package for comparison
		// if you need more complex version logic
	}

	// Check database schema version
	if schemaVer, ok := info.GetSchemaVersion("postgres_main"); ok {
		fmt.Printf("PostgreSQL schema is at version %s\n", schemaVer)
	}
}

// Example demonstrating non-singleton usage
func exampleNonSingleton() {
	// Create independent version info for a different service
	otherServiceManifest := []byte(`
manifest_version: "1.0"
project:
  name: "other-service"
  version: "1.0.0"
`)

	info, err := version.New(version.WithEmbedded(otherServiceManifest))
	if err != nil {
		log.Printf("Failed to load other service version: %v", err)
		return
	}

	fmt.Printf("Other service: %s v%s\n", info.Project.Name, info.Project.Version)
}

// Example: Building with ldflags injection
//
// You can inject version information at build time:
//
// go build -ldflags="-X github.com/itsatony/go-version.GitCommit=$(git rev-parse HEAD) \
//   -X github.com/itsatony/go-version.GitTag=$(git describe --tags --always) \
//   -X github.com/itsatony/go-version.BuildTime=$(date -u '+%Y-%m-%dT%H:%M:%SZ') \
//   -X github.com/itsatony/go-version.BuildUser=$(whoami)"
//
// Or use a Makefile:
//
// VERSION := $(shell git describe --tags --always --dirty)
// COMMIT := $(shell git rev-parse HEAD)
// BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
// BUILD_USER := $(shell whoami)
//
// LDFLAGS := -X github.com/itsatony/go-version.GitCommit=$(COMMIT)
// LDFLAGS += -X github.com/itsatony/go-version.GitTag=$(VERSION)
// LDFLAGS += -X github.com/itsatony/go-version.BuildTime=$(BUILD_TIME)
// LDFLAGS += -X github.com/itsatony/go-version.BuildUser=$(BUILD_USER)
//
// build:
//     go build -ldflags="$(LDFLAGS)" -o myapp ./cmd/myapp
