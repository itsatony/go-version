# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 2026-09-22

### ⚠️ BREAKING BEHAVIOUR CHANGE — read this before upgrading

**`DefaultGitTreeState` changed from `"clean"` to `"unknown"`.** Every binary that
does not inject a tree state via ldflags and carries no VCS build info now serves
`{"git":{"tree_state":"unknown"}}` on `/version` where it previously served
`"clean"`.

**If you have a test asserting `tree_state == "clean"` on an un-stamped build, it
will go red — and that redness is correct.** The old value was not a modest
fallback, it was a false statement: a binary with `"commit":"dev"` and
`"tree_state":"clean"` was telling an operator something it could not know.
`clean` and `unknown` are two different facts and must not share a field.

**Why this is a MINOR and not a MAJOR bump.** The Go API surface is unchanged and
purely additive — no identifier is removed, renamed or retyped, every existing
`switch`/comparison still compiles, and `GitTreeStateUnknown` is new. What changed
is the *value* of a documented constant and therefore the JSON a service emits.
A major bump on this module would require a `/v2` module path, forcing an import
rewrite in ~25 fleet repositories for a one-constant change, and without the path
change a `v2.0.0` tag would simply never be resolved by `go get` — the worse of
the two failure modes. The compatibility break is real but is in the *output*
contract, so it is announced loudly here instead.

**Remedy for anyone who wants `clean` back:** stamp it truthfully at build time —
`-X github.com/itsatony/go-version.GitTreeState=clean` now actually takes effect
(see below). Container images built from a source tree without `.git/` have no
other channel.

### Added
- `GitTreeStateUnknown = "unknown"` — the third, distinct tree state, meaning the
  build could not determine whether the working tree was clean or dirty.
- `VCSValueFalse` — the `"false"` reading of `runtime/debug.BuildInfo`'s
  `vcs.modified` setting, which is positive evidence of a *clean* tree.
- `git status --porcelain` fallback (`getGitTreeState`): when no tree state was
  injected and no VCS build info is present, the tree state is now resolved from a
  reachable git working copy. Uses the same hardened path as the existing git
  calls — whitelisted binary location, constant arguments, 5s `CommandContext`
  timeout — and reads only whether the output is empty, never its content. The
  `GitCmdStatus` / `GitArgPorcelain` constants had been declared since 1.0.0 and
  used by nothing.
- `applyVCSSettings(*Info, []debug.BuildSetting)`, split out of
  `applyBuildInfoGitData` so the tree-state decision is drivable from tests;
  `debug.ReadBuildInfo()` cannot be controlled from a test binary.
- `treeStateFromStatusOutput(string) string`, a pure mapping from git status
  output to a tree state.
- `Makefile` now stamps `GitTreeState`, so this repo's own CLI stops making a
  claim it could not support.
- README section "Git Tree State: clean, dirty, unknown" documenting the three
  readings and the four-step resolution order.

### Changed
- `DefaultGitTreeState` is now `GitTreeStateUnknown` (see the banner above).
- **`-X ...GitTreeState=clean` now takes effect.** `applyLdflagsGitInfo` decides
  "was anything injected?" by comparing the injected value against the default,
  which doubles as the not-injected sentinel. While that default was `"clean"`, a
  build system that *truthfully* injected `clean` was byte-indistinguishable from
  one that injected nothing and its statement was silently discarded; the output
  happened to be `clean` either way, so nothing ever showed it. With an `unknown`
  sentinel the comparison finally says what it was always meant to say. The
  comparison was deliberately kept rather than rewritten, and the reasoning is now
  a comment above it so the old default cannot be restored without noticing.
- **`vcs.modified=false` is now read as evidence of a clean tree.** The VCS arm
  previously read only the `"true"` case, so on a module build from a genuinely
  clean, VCS-stamped clone the toolchain had answered the question and the answer
  was dropped. Before this release that was invisible (the default was `clean`
  anyway); after it, leaving the arm unread would have traded a false `clean` for
  a lost true one. An explicit ldflags injection is never overridden by this arm.
- `GitInfo.TreeState`'s documentation now states all three readings and that
  `unknown` must never be read as `clean`. The field keeps no `omitempty`, so
  `unknown` is always on the wire and is never degraded to `""` — `""` remains a
  distinct (and invalid) reading meaning "nobody set this field".

### Fixed
- A binary that cannot determine its own tree state no longer claims it was clean.
  This is the whole point of the release: `/version` is the remedy for weak build
  provenance and it was able to lie.

### Tests
- New `treestate_test.go` with 12 guards covering: the default not being a
  confident `clean`; the three states being distinct, non-empty wire values; an
  injected `clean` reaching the `Info`; both `vcs.modified` readings; an absent
  and an unrecognised `vcs.modified` leaving the state `unknown`; ldflags
  precedence over the VCS clean arm; the status-output mapping; the git fallback
  not overriding an already-answered state; and `unknown` surviving JSON
  marshal/unmarshal and `LogFields()`.
