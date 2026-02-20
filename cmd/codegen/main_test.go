package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_Good(t *testing.T) {
	input := strings.NewReader(`{"H":"nav-bar","C":"main-content"}`)
	var output bytes.Buffer

	err := run(input, &output)
	require.NoError(t, err)

	js := output.String()
	assert.Contains(t, js, "NavBar")
	assert.Contains(t, js, "MainContent")
	assert.Contains(t, js, "customElements.define")
	assert.Equal(t, 2, strings.Count(js, "extends HTMLElement"))
}

func TestRun_Bad_InvalidJSON(t *testing.T) {
	input := strings.NewReader(`not json`)
	var output bytes.Buffer

	err := run(input, &output)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid JSON")
}

func TestRun_Bad_InvalidTag(t *testing.T) {
	input := strings.NewReader(`{"H":"notag"}`)
	var output bytes.Buffer

	err := run(input, &output)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hyphen")
}

func TestRun_Good_Empty(t *testing.T) {
	input := strings.NewReader(`{}`)
	var output bytes.Buffer

	err := run(input, &output)
	require.NoError(t, err)
	assert.Empty(t, output.String())
}
