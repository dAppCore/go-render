//go:build !js

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildComponentJS_ValidJSON_Good(t *testing.T) {
	slotsJSON := `{"H":"nav-bar","C":"main-content"}`
	js, err := buildComponentJS(slotsJSON)
	require.NoError(t, err)
	assert.Contains(t, js, "NavBar")
	assert.Contains(t, js, "MainContent")
	assert.Contains(t, js, "customElements.define")
}

func TestBuildComponentJS_InvalidJSON_Bad(t *testing.T) {
	_, err := buildComponentJS("not json")
	assert.Error(t, err)
}
