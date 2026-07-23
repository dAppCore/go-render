// SPDX-Licence-Identifier: EUPL-1.2

package image_test

import (
	"image"
	"image/color"
	"strings"
	"testing"

	tuiimage "dappco.re/go/html/tui/image"
)

// newCheckerboard builds a deterministic n-by-n RGBA image — a two-colour
// checkerboard — so every test run renders byte-identical pixels: no
// external image files, no timestamps.
func newCheckerboard(n int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, n, n))
	for y := range n {
		for x := range n {
			c := color.RGBA{R: 220, G: 40, B: 40, A: 255}
			if (x+y)%2 == 0 {
				c = color.RGBA{R: 40, G: 60, B: 220, A: 255}
			}
			img.Set(x, y, c)
		}
	}
	return img
}

// wantLines is the number of terminal rows Render emits for a given output
// Height. mosaic walks the (rescaled) image in 2-pixel-tall steps, one
// character per step, so a requested Height of pixel rows halves into that
// many terminal lines, rounded up — not Height itself.
func wantLines(height int) int {
	return (height + 1) / 2
}

// TestNew builds a Mosaic through New and configures it fluently — proving
// the Mosaic alias carries Width and Height, and that Render (a
// pointer-receiver method) works when called on the named variable m. A
// fully chained New().Width(w).Height(h).Render(img) would not compile —
// the chain's result isn't addressable — which is why m is assigned first.
func TestNew(t *testing.T) {
	const width, height = 8, 4

	m := tuiimage.New().Width(width).Height(height)
	out := m.Render(newCheckerboard(16))

	if out == "" {
		t.Fatal("Render() returned empty output")
	}
	got := len(strings.Split(strings.TrimRight(out, "\n"), "\n"))
	if want := wantLines(height); got != want {
		t.Fatalf("Render() produced %d lines, want %d (Height %d halved by half-block rendering)", got, want, height)
	}
}

// TestRender proves the one-shot Render function renders output identical
// to the fluent New().Width().Height().Render(img) path it wraps, and that
// its line count follows the same Height-halving rule as TestNew.
func TestRender(t *testing.T) {
	const width, height = 6, 5

	img := newCheckerboard(12)
	out := tuiimage.Render(img, width, height)

	if out == "" {
		t.Fatal("Render() returned empty output")
	}
	got := len(strings.Split(strings.TrimRight(out, "\n"), "\n"))
	if want := wantLines(height); got != want {
		t.Fatalf("Render() produced %d lines, want %d (Height %d halved by half-block rendering)", got, want, height)
	}

	m := tuiimage.New().Width(width).Height(height)
	if fluent := m.Render(img); out != fluent {
		t.Fatalf("Render(img, %d, %d) = %q, want %q (must match the fluent New().Width().Height().Render(img) path)", width, height, out, fluent)
	}
}
