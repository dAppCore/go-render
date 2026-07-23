//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ctmltest

import (
	"bytes"
	"image/png"
	"testing"

	coreio "dappco.re/go/io"
	"github.com/charmbracelet/x/vttest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildImage(t *testing.T) {
	t.Run("good: produces a valid, non-empty PNG sized cols x rows cells", func(t *testing.T) {
		data, err := buildImage("Welcome\nBuild 1.0", 10)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		img, err := png.Decode(bytes.NewReader(data))
		require.NoError(t, err, "buildImage must produce a decodable PNG")
		bounds := img.Bounds()
		assert.Equal(t, 10*vttest.DefaultDrawer.CellWidth, bounds.Dx())
		assert.Equal(t, 2*vttest.DefaultDrawer.CellHeight, bounds.Dy())
	})

	t.Run("good: the same frame and cols encode to byte-identical PNGs, run to run", func(t *testing.T) {
		a, err := buildImage("Welcome", 20)
		require.NoError(t, err)
		b, err := buildImage("Welcome", 20)
		require.NoError(t, err)
		assert.Equal(t, a, b)
	})

	t.Run("ugly: a single-cell frame still encodes to a valid, non-empty PNG", func(t *testing.T) {
		data, err := buildImage("x", 1)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		img, err := png.Decode(bytes.NewReader(data))
		require.NoError(t, err)
		assert.Equal(t, vttest.DefaultDrawer.CellWidth, img.Bounds().Dx())
		assert.Equal(t, vttest.DefaultDrawer.CellHeight, img.Bounds().Dy())
	})
}

func TestWriteImage(t *testing.T) {
	t.Run("good: writes the exact PNG bytes given", func(t *testing.T) {
		dir := t.TempDir()
		cmd := command{Args: []string{"out.png"}, Line: 1}
		data, err := buildImage("hi", 4)
		require.NoError(t, err)

		require.NoError(t, writeImage(dir, cmd, data))

		raw, err := coreio.Local.Read(dir + "/out.png")
		require.NoError(t, err)
		assert.Equal(t, data, []byte(raw))
	})

	t.Run("good: creates a missing image directory", func(t *testing.T) {
		dir := t.TempDir() + "/nested"
		cmd := command{Args: []string{"out.png"}, Line: 1}
		data, err := buildImage("hi", 4)
		require.NoError(t, err)

		require.NoError(t, writeImage(dir, cmd, data))

		raw, err := coreio.Local.Read(dir + "/out.png")
		require.NoError(t, err)
		assert.Equal(t, data, []byte(raw))
	})
}