- Each guard was verified by mutation — six mutations, each confirmed to COMPILE
  before being judged, each reddening a named test.
- `TestEnrichWithGitInfo` (pre-existing) passed vacuously before this release: it
  seeded the `Info` with `DefaultGitTreeState`, injected `GitTreeStateClean` and
  asserted `clean`, so the injection branch could have been deleted entirely
  without reddening it. It now has teeth, and carries a comment saying so.

## [1.0.0] - 2025-11-02

### Breaking Changes
- **CRITICAL**: Fixed YAML struct tags to use lowercase keys (aligning with YAML conventions)
  - **Migration required**: Update your `versions.yaml` files to use lowercase keys
  - Change `Schemas:` → `schemas:`
  - Change `APIs:` → `apis:`
  - Change `Components:` → `components:`
  - Change `Custom:` → `custom:`
  - This fixes a bug where the implementation used non-standard capitalized YAML keys
  - Template file (_templates/versions.yaml.tmpl) and documentation were already using correct lowercase format
  - **Why this change**: YAML convention uses lowercase/snake_case keys (like `manifest_version` and `project`), not mixed case
  - **Impact**: All existing manifest files must be updated to use lowercase keys
- Info struct map fields now unexported (use GetSchemas(), GetAPIs(), GetComponents(), GetCustom() instead)
  - True immutability with defensive copies
  - Prevents external mutation of internal state
- `Reset()` now uses `testing.Testing()` instead of GO_ENV (automatically detects test environment)

### Security
- Git binary PATH validation (whitelists /usr/bin, /usr/local/bin, /opt/homebrew/bin)
- Command injection protection with constant arguments
- Commit hash validation (hexadecimal format checks)
- Git command timeouts (5s) to prevent hangs
- Updated Go version requirement to 1.24.6+ (fixes 4 stdlib CVEs in Go 1.24.0)
- HTTP request size limits using `MaxBytesReader` (1KB limit, defense in depth)
- Production safeguards: `Reset()` panics if called outside test environment

### Added
- Core version information types and structures (Info, Manifest, ProjectVersion, GitInfo, BuildInfo)
- Multi-dimensional version manifest support (project, schemas, APIs, components, custom)
- Thread-safe singleton pattern with lock-free atomic reads
- Build-time injection via ldflags (GitCommit, GitTag, BuildTime, BuildUser)
- File-based manifest loading (versions.yaml)
- Embedded manifest support (go:embed)
- Build info extraction from runtime/debug
- Git info extraction (commit, tag, tree state)
- Context propagation to validators (WithContext option)
- Public semantic version API (ParseSemVer, MustParseSemVer, CompareVersions, IsNewerVersion)
- SemVer struct with full comparison methods (Compare, LessThan, GreaterThan, Equal, etc.)
- IsInitialized() to check state without triggering auto-init
- WithStrictMode() for production validation requirements
- Defensive copy getters (GetSchemas(), GetAPIs(), GetComponents(), GetCustom())
- Custom JSON marshaling for unexported fields
- Validation framework with Validator interface
- Built-in validators (NewSchemaValidator, NewAPIValidator, NewComponentValidator)
- HTTP handlers (/version endpoint, /health endpoint)
- HTTP middleware (adds version headers to responses)
- CLI tool with multiple output formats (full, compact, JSON, schemas, apis, components, git, build)
- Structured logging support with LogFields() for zap
- Comprehensive test suite (140+ tests, 90.5% coverage)
- Race detection verified (all tests pass with -race flag)
- 11 benchmark tests for performance tracking
- Security tests for git validation and command injection protection
- Complete documentation in README.md with examples
- YAML manifest template (_templates/versions.yaml.tmpl)
- Working examples (simple and advanced)

### Fixed
- YAML struct tags now use lowercase keys (standards compliance)
- All example and test YAML files updated to use lowercase keys
- Thread safety with true immutability (defensive copies)
- Reset() now properly detects test environment automatically

## [0.1.0] - 2025-10-11

### Added
- Initial project structure
- Go module initialization
- Core dependencies (go-cuserr, zap, yaml.v3)
- Makefile with build, test, lint targets
- golangci-lint configuration
- CLAUDE.md for development guidance
- Project plan documentation
- MIT License
- VERSION file
- Directory structure for implementation

[Unreleased]: https://github.com/itsatony/go-version/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/itsatony/go-version/compare/v0.1.0...v1.0.0
[0.1.0]: https://github.com/itsatony/go-version/releases/tag/v0.1.0
