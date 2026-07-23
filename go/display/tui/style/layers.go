package style

import "charm.land/lipgloss/v2"

// Layer is a positioned block of content within a z-ordered composition —
// the primitive behind floating overlays, modals and panels: build one per
// visual element, place it with X/Y, and stack several with Z (a higher Z
// paints over a lower one at the same position). A Layer can nest children
// (AddLayers) that are positioned relative to it, e.g.
//
//	background := style.NewLayer(page)
//	modal := style.NewLayer(dialog).X(4).Y(2).Z(1) // painted above the page
//	out := style.NewCompositor(background, modal).Render()
type Layer = lipgloss.Layer

// NewLayer builds a Layer from content and optional child layers. Position
// it with X/Y/Z (see Layer) before compositing.
func NewLayer(content string, layers ...*Layer) *Layer {
	return lipgloss.NewLayer(content, layers...)
}

// LayerHit is the result of a Compositor.Hit lookup — the topmost identified
// Layer at a point, or Empty() if the point hit nothing.
type LayerHit = lipgloss.LayerHit

// Compositor flattens a Layer hierarchy once, then renders or hit-tests it —
// build the base view and any overlay as Layers, hand both to NewCompositor,
// and Render composites them in z-order into a single string.
type Compositor = lipgloss.Compositor

// NewCompositor builds a Compositor from the given root-level Layers.
func NewCompositor(layers ...*Layer) *Compositor {
	return lipgloss.NewCompositor(layers...)
}

// Canvas is a cell buffer that Layers (or anything else drawable) can be
// composed onto directly — later composes paint over earlier ones — then
// rendered to a styled string with Render. Compositor.Render uses a Canvas
// internally; reach for Canvas directly to compose onto a fixed-size buffer
// rather than a self-sizing Layer tree.
type Canvas = lipgloss.Canvas

// NewCanvas builds an empty Canvas of the given width and height, in cells.
func NewCanvas(width, height int) *Canvas {
	return lipgloss.NewCanvas(width, height)
}
