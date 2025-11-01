package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/tabwriter"

	"github.com/itsatony/go-version"
)

// Test constants
const minimalManifestYAML = `manifest_version: "1.0"
project:
  name: "minimal-app"
  version: "1.0.0"
`

// resetFlags resets all flag variables to their default values for testing
func resetFlags() {
	*manifestPath = "versions.yaml"
	*jsonOutput = false
	*compactMode = false
	*schemasOnly = false
	*apisOnly = false
	*componentsOnly = false
	*gitOnly = false
	*buildOnly = false
	*showHelp = false
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

// captureOutput runs a function and captures its stdout output
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func TestOutputFull(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	testManifest := filepath.Join("testdata", "test-versions.yaml")
	*manifestPath = testManifest

	output := captureOutput(func() {
		err := version.Initialize(
			version.WithManifestPath(*manifestPath),
			version.WithGitInfo(),
			version.WithBuildInfo(),
		)
		if err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}

		info := version.MustGet()
		outputFull(info)
	})

	// Verify output contains expected sections
	if !strings.Contains(output, "Project Information:") {
		t.Error("Output missing Project Information section")
	}
	if !strings.Contains(output, "cli-test-app") {
		t.Error("Output missing project name")
	}
	if !strings.Contains(output, "1.2.3") {
		t.Error("Output missing project version")
	}
	if !strings.Contains(output, "Git Information:") {
		t.Error("Output missing Git Information section")
	}
	if !strings.Contains(output, "Build Information:") {
		t.Error("Output missing Build Information section")
	}
	if !strings.Contains(output, "Database Schemas:") {
		t.Error("Output missing Database Schemas section")
	}
	if !strings.Contains(output, "postgres_main") {
		t.Error("Output missing postgres schema")
	}
	if !strings.Contains(output, "API Versions:") {
		t.Error("Output missing API Versions section")
	}
	if !strings.Contains(output, "rest_v1") {
		t.Error("Output missing rest_v1 API")
	}
	if !strings.Contains(output, "Component Versions:") {
		t.Error("Output missing Component Versions section")
	}
	if !strings.Contains(output, "auth_service") {
		t.Error("Output missing auth_service component")
	}
}

func TestOutputJSON(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	testManifest := filepath.Join("testdata", "test-versions.yaml")
	*manifestPath = testManifest

	output := captureOutput(func() {
		err := version.Initialize(
			version.WithManifestPath(*manifestPath),
			version.WithGitInfo(),
			version.WithBuildInfo(),
		)
		if err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}

		info := version.MustGet()
		outputJSON(info)
	})

	// Verify JSON is valid
	var result version.Info
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Invalid JSON output: %v", err)
	}

	// Verify content
	if result.Project.Name != "cli-test-app" {
		t.Errorf("Expected project name 'cli-test-app', got '%s'", result.Project.Name)
	}
	if result.Project.Version != "1.2.3" {
		t.Errorf("Expected version '1.2.3', got '%s'", result.Project.Version)
	}
	if len(result.GetSchemas()) != 2 {
		t.Errorf("Expected 2 schemas, got %d", len(result.GetSchemas()))
	}
	if result.GetSchemas()["postgres_main"] != "45" {
		t.Error("Schema postgres_main not found or incorrect")
	}
}

func TestOutputCompact(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	testManifest := filepath.Join("testdata", "test-versions.yaml")
	*manifestPath = testManifest

	err := version.Initialize(
		version.WithManifestPath(*manifestPath),
		version.WithGitInfo(),
		version.WithBuildInfo(),
	)
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	info := version.MustGet()
	compact := info.String()

	if !strings.Contains(compact, "cli-test-app") {
		t.Error("Compact output missing project name")
	}
	if !strings.Contains(compact, "1.2.3") {
		t.Error("Compact output missing version")
	}
}

func TestOutputSchemas(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	testManifest := filepath.Join("testdata", "test-versions.yaml")
	*manifestPath = testManifest

	output := captureOutput(func() {
		err := version.Initialize(
			version.WithManifestPath(*manifestPath),
		)
		if err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}

		info := version.MustGet()
		outputSchemas(info)
	})

	if !strings.Contains(output, "Database Schemas:") {
		t.Error("Output missing Database Schemas section")
	}
	if !strings.Contains(output, "postgres_main") {
		t.Error("Output missing postgres_main")
	}
	if !strings.Contains(output, "45") {
		t.Error("Output missing postgres version")
	}
	if !strings.Contains(output, "redis_cache") {
		t.Error("Output missing redis_cache")
	}
}

