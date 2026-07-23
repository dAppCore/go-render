//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ctmltest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTape_Good(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []command
	}{
		{
			name: "good: Source takes one path argument",
			src:  "Source settings.ctml",
			want: []command{{Verb: "Source", Args: []string{"settings.ctml"}, Line: 1}},
		},
		{
			name: "good: Set Width takes a key and a positive integer",
			src:  "Set Width 80",
			want: []command{{Verb: "Set", Args: []string{"Width", "80"}, Line: 1}},
		},
		{
			name: "good: Set Height takes a key and a positive integer",
			src:  "Set Height 24",
			want: []command{{Verb: "Set", Args: []string{"Height", "24"}, Line: 1}},
		},
		{
			name: "good: Set Theme takes a key and a non-empty name",
			src:  "Set Theme midnight",
			want: []command{{Verb: "Set", Args: []string{"Theme", "midnight"}, Line: 1}},
		},
		{
			name: "good: Data takes a key and a quoted value",
			src:  `Data session.title "Welcome"`,
			want: []command{{Verb: "Data", Args: []string{"session.title", "Welcome"}, Line: 1}},
		},
		{
			name: "good: Data value need not be quoted when it has no spaces",
			src:  "Data version 1.0",
			want: []command{{Verb: "Data", Args: []string{"version", "1.0"}, Line: 1}},
		},
		{
			name: "good: Rows takes a sequence name and a row count",
			src:  "Rows items 3",
			want: []command{{Verb: "Rows", Args: []string{"items", "3"}, Line: 1}},
		},
		{
			name: "good: Expect Text takes a quoted substring",
			src:  `Expect Text "hello world"`,
			want: []command{{Verb: "Expect", Args: []string{"Text", "hello world"}, Line: 1}},
		},
		{
			name: "good: Expect Box takes a block id",
			src:  "Expect Box row-2",
			want: []command{{Verb: "Expect", Args: []string{"Box", "row-2"}, Line: 1}},
		},
		{
			name: "good: Expect Fits takes no further argument",
			src:  "Expect Fits",
			want: []command{{Verb: "Expect", Args: []string{"Fits"}, Line: 1}},
		},
		{
			name: "good: Golden takes one path argument",
			src:  "Golden settings.golden",
			want: []command{{Verb: "Golden", Args: []string{"settings.golden"}, Line: 1}},
		},
		{
			name: "good: a trailing # comment after real arguments is dropped",
			src:  "Set Width 80 # eighty columns",
			want: []command{{Verb: "Set", Args: []string{"Width", "80"}, Line: 1}},
		},
		{
			name: "good: a whole-line comment produces no command",
			src:  "# this is a comment\nSource a.ctml",
			want: []command{{Verb: "Source", Args: []string{"a.ctml"}, Line: 2}},
		},
		{
			name: "good: blank lines produce no command and still advance line numbers",
			src:  "Source a.ctml\n\n\nExpect Fits",
			want: []command{
				{Verb: "Source", Args: []string{"a.ctml"}, Line: 1},
				{Verb: "Expect", Args: []string{"Fits"}, Line: 4},
			},
		},
		{
			name: "good: a quoted argument may contain a literal # and spaces",
			src:  `Expect Text "50% #1"`,
			want: []command{{Verb: "Expect", Args: []string{"Text", "50% #1"}, Line: 1}},
		},
		{
			name: "good: CRLF line endings are normalised",
			src:  "Source a.ctml\r\nExpect Fits\r\n",
			want: []command{
				{Verb: "Source", Args: []string{"a.ctml"}, Line: 1},
				{Verb: "Expect", Args: []string{"Fits"}, Line: 2},
			},
		},
		{
			name: "good: Rows accepts a row count of exactly zero",
			src:  "Rows items 0",
			want: []command{{Verb: "Rows", Args: []string{"items", "0"}, Line: 1}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseTape([]byte(tc.src))
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestParseTape_Bad(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		wantLine int
	}{
		{"bad: unknown verb", "Bogus arg", 1},
		{"bad: Source with no argument", "Source", 1},
		{"bad: Source with two arguments", "Source a.ctml extra", 1},
		{"bad: Set with one argument", "Set Width", 1},
		{"bad: Set with three arguments", "Set Width 80 extra", 1},
		{"bad: Set with an unrecognised key", "Set Speed 80", 1},
		{"bad: Set Width with a non-integer value", "Set Width wide", 1},
		{"bad: Set Width with a negative value", "Set Width -1", 1},
		{"bad: Set Width with zero", "Set Width 0", 1},
		{"bad: Set Height with a non-integer value", "Set Height tall", 1},
		{"bad: Set Theme with an empty value", `Set Theme ""`, 1},
		{"bad: Data with one argument", "Data key", 1},
		{"bad: Data with three arguments", "Data key value extra", 1},
		{"bad: Rows with one argument", "Rows items", 1},
		{"bad: Rows with a non-integer count", "Rows items many", 1},
		{"bad: Rows with a negative count", "Rows items -1", 1},
		{"bad: bare Expect with no kind", "Expect", 1},
		{"bad: Expect with an unrecognised kind", "Expect Colour red", 1},
		{"bad: Expect Text with no substring", "Expect Text", 1},
		{"bad: Expect Text with two arguments", "Expect Text a b", 1},
		{"bad: Expect Box with no id", "Expect Box", 1},
		{"bad: Expect Fits with an extra argument", "Expect Fits now", 1},
		{"bad: Golden with no argument", "Golden", 1},
		{"bad: Golden with two arguments", "Golden a.golden extra", 1},
		{"bad: an unterminated quoted argument", `Expect Text "unterminated`, 1},
		{"bad: a later line reports its own line number", "Source a.ctml\nSet Width 80\nBogus", 3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseTape([]byte(tc.src))
			require.Error(t, err)
			assert.Nil(t, got)
			var tapeErr *TapeError
			require.ErrorAs(t, err, &tapeErr)
			assert.Equal(t, tc.wantLine, tapeErr.Line)
		})
	}
}

func TestParseTape_Ugly(t *testing.T) {
	t.Run("ugly: an empty tape parses to zero commands, not an error", func(t *testing.T) {
		got, err := parseTape([]byte(""))
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("ugly: a tape that is only comments and blank lines parses to zero commands", func(t *testing.T) {
		got, err := parseTape([]byte("# nothing here\n\n   \n# still nothing"))
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("ugly: leading and trailing whitespace around tokens is insignificant", func(t *testing.T) {
		got, err := parseTape([]byte("   Set   Width   80   "))
		require.NoError(t, err)
		assert.Equal(t, []command{{Verb: "Set", Args: []string{"Width", "80"}, Line: 1}}, got)
	})
}

func TestTapeError_Error(t *testing.T) {
	t.Run("good: formats as ctmltest:line: msg", func(t *testing.T) {
		err := &TapeError{Line: 5, Msg: `unknown verb "Bogus"`}
		assert.Equal(t, `ctmltest:5: unknown verb "Bogus"`, err.Error())
	})

	t.Run("ugly: a nil receiver returns an empty string, not a panic", func(t *testing.T) {
		var err *TapeError
		assert.Equal(t, "", err.Error())
	})
}

func TestTapeError_Unwrap(t *testing.T) {
	t.Run("good: exposes the wrapped cause", func(t *testing.T) {
		cause := assert.AnError
		err := &TapeError{Line: 1, Msg: "x", Cause: cause}
		assert.Same(t, cause, err.Unwrap())
	})

	t.Run("ugly: a nil receiver returns a nil cause, not a panic", func(t *testing.T) {
		var err *TapeError
		assert.Nil(t, err.Unwrap())
	})
}
