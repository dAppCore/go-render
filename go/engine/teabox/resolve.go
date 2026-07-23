// SPDX-Licence-Identifier: EUPL-1.2

// Package teabox resolves a terminal (x, y) coordinate to the go-html node
// occupying it, using the box map html.RenderTermBoxes produces. It has no
// dependency on bubbletea: a caller feeds plain ints, so the dependency is
// only ever pulled in by callers that already need the type --
//
//	out, boxes := html.RenderTermBoxes(page, ctx, html.TermOptions{Width: 100})
//	// inside a bubbletea Update:
//	case tea.MouseMsg:
//	    m := msg.Mouse()
//	    if hit, ok := teabox.Resolve(boxes, m.X, m.Y); ok {
//	        return m.handle(hit.BlockID, hit.Box.Node)
//	    }
package teabox

import html "dappco.re/go/render/engine/html"

// Hit is one resolved box map entry.
type Hit struct {
	BlockID string
	Box     html.Box
}

// Resolve finds the box in boxes that contains the 0-based screen
// coordinate (x, y) -- matching html.Box's Row/Col/Width/Height. ok is
// false when no box contains the point. Boxes can overlap (a nested
// layout's slot renders inside its enclosing slot's rectangle); the
// smallest-area match wins, so a click inside a nested card resolves to
// the card, not the whole page. A tie on area resolves to the
// lexicographically smaller block ID, so the result is deterministic.
func Resolve(boxes html.BoxMap, x, y int) (Hit, bool) {
	var best Hit
	found := false
	bestArea := 0

	for id, box := range boxes {
		if !contains(box, x, y) {
			continue
		}
		area := box.Width * box.Height
		if !found || area < bestArea || (area == bestArea && id < best.BlockID) {
			best = Hit{BlockID: id, Box: box}
			bestArea = area
			found = true
		}
	}
	return best, found
}

// ResolveNode is Resolve, returning just the node -- the common case for
// wiring a click handler that does not need the block ID or rectangle.
func ResolveNode(boxes html.BoxMap, x, y int) (html.Node, bool) {
	hit, ok := Resolve(boxes, x, y)
	if !ok {
		return nil, false
	}
	return hit.Box.Node, true
}

func contains(box html.Box, x, y int) bool {
	if box.Width <= 0 || box.Height <= 0 {
		return false
	}
	return x >= box.Col && x < box.Col+box.Width && y >= box.Row && y < box.Row+box.Height
}