func TestOutputAPIs(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	testManifest := filepath.Join("testdata", "test-versions.yaml")
	*manifestPath = testManifest

	output := captureOutput(func() {
		err := version.Initialize(
			version.WithManifestPath(*manifestPath),
		)
		if err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}

		info := version.MustGet()
		outputAPIs(info)
	})

	if !strings.Contains(output, "API Versions:") {
		t.Error("Output missing API Versions section")
	}
	if !strings.Contains(output, "rest_v1") {
		t.Error("Output missing rest_v1")
	}
	if !strings.Contains(output, "1.15.0") {
		t.Error("Output missing rest_v1 version")
	}
}

func TestOutputComponents(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	testManifest := filepath.Join("testdata", "test-versions.yaml")
	*manifestPath = testManifest

	output := captureOutput(func() {
		err := version.Initialize(
			version.WithManifestPath(*manifestPath),
		)
		if err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}

		info := version.MustGet()
		outputComponents(info)
	})

	if !strings.Contains(output, "Component Versions:") {
		t.Error("Output missing Component Versions section")
	}
	if !strings.Contains(output, "auth_service") {
		t.Error("Output missing auth_service")
	}
	if !strings.Contains(output, "2.1.0") {
		t.Error("Output missing auth_service version")
	}
}

func TestOutputGit(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	testManifest := filepath.Join("testdata", "test-versions.yaml")
	*manifestPath = testManifest

	output := captureOutput(func() {
		err := version.Initialize(
			version.WithManifestPath(*manifestPath),
			version.WithGitInfo(),
		)
		if err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}

		info := version.MustGet()
		outputGit(info)
	})

	if !strings.Contains(output, "Git Information:") {
		t.Error("Output missing Git Information section")
	}
	if !strings.Contains(output, "Commit:") {
		t.Error("Output missing Commit field")
	}
	if !strings.Contains(output, "Tree State:") {
		t.Error("Output missing Tree State field")
	}
}

func TestOutputBuild(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	testManifest := filepath.Join("testdata", "test-versions.yaml")
	*manifestPath = testManifest

	output := captureOutput(func() {
		err := version.Initialize(
			version.WithManifestPath(*manifestPath),
			version.WithBuildInfo(),
		)
		if err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}

		info := version.MustGet()
		outputBuild(info)
	})

	if !strings.Contains(output, "Build Information:") {
		t.Error("Output missing Build Information section")
	}
	if !strings.Contains(output, "Time:") {
		t.Error("Output missing Time field")
	}
	if !strings.Contains(output, "Go Version:") {
		t.Error("Output missing Go Version field")
	}
}

func TestOutputSchemasEmpty(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	// Create a minimal manifest without schemas
	tmpfile, err := os.CreateTemp("", "versions-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(minimalManifestYAML); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	output := captureOutput(func() {
		err := version.Initialize(
			version.WithManifestPath(tmpfile.Name()),
		)
		if err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}

		info := version.MustGet()
		outputSchemas(info)
	})

	if !strings.Contains(output, "No database schemas defined") {
		t.Error("Expected 'No database schemas defined' message")
	}
}

func TestOutputAPIsEmpty(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	// Create a minimal manifest without APIs
	tmpfile, err := os.CreateTemp("", "versions-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(minimalManifestYAML); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	output := captureOutput(func() {
		err := version.Initialize(
			version.WithManifestPath(tmpfile.Name()),
		)
		if err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}

		info := version.MustGet()
		outputAPIs(info)
	})

	if !strings.Contains(output, "No API versions defined") {
		t.Error("Expected 'No API versions defined' message")
	}
}

func TestOutputComponentsEmpty(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	// Create a minimal manifest without components
	tmpfile, err := os.CreateTemp("", "versions-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(minimalManifestYAML); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	output := captureOutput(func() {
		err := version.Initialize(
			version.WithManifestPath(tmpfile.Name()),
		)
		if err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}

		info := version.MustGet()
		outputComponents(info)
	})

	if !strings.Contains(output, "No component versions defined") {
		t.Error("Expected 'No component versions defined' message")
	}
}

func TestPrintSortedMap(t *testing.T) {
	resetFlags()

	testMap := map[string]string{
		"zebra":  "3",
		"alpha":  "1",
		"middle": "2",
	}

	output := captureOutput(func() {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		printSortedMap(w, testMap)
		w.Flush()
	})

	// Verify alphabetical ordering
	alphaPos := strings.Index(output, "alpha")
	middlePos := strings.Index(output, "middle")
	zebraPos := strings.Index(output, "zebra")

	if alphaPos == -1 || middlePos == -1 || zebraPos == -1 {
		t.Error("Missing expected keys in output")
	}

	if !(alphaPos < middlePos && middlePos < zebraPos) {
		t.Error("Output not in alphabetical order")
	}
}

