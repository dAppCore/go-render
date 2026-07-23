//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ctmltest

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRun_Good and TestRunFile_Good nest a real (sub)test that is expected
// to PASS -- safe, because only a FAILING (sub)test propagates its failure
// to every ancestor and the process exit code (confirmed empirically: Go's
// testing package has no way to run a real *testing.T to a deliberate
// failure without that failure poisoning the outer run). Every failure
// path below is therefore tested via matchTapes/prepareRun directly --
// plain functions returning a plain error, not *testing.T -- see
// runExpect's doc comment in expect.go for the same reasoning applied to
// Expect/Golden.

func TestRun_Good(t *testing.T) {
	t.Run("good: a glob matching one tape runs it as a passing subtest", func(t *testing.T) {
		ok := t.Run("inner", func(t *testing.T) {
			Run(t, "testdata/*_test.ctml")
		})
		assert.True(t, ok, "expected Run(testdata/*_test.ctml) to pass")
	})
}

func TestRunFile_Good(t *testing.T) {
	t.Run("good: a well-formed tape runs clean", func(t *testing.T) {
		ok := t.Run("inner", func(t *testing.T) {
			RunFile(t, "testdata/sample_test.ctml")
		})
		assert.True(t, ok, "expected RunFile(testdata/sample_test.ctml) to pass")
	})
}

func TestMatchTapes(t *testing.T) {
	t.Run("good: a matching glob returns a sorted, non-empty path list", func(t *testing.T) {
		paths, err := matchTapes("testdata/*_test.ctml")
		require.NoError(t, err)
		assert.Contains(t, paths, "testdata/sample_test.ctml")
		assert.True(t, sort.StringsAreSorted(paths))
	})

	t.Run("bad: a glob matching nothing is an error naming the glob", func(t *testing.T) {
		paths, err := matchTapes("testdata/does-not-exist/*_test.ctml")
		require.Error(t, err)
		assert.Nil(t, paths)
		assert.Contains(t, err.Error(), "testdata/does-not-exist/*_test.ctml")
	})
}

func TestPrepareRun_Good(t *testing.T) {
	t.Run("good: a well-formed tape resolves a full render result", func(t *testing.T) {
		result, err := prepareRun("testdata/sample_test.ctml")
		require.NoError(t, err)
		assert.NotEmpty(t, result.cmds)
		assert.Equal(t, "testdata", result.tapeDir)
		assert.Contains(t, result.frame, "Welcome")
		assert.Contains(t, result.boxes, "banner")
		assert.Equal(t, 40, result.fitWidth)
	})
}

func TestPrepareRun_Bad(t *testing.T) {
	tests := []struct {
		name        string
		tapePath    string
		wantMessage string
	}{
		{
			name:        "bad: the tape file itself does not exist",
			tapePath:    "testdata/edge/does-not-exist_test.ctml",
			wantMessage: "reading tape",
		},
		{
			name:        "bad: a tape with a parse error",
			tapePath:    "testdata/edge/bad_verb_test.ctml",
			wantMessage: `unknown verb "Bogus"`,
		},
		{
			name:        "bad: a tape missing the Source verb",
			tapePath:    "testdata/edge/no_source_test.ctml",
			wantMessage: "missing Source verb",
		},
		{
			name:        "bad: a tape whose Source file does not exist",
			tapePath:    "testdata/edge/missing_source_file_test.ctml",
			wantMessage: "reading Source",
		},
		{
			name:        "bad: a tape whose Source .ctml fails to parse",
			tapePath:    "testdata/edge/bad_ctml_test.ctml",
			wantMessage: "parsing Source",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := prepareRun(tc.tapePath)
			require.Error(t, err)
			assert.Zero(t, result)
			assert.Contains(t, err.Error(), tc.wantMessage)
		})
	}
}
