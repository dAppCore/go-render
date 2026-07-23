//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ctmltest

import (
	"testing"

	html "dappco.re/go/html"
	"github.com/stretchr/testify/assert"
)

func TestEvalClick(t *testing.T) {
	boxes := html.BoxMap{"banner": {Row: 0, Col: 0, Width: 10, Height: 1}}
	frame := "Welcome"

	t.Run("good: a recorded, non-empty box is hit-testable", func(t *testing.T) {
		ok, msg := evalClick("tape.ctml", command{Args: []string{"banner"}, Line: 3}, frame, boxes)
		assert.True(t, ok)
		assert.Empty(t, msg)
	})

	t.Run("bad: an id absent from the box map fails, names the tape line, and shows the frame", func(t *testing.T) {
		ok, msg := evalClick("tape.ctml", command{Args: []string{"nope"}, Line: 3}, frame, boxes)
		assert.False(t, ok)
		assert.Contains(t, msg, "tape.ctml:3:")
		assert.Contains(t, msg, `Click "nope"`)
		assert.Contains(t, msg, "banner") // lists what WAS recorded
		assert.Contains(t, msg, "frame:")
		assert.Contains(t, msg, frame)
	})

	t.Run("ugly: a recorded id with a zero-area rectangle fails like an absent one", func(t *testing.T) {
		zero := html.BoxMap{"empty": {Row: 0, Col: 0, Width: 0, Height: 0}}
		ok, _ := evalClick("tape.ctml", command{Args: []string{"empty"}, Line: 1}, frame, zero)
		assert.False(t, ok)
	})
}
