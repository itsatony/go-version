package version

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Variables for -ldflags injection at build time.
// These can be set using:
//
//	go build -ldflags="-X github.com/itsatony/go-version.GitCommit=abc123"
//
// Example Makefile:
//
//	LDFLAGS := -X github.com/itsatony/go-version.GitCommit=$(shell git rev-parse HEAD)
//	LDFLAGS += -X github.com/itsatony/go-version.GitTag=$(shell git describe --tags --always)
//	LDFLAGS += -X github.com/itsatony/go-version.BuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
//	LDFLAGS += -X github.com/itsatony/go-version.BuildUser=$(shell whoami)
var (
	// GitCommit is the git commit hash, injected via ldflags
	GitCommit = DefaultGitCommit

	// GitTag is the git tag, injected via ldflags
	GitTag = ""

	// GitTreeState is the git tree state (clean/dirty), injected via ldflags
	GitTreeState = DefaultGitTreeState

	// BuildTime is the build timestamp, injected via ldflags
	BuildTime = DefaultBuildTime

	// BuildUser is the user who built the binary, injected via ldflags
	BuildUser = DefaultBuildUser
)

// gitCommandTimeout is the maximum time allowed for git command execution.
// This prevents hanging if git is unresponsive.
const gitCommandTimeout = 5 * time.Second

// loadVersionInfo is the main entry point for loading version information.
// It applies options, loads the manifest, enriches with git/build info, and validates.
//
// Load precedence:
//  1. Embedded manifest (if provided via WithEmbedded)
//  2. File manifest (from WithManifestPath or default "versions.yaml")
//  3. Defaults (if no manifest found)
//
// Then enriches with:
//   - Git info (if WithGitInfo, default true)
//   - Build info (if WithBuildInfo, default true)
//
// Finally runs validators (if any).
func loadVersionInfo(opts ...Option) (*Info, error) {
	// Apply options
	options := defaultLoadOptions()
	for _, opt := range opts {
		opt(options)
	}

	// Load manifest
	manifest, err := loadManifest(options)
	if err != nil {
		// In strict mode, all manifest errors are fatal
		if options.strictMode {
			if os.IsNotExist(err) {
				return nil, newCategoryErrorWithHint(CategoryManifest, ErrMsgStrictModeManifestRequired, ErrHintStrictMode)
			}
			return nil, wrapErrorWithHint(err, CategoryManifest, ErrMsgLoadManifest, ErrHintParseYAML)
		}

		// Non-strict mode: use defaults if manifest not found
		if !os.IsNotExist(err) {
			return nil, wrapErrorWithHint(err, CategoryManifest, ErrMsgLoadManifest, ErrHintParseYAML)
		}
		// Use default manifest
		manifest = defaultManifest()
	}

	// Convert manifest to Info
	info := manifestToInfo(manifest)

	// Set loaded time immediately to ensure immutability
	// This timestamp marks when version info loading began
	info.loadedAt = time.Now()

	// Enrich with git info
	if options.includeGit {
		enrichWithGitInfo(info)
	}

	// Enrich with build info
	if options.includeBuild {
		enrichWithBuildInfo(info)
	}

	// Run validators with provided context
	ctx := options.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	for _, validator := range options.validators {
		if err := validator.Validate(ctx, info); err != nil {
			return nil, wrapError(err, CategoryValidation, ErrMsgValidationFailedWrap)
		}
	}

	return info, nil
}

// loadManifest loads the manifest from embedded data or file.
// Precedence: embedded > file
func loadManifest(options *LoadOptions) (*Manifest, error) {
	// Try embedded first
	if len(options.manifestEmbed) > 0 {
		return parseManifest(options.manifestEmbed)
	}

	// Try file
	if options.manifestPath != "" {
		data, err := os.ReadFile(options.manifestPath)
		if err != nil {
			return nil, err
		}
		return parseManifest(data)
	}

	return nil, os.ErrNotExist
}

