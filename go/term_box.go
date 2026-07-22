//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	"strconv"
	"strings"
)

// term_box.go: the terminal box map. RenderTermBoxes renders exactly like
// RenderTerm, but also records the rendered rectangle of every identified
// block, so a caller can resolve a mouse event's (x, y) back to the node
// that occupies it -- see go/teabox for the resolver.
//
// A block is identified two ways, matching the two places blocks already
// have addresses in this codebase: a Layout slot (H/L/C/R/F, including
// nested layouts) is keyed by the same blockID string the HTML renderer
// already uses for its data-block attribute; any element carrying an
// explicit id attribute is keyed by that id. Recording is opt-in and
// adds no cost to plain RenderTerm/Layout.RenderTerm/Responsive.RenderTerm:
// the recorder is nil on those paths, and every new code path below is
// guarded on it being non-nil.

// Box is the rendered terminal rectangle for one identified block.
// Row/Col are 0-based screen coordinates of the top-left cell.
type Box struct {
	Row, Col      int
	Width, Height int
	Node          Node
}

// BoxMap maps a block ID to the rectangle it rendered at.
type BoxMap map[string]Box

// termBoxRecorder accumulates boxes during one render. originRow/originCol
// track the absolute position the renderer is currently descending from,
// so nested calls (a Layout inside a Layout's C slot, an id'd element
// inside a slot) record absolute, page-relative coordinates. Element-id
// boxes recorded inside a themed slot (Header/Sidebar/Aside/Footer/Card)
// use that slot's own outer origin as an approximation -- the slot's
// border/padding can offset true content position by a cell or two; the
// slot's own box (always exact) is the precise target for the slot itself.
type termBoxRecorder struct {
	boxes                BoxMap
	originRow, originCol int
	frame                int // 0 = the outermost layout; incremented per renderTermFrame call
}

func (rec *termBoxRecorder) record(id string, row, col, width, height int, n Node) {
	if rec == nil || id == "" || height <= 0 || width <= 0 {
		return
	}
	rec.boxes[id] = Box{Row: row, Col: col, Width: width, Height: height, Node: n}
}

// framePrefix returns the key prefix for the layout frame currently being
// entered, and advances the frame counter. The outermost layout (the
// common case) keeps the clean "H"/"L"/"C"/"R"/"F" keys that already
// match the HTML renderer's data-block attribute; a nested layout (one
// Layout inside another's slot) gets an "L<n>." prefix so its slots
// cannot collide with the outer layout's -- term rendering never threads
// the HTML side's clone-on-render path IDs (Layout.path is always ""
// here), so this is box-map-local disambiguation, not a reuse of that
// scheme.
func (r *termRenderer) framePrefix() string {
	if r.rec == nil {
		return ""
	}
	n := r.rec.frame
	r.rec.frame++
	if n == 0 {
		return ""
	}
	return "L" + strconv.Itoa(n) + "."
}

// originRow and originCol read the recorder's current origin, or (0, 0)
// when boxes are not being recorded -- every caller can use them
// unconditionally instead of nil-checking r.rec itself.
func (r *termRenderer) originRow() int {
	if r.rec == nil {
		return 0
	}
	return r.rec.originRow
}

func (r *termRenderer) originCol() int {
	if r.rec == nil {
		return 0
	}
	return r.rec.originCol
}

// recordElBox records the box for a block-level *elNode carrying an id
// attribute, keyed by that id -- see the recorder-origin note above for
// the approximation this makes inside a padded/bordered slot.
func (r *termRenderer) recordElBox(n Node, startRow, width int, lines []string) {
	if r.rec == nil {
		return
	}
	el, ok := n.(*elNode)
	if !ok || el == nil {
		return
	}
	id := el.attrs["id"]
	if id == "" {
		return
	}
	height := len(lines)
	for height > 0 && lines[height-1] == "" {
		height--
	}
	r.rec.record(id, r.originRow()+startRow, r.originCol(), width, height, n)
}

// withOrigin temporarily moves the recorder's origin for the duration of
// fn, restoring it afterwards. Rendering is single-threaded and strictly
// sequential, so a plain save/restore is safe -- no recorder means no-op.
func (r *termRenderer) withOrigin(row, col int, fn func()) {
	if r.rec == nil {
		fn()
		return
	}
	savedRow, savedCol := r.rec.originRow, r.rec.originCol
	r.rec.originRow, r.rec.originCol = row, col
	fn()
	r.rec.originRow, r.rec.originCol = savedRow, savedCol
}

// termLineCount returns how many terminal rows s occupies -- 0 for an
// empty string (nothing rendered, so nothing to record), otherwise the
// newline count plus one.
func termLineCount(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n") + 1
}

// RenderTermBoxes renders n like RenderTerm and additionally returns the
// box map of every identified block's rendered bounds.
// Example: out, boxes := RenderTermBoxes(page, ctx, TermOptions{Width: 100})
func RenderTermBoxes(n Node, ctx *Context, opts ...TermOptions) (string, BoxMap) {
	boxes := BoxMap{}
	if n == nil {
		return "", boxes
	}
	width, theme, fit := resolveTermOptions(opts)
	rec := &termBoxRecorder{boxes: boxes}
	r := &termRenderer{ctx: termContext(ctx), theme: theme, fit: fit, rec: rec}
	out := strings.TrimRight(strings.Join(r.blocks([]Node{n}, width), "\n"), "\n")
	return out, boxes
}

// RenderTermBoxes is Layout.RenderTerm plus a box map of every slot's
// (and any id'd element's) rendered bounds, keyed the same way as
// package-level RenderTermBoxes.
// Example: out, boxes := layout.RenderTermBoxes(ctx, TermOptions{Width: 120})
func (l *Layout) RenderTermBoxes(ctx *Context, opts ...TermOptions) (string, BoxMap) {
	boxes := BoxMap{}
	if l == nil {
		return "", boxes
	}
	width, theme, fit := resolveTermOptions(opts)
	rec := &termBoxRecorder{boxes: boxes}
	r := &termRenderer{ctx: termContext(ctx), theme: theme, fit: fit, rec: rec}
	out := l.renderTermFrame(r, width)
	return out, boxes
}

// RenderTermBoxes is Responsive.RenderTerm plus a box map for whichever
// variant width picks -- see RenderTermBoxes.
// Example: out, boxes := resp.RenderTermBoxes(ctx, TermOptions{Width: 72})
func (resp *Responsive) RenderTermBoxes(ctx *Context, opts ...TermOptions) (string, BoxMap) {
	boxes := BoxMap{}
	if resp == nil || len(resp.variants) == 0 {
		return "", boxes
	}
	width, theme, fit := resolveTermOptions(opts)
	rec := &termBoxRecorder{boxes: boxes}
	r := &termRenderer{ctx: termContext(ctx), theme: theme, fit: fit, rec: rec}
	out := resp.renderTermPick(r, width)
	return out, boxes
}
