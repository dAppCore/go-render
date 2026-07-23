// SPDX-Licence-Identifier: EUPL-1.2

// Package image re-exports charmbracelet/x/mosaic through go-html's tui
// seam — a Unicode image renderer that turns a Go image.Image into a
// terminal string, pixels-per-cell (the surface the lem tui uses to show
// vision output inline — gemma4 is multimodal). Swap the import path
// (x/mosaic → html/tui/image) and keep every Mosaic / New / Render
// reference unchanged. mosaic is lipgloss-free — ansi and golang.org/x/image
// only — so this package never drags lipgloss along.
//
// Width, Height and the other builder methods (Scale, Dither, Threshold,
// InvertColors, IgnoreBlockSymbols, Symbol) are Mosaic methods, not package
// functions, so they need no re-export here — they come along for free
// since Mosaic is a genuine alias. Render, though, has a pointer receiver:
// build into a named variable before calling it, since a fully chained
// New().Width(w).Height(h).Render(img) will not compile (the chain's result
// isn't addressable). Render (the package function below) is the one-shot
// form that sidesteps this for the common case.
//
// Usage example:
//
//	m := image.New().Width(80).Height(40)
//	out := m.Render(img)
//
//	// or, one-shot:
//	out := image.Render(img, 80, 40)
package image

import "github.com/charmbracelet/x/mosaic"

// Mosaic renders a Go image.Image to a terminal string, pixels-per-cell.
// Build one with New, configure it fluently (Width, Height, Scale, Dither,
// Threshold, InvertColors, IgnoreBlockSymbols, Symbol), then call Render.
type Mosaic = mosaic.Mosaic

var (
	// New returns a zero-configured Mosaic, ready for Width/Height/Render.
	New = mosaic.New

	// Render renders img to a terminal string sized width-by-height cells
	// in one call — the one-shot form of
	// New().Width(width).Height(height).Render(img) for a consumer that
	// doesn't need the other Mosaic options.
	Render = mosaic.Render
)