// parseManifest parses YAML manifest data.
func parseManifest(data []byte) (*Manifest, error) {
	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, wrapError(err, CategoryManifest, ErrMsgParseYAML)
	}

	// Validate manifest version
	if manifest.ManifestVersion == "" {
		manifest.ManifestVersion = ManifestVersion
	}

	// Validate project section
	if manifest.Project.Name == "" {
		return nil, newCategoryErrorWithHint(CategoryManifest, ErrMsgProjectNameRequired, ErrHintProjectNameRequired)
	}
	if manifest.Project.Version == "" {
		return nil, newCategoryErrorWithHint(CategoryManifest, ErrMsgProjectVersionRequired, ErrHintProjectVersionRequired)
	}

	return &manifest, nil
}

// defaultManifest returns a default manifest when no file is found.
func defaultManifest() *Manifest {
	return &Manifest{
		ManifestVersion: ManifestVersion,
		Project: ProjectManifest{
			Name:    DefaultProjectName,
			Version: DefaultProjectVersion,
		},
	}
}

// manifestToInfo converts a Manifest to an Info struct.
func manifestToInfo(m *Manifest) *Info {
	info := &Info{
		Project: ProjectVersion{
			Name:    m.Project.Name,
			Version: m.Project.Version,
		},
		Git: GitInfo{
			Commit:    DefaultGitCommit,
			TreeState: DefaultGitTreeState,
		},
		Build: BuildInfo{
			Time:      DefaultBuildTime,
			GoVersion: runtime.Version(),
		},
	}

	// Copy maps to unexported fields (defensive copies for immutability)
	if m.Schemas != nil {
		info.schemas = make(map[string]string, len(m.Schemas))
		for k, v := range m.Schemas {
			info.schemas[k] = v
		}
	}

	if m.APIs != nil {
		info.apis = make(map[string]string, len(m.APIs))
		for k, v := range m.APIs {
			info.apis[k] = v
		}
	}

	if m.Components != nil {
		info.components = make(map[string]string, len(m.Components))
		for k, v := range m.Components {
			info.components[k] = v
		}
	}

	if m.Custom != nil {
		info.custom = make(map[string]interface{}, len(m.Custom))
		for k, v := range m.Custom {
			info.custom[k] = v
		}
	}

	return info
}

// enrichWithGitInfo enriches Info with git metadata.
// Tries multiple sources in order:
//  1. Injected ldflags variables
//  2. runtime/debug.BuildInfo (Go 1.18+)
//  3. Direct git command execution (fallback)
func enrichWithGitInfo(info *Info) {
	// Use injected ldflags if available
	applyLdflagsGitInfo(info)

	// Try runtime/debug.BuildInfo
	applyBuildInfoGitData(info)

	// Fallback: try git command (if in git repo)
	applyGitCommandFallback(info)
}

// applyLdflagsGitInfo applies git information from ldflags injection
func applyLdflagsGitInfo(info *Info) {
	if GitCommit != DefaultGitCommit {
		info.Git.Commit = GitCommit
	}
	if GitTag != "" {
		info.Git.Tag = GitTag
	}
	if GitTreeState != DefaultGitTreeState {
		info.Git.TreeState = GitTreeState
	}
}

// applyBuildInfoGitData extracts git information from runtime/debug.BuildInfo
func applyBuildInfoGitData(info *Info) {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case VCSKeyRevision:
			if info.Git.Commit == DefaultGitCommit {
				info.Git.Commit = setting.Value
			}
		case VCSKeyTime:
			if info.Git.CommitTime == "" {
				info.Git.CommitTime = setting.Value
			}
		case VCSKeyModified:
			if setting.Value == VCSValueTrue {
				info.Git.TreeState = GitTreeStateDirty
			}
		}
	}
}

// applyGitCommandFallback tries to get git information via command execution
func applyGitCommandFallback(info *Info) {
	if info.Git.Commit == DefaultGitCommit {
		if commit := getGitCommit(); commit != "" {
			info.Git.Commit = commit
		}
	}
	if info.Git.Tag == "" {
		if tag := getGitTag(); tag != "" {
			info.Git.Tag = tag
		}
	}
}

