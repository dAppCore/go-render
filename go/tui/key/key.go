// Package key re-exports charmbracelet/bubbles/key through go-html's tui seam —
// key bindings and matching. Swap the import path (bubbles/key →
// html/tui/key) and keep every key.Binding / key.Matches reference unchanged.
package key

import (
	"fmt"

	"charm.land/bubbles/v2/key"
)

type (
	Binding = key.Binding
	Help    = key.Help
)

var (
	NewBinding = key.NewBinding
	WithHelp   = key.WithHelp
	WithKeys   = key.WithKeys
)

// Matches reports whether k (any fmt.Stringer key message) matches any of the
// given bindings. A generic forward of key.Matches, so the call site is
// unchanged.
func Matches[Key fmt.Stringer](k Key, b ...Binding) bool { return key.Matches(k, b...) }
