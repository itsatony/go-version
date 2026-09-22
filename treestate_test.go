package version

import (
	"encoding/json"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDefaultGitTreeStateIsNotAConfidentClean pins the headline defect: the value
// used when the tree state could NOT be determined must not be a confident "clean".
// "clean" and "unknown" are two different facts and must not share a field.
func TestDefaultGitTreeStateIsNotAConfidentClean(t *testing.T) {
	assert.NotEqual(t, GitTreeStateClean, DefaultGitTreeState,
		"the fallback for an undeterminable tree state must not claim the tree was clean")
	assert.NotEqual(t, GitTreeStateDirty, DefaultGitTreeState,
		"the fallback for an undeterminable tree state must not claim the tree was dirty")
	assert.NotEmpty(t, DefaultGitTreeState,
		"the fallback must be a stated value, not an empty string a caller cannot distinguish from unset")
}

// TestManifestDefaultTreeStateIsNotClean covers the second place the default is
// written: the Info built from a manifest, before any git enrichment runs.
func TestManifestDefaultTreeStateIsNotClean(t *testing.T) {
	info := manifestToInfo(defaultManifest())
	assert.NotEqual(t, GitTreeStateClean, info.Git.TreeState,
		"an un-enriched Info must not claim a clean tree")
	assert.Equal(t, DefaultGitTreeState, info.Git.TreeState)
}

// TestLdflagsInjectedCleanTreeStateIsApplied is the loader's sentinel test.
// applyLdflagsGitInfo skips injection when the injected value equals the default,
// so while the default IS "clean" a build system that truthfully injects
// -X ...GitTreeState=clean is indistinguishable from one that injects nothing,
// and its statement is silently discarded.
func TestLdflagsInjectedCleanTreeStateIsApplied(t *testing.T) {
	orig := GitTreeState
	t.Cleanup(func() { GitTreeState = orig })

	info := &Info{Git: GitInfo{Commit: DefaultGitCommit, TreeState: "unknown"}}
	GitTreeState = GitTreeStateClean

	applyLdflagsGitInfo(info)

	assert.Equal(t, GitTreeStateClean, info.Git.TreeState,
		"an explicitly injected clean tree state must reach the Info")
}

// TestTreeStateVocabularyIsThreeDistinctFacts keeps "could not be determined"
// from collapsing back into "clean" (or into the empty string, which a caller
// cannot distinguish from a field nobody set).
func TestTreeStateVocabularyIsThreeDistinctFacts(t *testing.T) {
	states := map[string]string{
		"clean":   GitTreeStateClean,
		"dirty":   GitTreeStateDirty,
		"unknown": GitTreeStateUnknown,
	}
	seen := map[string]string{}
	for name, value := range states {
		assert.NotEmpty(t, value, "%s must have a stated wire value", name)
		if other, dup := seen[value]; dup {
			t.Errorf("tree states %q and %q share the wire value %q", name, other, value)
		}
		seen[value] = name
	}
	assert.Equal(t, GitTreeStateUnknown, DefaultGitTreeState,
		"the undeterminable default must be the unknown member of the vocabulary")
}

// TestVCSModifiedFalseIsEvidenceOfACleanTree: vcs.modified carries both answers.
// Reading only the "true" arm would leave every clean module build reporting
// "unknown" while the toolchain had in fact answered the question.
func TestVCSModifiedFalseIsEvidenceOfACleanTree(t *testing.T) {
	info := &Info{Git: GitInfo{Commit: DefaultGitCommit, TreeState: DefaultGitTreeState}}

	applyVCSSettings(info, []debug.BuildSetting{{Key: VCSKeyModified, Value: VCSValueFalse}})

	assert.Equal(t, GitTreeStateClean, info.Git.TreeState)
}

func TestVCSModifiedTrueIsEvidenceOfADirtyTree(t *testing.T) {
	info := &Info{Git: GitInfo{Commit: DefaultGitCommit, TreeState: DefaultGitTreeState}}

	applyVCSSettings(info, []debug.BuildSetting{{Key: VCSKeyModified, Value: VCSValueTrue}})

	assert.Equal(t, GitTreeStateDirty, info.Git.TreeState)
}

// TestVCSModifiedAbsentLeavesTheTreeStateUnknown is the case that reaches the
// decision FIRST with no evidence at all: no vcs.modified setting means nobody
// answered, and the answer must stay "unknown".
func TestVCSModifiedAbsentLeavesTheTreeStateUnknown(t *testing.T) {
	info := &Info{Git: GitInfo{Commit: DefaultGitCommit, TreeState: DefaultGitTreeState}}

	applyVCSSettings(info, []debug.BuildSetting{
		{Key: VCSKeyRevision, Value: "0123456789abcdef0123456789abcdef01234567"},
		{Key: VCSKeyTime, Value: "2026-09-22T00:00:00Z"},
	})

	assert.Equal(t, GitTreeStateUnknown, info.Git.TreeState)
	assert.Equal(t, "0123456789abcdef0123456789abcdef01234567", info.Git.Commit)
}

// TestVCSModifiedFalseDoesNotOverrideAnInjectedTreeState reaches the clean branch
// with the state ALREADY answered by ldflags, which is the ordering that decides
// which source wins. An explicit injection is the operator's statement and must
// survive a weaker one.
func TestVCSModifiedFalseDoesNotOverrideAnInjectedTreeState(t *testing.T) {
	info := &Info{Git: GitInfo{Commit: DefaultGitCommit, TreeState: GitTreeStateDirty}}

	applyVCSSettings(info, []debug.BuildSetting{{Key: VCSKeyModified, Value: VCSValueFalse}})

	assert.Equal(t, GitTreeStateDirty, info.Git.TreeState,
		"build info saying 'not modified' must not erase an explicit dirty injection")
}

// TestVCSModifiedUnrecognisedValueLeavesTheStateUnknown: a value that is neither
// "true" nor "false" is not evidence of anything.
func TestVCSModifiedUnrecognisedValueLeavesTheStateUnknown(t *testing.T) {
	info := &Info{Git: GitInfo{Commit: DefaultGitCommit, TreeState: DefaultGitTreeState}}

	applyVCSSettings(info, []debug.BuildSetting{{Key: VCSKeyModified, Value: "maybe"}})

	assert.Equal(t, GitTreeStateUnknown, info.Git.TreeState)
}

func TestTreeStateFromStatusOutput(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   string
	}{
		{"empty output is a clean tree", "", GitTreeStateClean},
		{"whitespace only is a clean tree", "\n", GitTreeStateClean},
		{"modified file is a dirty tree", " M version.loader.go\n", GitTreeStateDirty},
		{"untracked file is a dirty tree", "?? new.go\n", GitTreeStateDirty},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, treeStateFromStatusOutput(tt.output))
		})
	}
}