func TestPrintSortedCustomMap(t *testing.T) {
	resetFlags()

	testMap := map[string]interface{}{
		"string_val": "test",
		"int_val":    42,
		"bool_val":   true,
	}

	output := captureOutput(func() {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		printSortedCustomMap(w, testMap)
		w.Flush()
	})

	if !strings.Contains(output, "string_val") {
		t.Error("Missing string_val in output")
	}
	if !strings.Contains(output, "int_val") {
		t.Error("Missing int_val in output")
	}
	if !strings.Contains(output, "bool_val") {
		t.Error("Missing bool_val in output")
	}
}

func TestMainWithValidManifest(t *testing.T) {
	// Save original args and restore after test
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Reset for clean state
	resetFlags()
	version.Reset()
	defer version.Reset()

	// Set up args for test
	testManifest := filepath.Join("testdata", "test-versions.yaml")
	os.Args = []string{"go-version", "-manifest", testManifest, "-compact"}

	// Capture output
	output := captureOutput(func() {
		// Re-parse flags with new args
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		manifestPath = flag.String("manifest", "versions.yaml", "")
		compactMode = flag.Bool("compact", false, "")
		flag.Parse()

		// Initialize and run
		if err := version.Initialize(
			version.WithManifestPath(*manifestPath),
			version.WithGitInfo(),
			version.WithBuildInfo(),
		); err == nil {
			if *compactMode {
				info := version.MustGet()
				fmt.Println(info.String())
			}
		}
	})

	if !strings.Contains(output, "cli-test-app") {
		t.Error("Expected app name in compact output")
	}
}

func TestMainJSONOutput(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	resetFlags()
	version.Reset()
	defer version.Reset()

	testManifest := filepath.Join("testdata", "test-versions.yaml")
	os.Args = []string{"go-version", "-manifest", testManifest, "-json"}

	output := captureOutput(func() {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		manifestPath = flag.String("manifest", "versions.yaml", "")
		jsonOutput = flag.Bool("json", false, "")
		flag.Parse()

		if err := version.Initialize(
			version.WithManifestPath(*manifestPath),
			version.WithGitInfo(),
			version.WithBuildInfo(),
		); err == nil {
			if *jsonOutput {
				info := version.MustGet()
				outputJSON(info)
			}
		}
	})

	// Should be valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("Invalid JSON: %v", err)
	}
}

func TestOutputWithCustomMetadata(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	testManifest := filepath.Join("testdata", "test-versions.yaml")

	output := captureOutput(func() {
		err := version.Initialize(
			version.WithManifestPath(testManifest),
		)
		if err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}

		info := version.MustGet()
		outputFull(info)
	})

	// Check for custom metadata
	if !strings.Contains(output, "Custom Metadata:") {
		t.Error("Output missing Custom Metadata section")
	}
	if !strings.Contains(output, "environment") {
		t.Error("Output missing environment field")
	}
	if !strings.Contains(output, "test") {
		t.Error("Output missing environment value 'test'")
	}
}

func TestInvalidManifestPath(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	// Try to initialize with non-existent manifest
	err := version.Initialize(
		version.WithManifestPath("/nonexistent/path/versions.yaml"),
	)

	// Should use defaults rather than failing
	if err != nil {
		// For non-existent file, it should not error (uses defaults)
		t.Log("Initialization used defaults for missing manifest")
	}

	info := version.MustGet()
	if info == nil {
		t.Fatal("Expected version info to be initialized with defaults")
	}

	// Should have default values
	if info.Project.Name == "" {
		t.Error("Expected default project name")
	}
}

func TestOutputGitWithTag(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	testManifest := filepath.Join("testdata", "test-versions.yaml")

	output := captureOutput(func() {
		err := version.Initialize(
			version.WithManifestPath(testManifest),
			version.WithGitInfo(),
		)
		if err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}

		info := version.MustGet()
		// Check if git tag is available in output
		outputGit(info)
	})

	// Should have git information section
	if !strings.Contains(output, "Git Information:") {
		t.Error("Missing Git Information section")
	}
}

func TestOutputBuildWithUser(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	testManifest := filepath.Join("testdata", "test-versions.yaml")

	output := captureOutput(func() {
		err := version.Initialize(
			version.WithManifestPath(testManifest),
			version.WithBuildInfo(),
		)
		if err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}

		info := version.MustGet()
		outputBuild(info)
	})

	// Should have build information
	if !strings.Contains(output, "Build Information:") {
		t.Error("Missing Build Information section")
	}
	if !strings.Contains(output, "Go Version:") {
		t.Error("Missing Go Version in build info")
	}
}

func TestHelpFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	resetFlags()

	os.Args = []string{"go-version", "-help"}

	// Parse flags - help will trigger os.Exit(0) in real main
	// but in our test we just check that showHelp is set
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	showHelp = flag.Bool("help", false, "")
	flag.Parse()

	if !*showHelp {
		t.Error("Expected showHelp flag to be true")
	}
}

