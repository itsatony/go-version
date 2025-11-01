package main

import (
	"fmt"
	"log"

	"github.com/itsatony/go-version"
)

// Simple example demonstrating zero-config usage of go-version.
// This example shows how to use the library with minimal setup.
//
// Run with:
//
//	go run main.go
func main() {
	fmt.Println("=== Simple go-version Example ===")
	fmt.Println()

	// Zero-config usage: just call MustGet()
	// If versions.yaml exists, it will be loaded automatically
	// Otherwise, defaults will be used
	info := version.MustGet()

	// Display basic version information
	fmt.Printf("Application: %s\n", info.Project.Name)
	fmt.Printf("Version:     %s\n", info.Project.Version)
	fmt.Printf("Git Commit:  %s\n", info.Git.Commit)
	fmt.Printf("Git Tag:     %s\n", info.Git.Tag)
	fmt.Printf("Tree State:  %s\n", info.Git.TreeState)
	fmt.Printf("Build Time:  %s\n", info.Build.Time)
	fmt.Printf("Go Version:  %s\n", info.Build.GoVersion)

	// Access schema versions
	if len(info.GetSchemas()) > 0 {
		fmt.Println("\nDatabase Schemas:")
		for name, ver := range info.GetSchemas() {
			fmt.Printf("  %s: %s\n", name, ver)
		}
	}

	// Access API versions
	if len(info.GetAPIs()) > 0 {
		fmt.Println("\nAPI Versions:")
		for name, ver := range info.GetAPIs() {
			fmt.Printf("  %s: %s\n", name, ver)
		}
	}

	// Access component versions
	if len(info.GetComponents()) > 0 {
		fmt.Println("\nComponent Versions:")
		for name, ver := range info.GetComponents() {
			fmt.Printf("  %s: %s\n", name, ver)
		}
	}

	// Use the String() method for a compact representation
	fmt.Printf("\nCompact Format: %s\n", info.String())

	// Example: Check if a specific schema version is available
	if schemaVer, ok := info.GetSchemaVersion("postgres_main"); ok {
		fmt.Printf("\nPostgreSQL main schema version: %s\n", schemaVer)
	} else {
		fmt.Println("\nNo PostgreSQL main schema version defined")
	}

	// Example: Using with structured logging (zap fields)
	// In a real application, you would use these fields with zap logger:
	// logger := zap.NewProduction()
	// logger = logger.With(info.LogFields()...)
	// logger.Info("Application started")
	fmt.Println("\nLog fields for structured logging:")
	fields := info.LogFields()
	for _, field := range fields {
		fmt.Printf("  %s: %v\n", field.Key, field.Interface)
	}
}

// Example output:
//
// === Simple go-version Example ===
//
// Application: unknown
// Version:     0.0.0-dev
// Git Commit:  dev
// Git Tag:
// Tree State:  unknown
// Build Time:  unknown
// Go Version:  go1.24.0
//
// Compact Format: unknown 0.0.0-dev (git: dev)
//
// No PostgreSQL main schema version defined
//
// Log fields for structured logging:
//   version: 0.0.0-dev
//   git_commit: dev
//   build_time: unknown

// Example with versions.yaml:
//
// If you create a versions.yaml file in this directory:
//
// manifest_version: "1.0"
// project:
//   name: "my-app"
//   version: "1.2.3"
// schemas:
//   postgres_main: "45"
//   redis_cache: "3"
// apis:
//   rest_v1: "1.15.0"
//
// The output will show:
//
// Application: my-app
// Version:     1.2.3
// ...
// Database Schemas:
//   postgres_main: 45
//   redis_cache: 3
//
// API Versions:
//   rest_v1: 1.15.0

// Example with graceful error handling:
func exampleWithErrorHandling() {
	// Use Get() instead of MustGet() for graceful error handling
	info, err := version.Get()
	if err != nil {
		log.Printf("Warning: Version info unavailable: %v", err)
		log.Println("Continuing with limited functionality...")
		return
	}

	log.Printf("Application %s v%s started", info.Project.Name, info.Project.Version)
}
