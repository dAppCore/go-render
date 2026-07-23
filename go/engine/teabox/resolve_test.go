// SPDX-Licence-Identifier: EUPL-1.2

package teabox

import (
	"testing"

	core "dappco.re/go"
	html "dappco.re/go/render/engine/html"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func syntheticBoxes() (html.BoxMap, html.Node, html.Node, html.Node) {
	page := html.El("div")
	nav := html.El("nav")
	card := html.El("div")
	return html.BoxMap{
		"H": {Row: 0, Col: 0, Width: 40, Height: 2, Node: page},
		"L": {Row: 2, Col: 0, Width: 10, Height: 8, Node: nav},
		"C": {Row: 2, Col: 10, Width: 30, Height: 8, Node: card},
	}, page, nav, card
}

func TestResolve_Good(t *testing.T) {
	boxes, _, nav, card := syntheticBoxes()

	tests := []struct {
		name    string
		x, y    int
		wantID  string
		wantHit bool
	}{
		{"good: click inside H", 5, 0, "H", true},
		{"good: click inside L", 3, 5, "L", true},
		{"good: click inside C", 25, 5, "C", true},
		{"good: top-left corner of a box is inside it", 0, 0, "H", true},
		{"good: bottom-right corner is exclusive (one past the box is outside)", 40, 2, "", false},
		{"good: last inclusive row/col of H is inside it", 39, 1, "H", true},
		{"bad: click outside every box", 5, 100, "", false},
		{"bad: negative coordinates are outside every box", -1, 0, "", false},
		{"ugly: click exactly on the boundary between L and C belongs to C, not L", 10, 5, "C", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hit, ok := Resolve(boxes, tc.x, tc.y)
			require.Equal(t, tc.wantHit, ok)
			if tc.wantHit {
				assert.Equal(t, tc.wantID, hit.BlockID)
			}
		})
	}

	navHit, ok := Resolve(boxes, 3, 5)
	require.True(t, ok)
	assert.Same(t, nav, navHit.Box.Node)

	cardHit, ok := Resolve(boxes, 25, 5)
	require.True(t, ok)
	assert.Same(t, card, cardHit.Box.Node)
}

func TestResolve_OverlappingBoxesPickSmallestArea(t *testing.T) {
	inner := html.El("button")
	outer := html.El("div")
	boxes := html.BoxMap{
		"C":       {Row: 0, Col: 0, Width: 40, Height: 20, Node: outer},
		"C.card1": {Row: 2, Col: 2, Width: 10, Height: 5, Node: inner},
	}

	hit, ok := Resolve(boxes, 5, 3)
	require.True(t, ok)
	assert.Equal(t, "C.card1", hit.BlockID, "the smaller, nested box wins over the larger enclosing one")
	assert.Same(t, inner, hit.Box.Node)

	hit, ok = Resolve(boxes, 30, 15)
	require.True(t, ok)
	assert.Equal(t, "C", hit.BlockID, "outside the nested box, the enclosing box still resolves")
}

func TestResolve_TieBreaksOnBlockID(t *testing.T) {
	a := html.El("a")
	b := html.El("b")
	boxes := html.BoxMap{
		"zzz": {Row: 0, Col: 0, Width: 5, Height: 5, Node: a},
		"aaa": {Row: 0, Col: 0, Width: 5, Height: 5, Node: b},
	}

	hit, ok := Resolve(boxes, 2, 2)
	require.True(t, ok)
	assert.Equal(t, "aaa", hit.BlockID, "an exact-area tie resolves deterministically to the lexicographically smaller ID")
}

func TestResolve_Ugly(t *testing.T) {
	t.Run("ugly: empty box map never matches", func(t *testing.T) {
		_, ok := Resolve(html.BoxMap{}, 0, 0)
		assert.False(t, ok)
	})

	t.Run("ugly: nil box map never matches", func(t *testing.T) {
		_, ok := Resolve(nil, 0, 0)
		assert.False(t, ok)
	})

	t.Run("ugly: zero-size box never matches, even at its own origin", func(t *testing.T) {
		boxes := html.BoxMap{"z": {Row: 0, Col: 0, Width: 0, Height: 0}}
		_, ok := Resolve(boxes, 0, 0)
		assert.False(t, ok)
	})

	t.Run("ugly: negative width/height never matches", func(t *testing.T) {
		boxes := html.BoxMap{"z": {Row: 0, Col: 0, Width: -5, Height: -5}}
		_, ok := Resolve(boxes, 0, 0)
		assert.False(t, ok)
	})
}

func TestResolveNode_Good(t *core.T) {
	boxes, _, nav, _ := syntheticBoxes()
	node, ok := ResolveNode(boxes, 3, 5)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, nav, node)
}

func TestResolveNode_Bad(t *core.T) {
	boxes, _, _, _ := syntheticBoxes()
	node, ok := ResolveNode(boxes, 999, 999)
	core.AssertFalse(t, ok)
	core.AssertNil(t, node)
}

func TestResolveNode_Ugly(t *core.T) {
	node, ok := ResolveNode(nil, 0, 0)
	core.AssertFalse(t, ok)
	core.AssertNil(t, node)
}