// enrichWithBuildInfo enriches Info with build metadata.
func enrichWithBuildInfo(info *Info) {
	// Use injected ldflags
	if BuildTime != DefaultBuildTime {
		info.Build.Time = BuildTime
	}
	if BuildUser != DefaultBuildUser {
		info.Build.User = BuildUser
	}

	// Go version is always available
	info.Build.GoVersion = runtime.Version()

	// Try runtime/debug.BuildInfo for build time if not injected
	if info.Build.Time == DefaultBuildTime {
		if _, ok := debug.ReadBuildInfo(); ok {
			// BuildInfo doesn't have build time, but we can use current time as fallback
			// In practice, build time should be injected via ldflags
			_ = ok // Suppress unused variable warning
		}
	}
}

// getGitBinary returns the git binary path after security validation.
//
// SECURITY: Validates the git binary location to prevent PATH hijacking attacks.
// Only accepts git from standard system locations on Unix-like systems.
// On Windows, relies on exec.LookPath which uses system PATH.
//
// Returns empty string if:
//   - Git is not found
//   - Git is found in an unsafe location (Unix only)
//   - exec.LookPath fails
//
// Safe locations on Unix:
//   - /usr/bin/git
//   - /usr/local/bin/git
//   - /opt/homebrew/bin/git (macOS Homebrew)
func getGitBinary() string {
	// Look up git in PATH
	gitPath, err := exec.LookPath(GitCmdName)
	if err != nil {
		return ""
	}

	// On Unix systems, validate the git binary location
	// to prevent PATH hijacking attacks
	if runtime.GOOS != "windows" {
		// Whitelist of safe git locations
		safeLocations := []string{
			"/usr/bin/git",
			"/usr/local/bin/git",
			"/opt/homebrew/bin/git", // macOS Homebrew
		}

		// Check if git is in a safe location
		isSafe := false
		for _, safeLoc := range safeLocations {
			if gitPath == safeLoc {
				isSafe = true
				break
			}
		}

		if !isSafe {
			// Git found in unexpected location - potential security risk
			// Return empty to skip git info enrichment
			return ""
		}
	}

	return gitPath
}

// isValidCommitHash validates that a string looks like a valid git commit hash.
//
// SECURITY: Prevents injection of malicious data into git commit field.
// Valid formats:
//   - Full SHA-1: 40 hexadecimal characters (e.g., "abc123...")
//   - Short SHA-1: 7-40 hexadecimal characters
//
// Returns true if the hash is valid, false otherwise.
func isValidCommitHash(hash string) bool {
	// Commit hash must be between 7 and 40 characters
	if len(hash) < 7 || len(hash) > 40 {
		return false
	}

	// Must contain only hexadecimal characters
	for _, c := range hash {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}

	return true
}

// getGitCommit tries to get the current git commit hash.
// Returns empty string if git is not available or not in a git repo.
//
// SECURITY:
//   - Uses getGitBinary() to validate git binary location
//   - Uses CommandContext with timeout to prevent hangs
//   - Git command arguments are constants to prevent command injection
//   - Validates output format with isValidCommitHash()
func getGitCommit() string {
	// Validate git binary location first
	gitBinary := getGitBinary()
	if gitBinary == "" {
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), gitCommandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, gitBinary, GitCmdRevParse, GitArgHead)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	hash := strings.TrimSpace(string(output))

	// Validate the hash format
	if !isValidCommitHash(hash) {
		return ""
	}

	return hash
}

// getGitTag tries to get the current git tag.
// Returns empty string if no tag or git is not available.
//
// SECURITY:
//   - Uses getGitBinary() to validate git binary location
//   - Uses CommandContext with timeout to prevent hangs
//   - Git command arguments are constants to prevent command injection
func getGitTag() string {
	// Validate git binary location first
	gitBinary := getGitBinary()
	if gitBinary == "" {
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), gitCommandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, gitBinary, GitCmdDescribe, GitArgTags, GitArgExactMatch)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
