// SPDX-Licence-Identifier: EUPL-1.2

package main

import html "dappco.re/go/render/engine/html"

// renderLayout renders an HLCRF layout from a slot map.
//
// Empty string values are meaningful: they create an explicit empty slot
// container rather than being treated as absent input.
func renderLayout(variant, locale string, slots map[string]string) string {
	if variant == "" {
		return ""
	}

	ctx := html.NewContext()
	if locale != "" {
		ctx.SetLocale(locale)
	}

	layout := html.NewLayout(variant)

	for _, slot := range []string{"H", "L", "C", "R", "F"} {
		content, ok := slots[slot]
		if !ok {
			continue
		}

		switch slot {
		case "H":
			layout.H(html.Raw(content))
		case "L":
			layout.L(html.Raw(content))
		case "C":
			layout.C(html.Raw(content))
		case "R":
			layout.R(html.Raw(content))
		case "F":
			layout.F(html.Raw(content))
		}
	}

	return layout.Render(ctx)
}
