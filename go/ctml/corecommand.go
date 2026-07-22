// SPDX-Licence-Identifier: EUPL-1.2

package ctml

import (
	"strings"

	core "dappco.re/go"
	html "dappco.re/go/html"
)

// SubcommandList generates a <ul> of a CoreCommand path's direct
// subcommands -- the well-founded half of the exploratory
// CoreCommand-derived default TUI; see docs/ctml.md S:S15 for why the
// other half (flags -> input chips) is a design write-up, not code, in
// this pass. A Hidden command is skipped, matching Command.Hidden's own
// meaning elsewhere in dappco.re/go.
//
// paths is every registered command path (typically core.Core.Commands());
// root selects which level to list -- "" for the top level, "deploy" for
// the direct children of "deploy". A command's label is its I18nKey(),
// the same i18n-key-as-text convention S:S6.1 uses for ordinary .ctml
// text -- Command already carries a real label mechanism, even though
// Option (a flag) does not.
//
// Usage example: list := ctml.SubcommandList(c, "", c.Commands())
func SubcommandList(c *core.Core, root string, paths []string) html.Node {
	var items []html.Node
	for _, path := range directChildPaths(root, paths) {
		cmd, ok := c.Command(path).Value.(*core.Command)
		if !ok || cmd == nil || cmd.Hidden {
			continue
		}
		items = append(items, html.El("li", html.Text(cmd.I18nKey())))
	}
	return html.El("ul", items...)
}

// directChildPaths filters paths to those exactly one path segment below
// root, deduplicated and in first-seen order.
func directChildPaths(root string, paths []string) []string {
	prefix := ""
	if root != "" {
		prefix = root + "/"
	}

	seen := make(map[string]bool)
	var out []string
	for _, p := range paths {
		rest := p
		if root == "" {
			if strings.Contains(p, "/") {
				continue
			}
		} else {
			if !strings.HasPrefix(p, prefix) {
				continue
			}
			rest = p[len(prefix):]
			if strings.Contains(rest, "/") {
				continue
			}
		}
		if seen[p] {
			continue
		}
		seen[p] = true
		out = append(out, p)
	}
	return out
}
