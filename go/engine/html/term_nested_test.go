//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The HLCRF frame is a box model: a Layout is a Node, so a frame nests inside
// another frame's slot and renders width-bounded within the parent column —
// the same clone-on-render nesting the HTML compositor performs.
func TestTermLayout_RenderTerm_Nested(t *testing.T) {
	inner := NewLayout("HCF").
		H(Text("in.h")).
		C(El("p", Text("in.c"))).
		F(Text("in.f"))
	outer := NewLayout("HLCF").
		H(Text("out.h")).
		L(El("ul", El("li", Text("nav")))).
		C(El("h2", Text("out.title")), inner).
		F(Text("out.f"))
	ctx := termTestContext(map[string]string{
		"out.h": "Outer header", "nav": "Menu", "out.title": "Outer content",
		"in.h": "Inner header", "in.c": "Inner body", "in.f": "Inner footer",
		"out.f": "Outer footer",
	})

	out := outer.RenderTerm(ctx, TermOptions{Width: 92})

	for _, want := range []string{
		"Outer header", "Menu", "Outer content",
		"Inner header", "Inner body", "Inner footer", "Outer footer",
	} {
		assert.Contains(t, out, want)
	}

	lines := strings.Split(out, "\n")
	innerLine := ""
	for _, line := range lines {
		if strings.Contains(line, "Inner body") {
			innerLine = line
			break
		}
	}
	require.NotEmpty(t, innerLine, "inner content line present")
	assert.True(t, strings.HasPrefix(strings.TrimRight(innerLine, " "), strings.Repeat(" ", 20)),
		"inner frame renders inside the outer content column, not at column zero")

	assert.Less(t, strings.Index(out, "Outer header"), strings.Index(out, "Inner header"),
		"outer frame opens before the nested frame")
	assert.Less(t, strings.Index(out, "Inner footer"), strings.Index(out, "Outer footer"),
		"nested frame closes before the outer footer band")
}