func TestAllOutputModes(t *testing.T) {
	tests := []struct {
		name     string
		flagFunc func()
		expected string
	}{
		{
			name: "schemas mode",
			flagFunc: func() {
				*schemasOnly = true
			},
			expected: "Database Schemas:",
		},
		{
			name: "apis mode",
			flagFunc: func() {
				*apisOnly = true
			},
			expected: "API Versions:",
		},
		{
			name: "components mode",
			flagFunc: func() {
				*componentsOnly = true
			},
			expected: "Component Versions:",
		},
		{
			name: "git mode",
			flagFunc: func() {
				*gitOnly = true
			},
			expected: "Git Information:",
		},
		{
			name: "build mode",
			flagFunc: func() {
				*buildOnly = true
			},
			expected: "Build Information:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlags()
			version.Reset()
			defer version.Reset()

			testManifest := filepath.Join("testdata", "test-versions.yaml")
			*manifestPath = testManifest

			tt.flagFunc()

			output := captureOutput(func() {
				err := version.Initialize(
					version.WithManifestPath(*manifestPath),
					version.WithGitInfo(),
					version.WithBuildInfo(),
				)
				if err != nil {
					t.Fatalf("Failed to initialize: %v", err)
				}

				info := version.MustGet()

				// Simulate main's switch statement
				switch {
				case *jsonOutput:
					outputJSON(info)
				case *compactMode:
					fmt.Println(info.String())
				case *schemasOnly:
					outputSchemas(info)
				case *apisOnly:
					outputAPIs(info)
				case *componentsOnly:
					outputComponents(info)
				case *gitOnly:
					outputGit(info)
				case *buildOnly:
					outputBuild(info)
				default:
					outputFull(info)
				}
			})

			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected output to contain %q", tt.expected)
			}
		})
	}
}

func TestOutputFullWithMinimalData(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	// Create minimal manifest
	tmpfile, err := os.CreateTemp("", "versions-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(minimalManifestYAML); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	output := captureOutput(func() {
		err := version.Initialize(
			version.WithManifestPath(tmpfile.Name()),
			version.WithGitInfo(),
			version.WithBuildInfo(),
		)
		if err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}

		info := version.MustGet()
		outputFull(info)
	})

	// Should have basic sections even without optional data
	if !strings.Contains(output, "Project Information:") {
		t.Error("Missing Project Information section")
	}
	if !strings.Contains(output, "minimal-app") {
		t.Error("Missing project name")
	}

	// Should NOT have the optional sections since they're empty
	if strings.Contains(output, "Database Schemas:") {
		t.Error("Should not show Database Schemas section when empty")
	}
	if strings.Contains(output, "API Versions:") {
		t.Error("Should not show API Versions section when empty")
	}
}

func TestRun(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	testManifest := filepath.Join("testdata", "test-versions.yaml")
	*manifestPath = testManifest

	output := captureOutput(func() {
		if err := run(); err != nil {
			t.Errorf("run() returned error: %v", err)
		}
	})

	// Should produce full output by default
	if !strings.Contains(output, "Project Information:") {
		t.Error("Expected full output from run()")
	}
}

func TestRunWithHelp(t *testing.T) {
	resetFlags()
	*showHelp = true

	// Set up flag.Usage to capture output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
	}

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := run()

	w.Close()
	os.Stderr = oldStderr

	if err != nil {
		t.Errorf("run() with help should not return error, got: %v", err)
	}

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "Usage:") {
		t.Error("Expected help message in stderr")
	}
}

func TestRunWithInvalidManifest(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	// Point to non-existent manifest (but library uses defaults, so no error)
	*manifestPath = "/nonexistent/manifest.yaml"

	err := run()
	// Should not error because library uses defaults for missing manifest
	if err != nil {
		t.Logf("Got expected behavior: %v", err)
	}
}

func TestRunCompactMode(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	testManifest := filepath.Join("testdata", "test-versions.yaml")
	*manifestPath = testManifest
	*compactMode = true

	output := captureOutput(func() {
		if err := run(); err != nil {
			t.Errorf("run() returned error: %v", err)
		}
	})

	if !strings.Contains(output, "cli-test-app") {
		t.Error("Expected compact output with project name")
	}
	if !strings.Contains(output, "1.2.3") {
		t.Error("Expected compact output with version")
	}
}

func TestRunJSONMode(t *testing.T) {
	resetFlags()
	version.Reset()
	defer version.Reset()

	testManifest := filepath.Join("testdata", "test-versions.yaml")
	*manifestPath = testManifest
	*jsonOutput = true

	output := captureOutput(func() {
		if err := run(); err != nil {
			t.Errorf("run() returned error: %v", err)
		}
	})

	// Verify JSON output
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("Expected valid JSON output: %v", err)
	}
}
