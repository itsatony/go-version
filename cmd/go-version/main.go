package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/itsatony/go-version"
)

const (
	appName = "go-version"
	usage   = `go-version - Display version information

Usage:
  go-version [options]

Options:
  -manifest string
        Path to versions.yaml manifest file (default: versions.yaml)
  -json
        Output in JSON format
  -compact
        Show compact single-line format
  -schemas
        Show only database schema versions
  -apis
        Show only API versions
  -components
        Show only component versions
  -git
        Show only git information
  -build
        Show only build information
  -help
        Show this help message

Examples:
  # Show all version information
  go-version

  # Show version in JSON format
  go-version -json

  # Use custom manifest file
  go-version -manifest ./config/versions.yaml

  # Show only schema versions
  go-version -schemas

  # Show compact format
  go-version -compact

  # Combine JSON with custom manifest
  go-version -json -manifest ./versions.yaml
`
)

var (
	manifestPath   = flag.String("manifest", "versions.yaml", "Path to versions.yaml manifest file")
	jsonOutput     = flag.Bool("json", false, "Output in JSON format")
	compactMode    = flag.Bool("compact", false, "Show compact single-line format")
	schemasOnly    = flag.Bool("schemas", false, "Show only database schema versions")
	apisOnly       = flag.Bool("apis", false, "Show only API versions")
	componentsOnly = flag.Bool("components", false, "Show only component versions")
	gitOnly        = flag.Bool("git", false, "Show only git information")
	buildOnly      = flag.Bool("build", false, "Show only build information")
	showHelp       = flag.Bool("help", false, "Show help message")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
	}
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// run executes the main CLI logic and returns any error.
// This function is separated from main() to make it testable.
func run() error {
	if *showHelp {
		flag.Usage()
		return nil
	}

	// Initialize version library
	err := version.Initialize(
		version.WithManifestPath(*manifestPath),
		version.WithGitInfo(),
		version.WithBuildInfo(),
	)
	if err != nil {
		return fmt.Errorf(
			"failed to load version information: %w\nMake sure %s exists or use -manifest to specify a different file",
			err,
			*manifestPath,
		)
	}

	info := version.MustGet()

	// Handle different output modes
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

	return nil
}

// outputJSON outputs version info as JSON
func outputJSON(info *version.Info) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(info); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

// outputFull displays complete version information in human-readable format
func outputFull(info *version.Info) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "Project Information:\n")
	fmt.Fprintf(w, "  Name:\t%s\n", info.Project.Name)
	fmt.Fprintf(w, "  Version:\t%s\n", info.Project.Version)
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "Git Information:\n")
	fmt.Fprintf(w, "  Commit:\t%s\n", info.Git.Commit)
	if info.Git.Tag != "" {
		fmt.Fprintf(w, "  Tag:\t%s\n", info.Git.Tag)
	}
	fmt.Fprintf(w, "  Tree State:\t%s\n", info.Git.TreeState)
	if info.Git.CommitTime != "" {
		fmt.Fprintf(w, "  Commit Time:\t%s\n", info.Git.CommitTime)
	}
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "Build Information:\n")
	fmt.Fprintf(w, "  Time:\t%s\n", info.Build.Time)
	if info.Build.User != "" {
		fmt.Fprintf(w, "  User:\t%s\n", info.Build.User)
	}
	fmt.Fprintf(w, "  Go Version:\t%s\n", info.Build.GoVersion)
	fmt.Fprintf(w, "\n")

	if len(info.GetSchemas()) > 0 {
		fmt.Fprintf(w, "Database Schemas:\n")
		printSortedMap(w, info.GetSchemas())
		fmt.Fprintf(w, "\n")
	}

	if len(info.GetAPIs()) > 0 {
		fmt.Fprintf(w, "API Versions:\n")
		printSortedMap(w, info.GetAPIs())
		fmt.Fprintf(w, "\n")
	}

	if len(info.GetComponents()) > 0 {
		fmt.Fprintf(w, "Component Versions:\n")
		printSortedMap(w, info.GetComponents())
		fmt.Fprintf(w, "\n")
	}

	if len(info.GetCustom()) > 0 {
		fmt.Fprintf(w, "Custom Metadata:\n")
		printSortedCustomMap(w, info.GetCustom())
	}

	w.Flush()
}

// outputSchemas displays only database schema versions
func outputSchemas(info *version.Info) {
	if len(info.GetSchemas()) == 0 {
		fmt.Println("No database schemas defined")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Database Schemas:\n")
	printSortedMap(w, info.GetSchemas())
	w.Flush()
}

// outputAPIs displays only API versions
func outputAPIs(info *version.Info) {
	if len(info.GetAPIs()) == 0 {
		fmt.Println("No API versions defined")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "API Versions:\n")
	printSortedMap(w, info.GetAPIs())
	w.Flush()
}

// outputComponents displays only component versions
func outputComponents(info *version.Info) {
	if len(info.GetComponents()) == 0 {
		fmt.Println("No component versions defined")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Component Versions:\n")
	printSortedMap(w, info.GetComponents())
	w.Flush()
}

// outputGit displays only git information
func outputGit(info *version.Info) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Git Information:\n")
	fmt.Fprintf(w, "  Commit:\t%s\n", info.Git.Commit)
	if info.Git.Tag != "" {
		fmt.Fprintf(w, "  Tag:\t%s\n", info.Git.Tag)
	}
	fmt.Fprintf(w, "  Tree State:\t%s\n", info.Git.TreeState)
	if info.Git.CommitTime != "" {
		fmt.Fprintf(w, "  Commit Time:\t%s\n", info.Git.CommitTime)
	}
	w.Flush()
}

// outputBuild displays only build information
func outputBuild(info *version.Info) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Build Information:\n")
	fmt.Fprintf(w, "  Time:\t%s\n", info.Build.Time)
	if info.Build.User != "" {
		fmt.Fprintf(w, "  User:\t%s\n", info.Build.User)
	}
	fmt.Fprintf(w, "  Go Version:\t%s\n", info.Build.GoVersion)
	w.Flush()
}

// printSortedMap prints a map in sorted order by keys
func printSortedMap(w *tabwriter.Writer, m map[string]string) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(w, "  %s:\t%s\n", k, m[k])
	}
}

// printSortedCustomMap prints a custom metadata map in sorted order
func printSortedCustomMap(w *tabwriter.Writer, m map[string]interface{}) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(w, "  %s:\t%v\n", k, m[k])
	}
}
