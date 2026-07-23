//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ctmltest

import (
	"testing"

	html "dappco.re/go/render/engine/html"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsRenderRead(t *testing.T) {
	tests := []struct {
		verb string
		want bool
	}{
		{"Expect", true},
		{"Golden", true},
		{"Click", true},
		{"Snapshot", true},
		{"Image", true},
		{"Source", false},
		{"Set", false},
		{"Data", false},
		{"Rows", false},
		{"Bogus", false},
	}
	for _, tc := range tests {
		t.Run(tc.verb, func(t *testing.T) {
			assert.Equal(t, tc.want, isRenderRead(tc.verb))
		})
	}
}

func TestCmdsBeforeFirstAssertion(t *testing.T) {
	t.Run("good: stops before the first Expect, leaving later Data out", func(t *testing.T) {
		cmds := []command{
			{Verb: "Source", Args: []string{"a.ctml"}, Line: 1},
			{Verb: "Data", Args: []string{"k", "v"}, Line: 2},
			{Verb: "Expect", Args: []string{"Fits"}, Line: 3},
			{Verb: "Data", Args: []string{"k", "v2"}, Line: 4},
		}
		assert.Equal(t, cmds[:2], cmdsBeforeFirstAssertion(cmds))
	})

	t.Run("good: stops before the first Click", func(t *testing.T) {
		cmds := []command{
			{Verb: "Source", Args: []string{"a.ctml"}, Line: 1},
			{Verb: "Click", Args: []string{"banner"}, Line: 2},
		}
		assert.Equal(t, cmds[:1], cmdsBeforeFirstAssertion(cmds))
	})

	t.Run("good: stops before the first Golden", func(t *testing.T) {
		cmds := []command{
			{Verb: "Source", Args: []string{"a.ctml"}, Line: 1},
			{Verb: "Golden", Args: []string{"a.golden"}, Line: 2},
		}
		assert.Equal(t, cmds[:1], cmdsBeforeFirstAssertion(cmds))
	})

	t.Run("ugly: no render-reading command yields the whole slice", func(t *testing.T) {
		cmds := []command{
			{Verb: "Source", Args: []string{"a.ctml"}, Line: 1},
			{Verb: "Data", Args: []string{"k", "v"}, Line: 2},
		}
		assert.Equal(t, cmds, cmdsBeforeFirstAssertion(cmds))
	})

	t.Run("ugly: a render-reading command first yields an empty prefix", func(t *testing.T) {
		cmds := []command{{Verb: "Expect", Args: []string{"Fits"}, Line: 1}}
		assert.Empty(t, cmdsBeforeFirstAssertion(cmds))
	})

	t.Run("ugly: an empty tape yields an empty slice", func(t *testing.T) {
		assert.Empty(t, cmdsBeforeFirstAssertion(nil))
	})
}

func TestNewDriveState(t *testing.T) {
	t.Run("good: copies every field from tapePath and result", func(t *testing.T) {
		result := renderResult{
			frame: "hello",
			boxes: html.BoxMap{"banner": {Width: 1, Height: 1}},
			cfg: tapeConfig{
				sourcePath: "a.ctml",
				width:      40,
				values:     map[string]any{"x": "1"},
				sequences:  map[string][]map[string]any{"items": {{"n": 1}}},
			},
			ctmlSrc: []byte("<p>x</p>"),
		}
		drive := newDriveState("tape_test.ctml", result)
		assert.Equal(t, "tape_test.ctml", drive.tapePath)
		assert.Equal(t, "a.ctml", drive.sourcePath)
		assert.Equal(t, []byte("<p>x</p>"), drive.ctmlSrc)
		assert.Equal(t, 40, drive.width)
		assert.Equal(t, "hello", drive.frame)
		assert.Equal(t, result.boxes, drive.boxes)
		assert.Equal(t, result.cfg.values, drive.values)
		assert.Equal(t, result.cfg.sequences, drive.sequences)
	})

	t.Run("good: the live values map is result.cfg's own, not a copy", func(t *testing.T) {
		values := map[string]any{"x": "1"}
		result := renderResult{cfg: tapeConfig{values: values, sequences: map[string][]map[string]any{}}}
		drive := newDriveState("tape_test.ctml", result)
		drive.values["x"] = "2"
		assert.Equal(t, "2", values["x"])
	})
}

func TestDriveStateRedrive(t *testing.T) {
	t.Run("good: merges the Data command into values and re-renders", func(t *testing.T) {
		drive := &driveState{
			tapePath:   "tape.ctml",
			sourcePath: "widget.ctml",
			ctmlSrc:    []byte("<p>{{msg}}</p>"),
			width:      40,
			values:     map[string]any{"msg": "old"},
			sequences:  map[string][]map[string]any{},
			frame:      "old",
		}
		err := drive.redrive(command{Args: []string{"msg", "new"}, Line: 9})
		require.NoError(t, err)
		assert.Equal(t, "new", drive.values["msg"])
		assert.Contains(t, drive.frame, "new")
		assert.NotContains(t, drive.frame, "old")
	})

	t.Run("good: a dotted key merges via setDotted, same as a leading Data line", func(t *testing.T) {
		drive := &driveState{
			ctmlSrc:   []byte("<p>{{user.name}}</p>"),
			values:    map[string]any{},
			sequences: map[string][]map[string]any{},
		}
		err := drive.redrive(command{Args: []string{"user.name", "Ada"}, Line: 1})
		require.NoError(t, err)
		assert.Equal(t, map[string]any{"user": map[string]any{"name": "Ada"}}, drive.values)
		assert.Contains(t, drive.frame, "Ada")
	})

	t.Run("good: boxes are replaced with the new render's box map", func(t *testing.T) {
		drive := &driveState{
			ctmlSrc:   []byte(`<p id="msg">{{msg}}</p>`),
			values:    map[string]any{"msg": "hi"},
			sequences: map[string][]map[string]any{},
		}
		err := drive.redrive(command{Args: []string{"msg", "hello"}, Line: 1})
		require.NoError(t, err)
		assert.Contains(t, drive.boxes, "msg")
	})

	t.Run("bad: a re-Parse failure is wrapped and names the Data line", func(t *testing.T) {
		drive := &driveState{
			tapePath:   "tape.ctml",
			sourcePath: "bad.ctml",
			ctmlSrc:    []byte("<if><p>x</p></if>"), // same deliberately-invalid markup as testdata/edge/bad.ctml
			values:     map[string]any{},
			sequences:  map[string][]map[string]any{},
		}
		err := drive.redrive(command{Args: []string{"k", "v"}, Line: 7})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "tape.ctml:7")
		assert.Contains(t, err.Error(), "parsing Source bad.ctml")
	})
}
