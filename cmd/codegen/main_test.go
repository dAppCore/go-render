//go:build !js

package main

import (
	"testing"

	core "dappco.re/go/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_WritesBundle(t *testing.T) {
	input := core.NewReader(`{"H":"nav-bar","C":"main-content"}`)
	output := core.NewBuilder()

	err := run(input, output)
	require.NoError(t, err)

	js := output.String()
	assert.Contains(t, js, "NavBar")
	assert.Contains(t, js, "MainContent")
	assert.Contains(t, js, "customElements.define")
	assert.Equal(t, 2, countSubstr(js, "extends HTMLElement"))
}

func TestRun_InvalidJSON(t *testing.T) {
	input := core.NewReader(`not json`)
	output := core.NewBuilder()

	err := run(input, output)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid JSON")
}

func TestRun_InvalidTag(t *testing.T) {
	input := core.NewReader(`{"H":"notag"}`)
	output := core.NewBuilder()

	err := run(input, output)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hyphen")
}

func TestRun_EmptySlots(t *testing.T) {
	input := core.NewReader(`{}`)
	output := core.NewBuilder()

	err := run(input, output)
	require.NoError(t, err)
	assert.Empty(t, output.String())
}

func countSubstr(s, substr string) int {
	if substr == "" {
		return len(s) + 1
	}

	count := 0
	for i := 0; i <= len(s)-len(substr); {
		j := indexSubstr(s[i:], substr)
		if j < 0 {
			return count
		}
		count++
		i += j + len(substr)
	}

	return count
}

func indexSubstr(s, substr string) int {
	if substr == "" {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}

	return -1
}