// TestGitCommandFallbackDoesNotOverrideAnAnsweredTreeState reaches the fallback
// with the question already answered - the case that must reach the guard FIRST.
func TestGitCommandFallbackDoesNotOverrideAnAnsweredTreeState(t *testing.T) {
	// A sentinel git status can never produce, so this assertion is independent of
	// whether the checkout running the suite happens to be clean or dirty.
	const alreadyAnswered = "answered-by-an-earlier-source"
	info := &Info{Git: GitInfo{Commit: "0123456789abcdef0123456789abcdef01234567", TreeState: alreadyAnswered}}

	applyGitCommandFallback(info)

	assert.Equal(t, alreadyAnswered, info.Git.TreeState,
		"the git-status fallback must only speak when the state is still unknown")
}

// TestUnknownTreeStateSurvivesTheWire proves "unknown" is representable end to
// end: it is serialized, it is never omitted, and it never degrades to "".
func TestUnknownTreeStateSurvivesTheWire(t *testing.T) {
	info := manifestToInfo(defaultManifest())

	encoded, err := json.Marshal(info)
	require.NoError(t, err)
	assert.Contains(t, string(encoded), `"tree_state":"unknown"`)
	assert.NotContains(t, string(encoded), `"tree_state":"clean"`)

	var decoded Info
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	assert.Equal(t, GitTreeStateUnknown, decoded.Git.TreeState)

	// An empty tree state is a DIFFERENT (and invalid) reading from "unknown";
	// the two must not be collapsed by the omitempty/zero-value machinery.
	empty := &Info{Git: GitInfo{TreeState: ""}}
	emptyEncoded, err := json.Marshal(empty)
	require.NoError(t, err)
	assert.Contains(t, string(emptyEncoded), `"tree_state":""`)
}

// TestLogFieldsCarryTheUnknownTreeState covers the second reporting surface.
func TestLogFieldsCarryTheUnknownTreeState(t *testing.T) {
	info := manifestToInfo(defaultManifest())
	for _, f := range info.LogFields() {
		if f.Key == LogFieldGitTreeState {
			assert.Equal(t, GitTreeStateUnknown, f.String)
			return
		}
	}
	t.Fatalf("no %s field in LogFields()", LogFieldGitTreeState)
}
