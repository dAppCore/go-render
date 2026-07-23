// SPDX-Licence-Identifier: EUPL-1.2

// Package ctml parses .ctml documents into the same node trees the
// dappco.re/go/render builder API produces by hand -- see docs/ctml.md for
// the full grammar. HTML render, terminal render, entitlements, and i18n
// all work unchanged on a parsed tree, because it is the same tree.
//
//	tree, err := ctml.Parse(src)
//	if err != nil { ... }
//	out := html.Render(tree, html.NewContext())
package ctml
