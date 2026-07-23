//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ctmltest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildConfig_Good(t *testing.T) {
	t.Run("good: Source resolves relative to the tape directory", func(t *testing.T) {
		cmds := []command{{Verb: "Source", Args: []string{"a.ctml"}, Line: 1}}
		cfg := buildConfig("testdata", cmds)
		assert.Equal(t, "testdata/a.ctml", cfg.sourcePath)
		assert.Equal(t, 1, cfg.sourceLine)
	})

	t.Run("good: Set Width/Height/Theme populate their fields", func(t *testing.T) {
		cmds := []command{
			{Verb: "Set", Args: []string{"Width", "80"}, Line: 1},
			{Verb: "Set", Args: []string{"Height", "24"}, Line: 2},
			{Verb: "Set", Args: []string{"Theme", "midnight"}, Line: 3},
		}
		cfg := buildConfig(".", cmds)
		assert.Equal(t, 80, cfg.width)
		assert.Equal(t, 24, cfg.height)
		assert.Equal(t, "midnight", cfg.theme)
	})

	t.Run("good: a repeated Set key keeps its last value", func(t *testing.T) {
		cmds := []command{
			{Verb: "Set", Args: []string{"Width", "80"}, Line: 1},
			{Verb: "Set", Args: []string{"Width", "120"}, Line: 2},
		}
		cfg := buildConfig(".", cmds)
		assert.Equal(t, 120, cfg.width)
	})

	t.Run("good: Data seeds a flat value", func(t *testing.T) {
		cmds := []command{{Verb: "Data", Args: []string{"version", "1.0"}, Line: 1}}
		cfg := buildConfig(".", cmds)
		assert.Equal(t, map[string]any{"version": "1.0"}, cfg.values)
	})

	t.Run("good: Data seeds a dotted value as a nested map", func(t *testing.T) {
		cmds := []command{{Verb: "Data", Args: []string{"session.title", "Welcome"}, Line: 1}}
		cfg := buildConfig(".", cmds)
		assert.Equal(t, map[string]any{"session": map[string]any{"title": "Welcome"}}, cfg.values)
	})

	t.Run("good: Rows seeds a generated sequence under its own name", func(t *testing.T) {
		cmds := []command{{Verb: "Rows", Args: []string{"items", "2"}, Line: 1}}
		cfg := buildConfig(".", cmds)
		assert.Equal(t, buildRows("items", 2), cfg.sequences["items"])
		assert.Len(t, cfg.sequences["items"], 2)
	})

	t.Run("good: an empty tape resolves to a usable zero config", func(t *testing.T) {
		cfg := buildConfig(".", nil)
		assert.Empty(t, cfg.sourcePath)
		assert.Empty(t, cfg.width)
		assert.NotNil(t, cfg.values)
		assert.NotNil(t, cfg.sequences)
	})
}

func TestSetDotted_Good(t *testing.T) {
	t.Run("good: a bare key sets a top-level value", func(t *testing.T) {
		values := map[string]any{}
		setDotted(values, "count", "3")
		assert.Equal(t, map[string]any{"count": "3"}, values)
	})

	t.Run("good: a dotted key builds a nested map", func(t *testing.T) {
		values := map[string]any{}
		setDotted(values, "user.name", "Ada")
		assert.Equal(t, map[string]any{"user": map[string]any{"name": "Ada"}}, values)
	})

	t.Run("good: a two-dot key builds two nested levels", func(t *testing.T) {
		values := map[string]any{}
		setDotted(values, "a.b.c", "x")
		assert.Equal(t, map[string]any{"a": map[string]any{"b": map[string]any{"c": "x"}}}, values)
	})

	t.Run("good: two keys sharing a prefix share the same nested map", func(t *testing.T) {
		values := map[string]any{}
		setDotted(values, "session.title", "Welcome")
		setDotted(values, "session.id", "42")
		assert.Equal(t, map[string]any{"session": map[string]any{"title": "Welcome", "id": "42"}}, values)
	})

	t.Run("ugly: a later flat write overwrites an earlier nested map at the same key", func(t *testing.T) {
		values := map[string]any{}
		setDotted(values, "session.title", "Welcome")
		setDotted(values, "session", "flat")
		assert.Equal(t, map[string]any{"session": "flat"}, values)
	})
}

func TestBuildRows_Good(t *testing.T) {
	t.Run("good: each row carries index, n, name, and label", func(t *testing.T) {
		rows := buildRows("items", 3)
		assert.Equal(t, []map[string]any{
			{"index": 0, "n": 1, "name": "items-0", "label": "Row 1"},
			{"index": 1, "n": 2, "name": "items-1", "label": "Row 2"},
			{"index": 2, "n": 3, "name": "items-2", "label": "Row 3"},
		}, rows)
	})

	t.Run("ugly: a row count of zero returns an empty, non-nil slice", func(t *testing.T) {
		rows := buildRows("items", 0)
		assert.NotNil(t, rows)
		assert.Empty(t, rows)
	})
}
