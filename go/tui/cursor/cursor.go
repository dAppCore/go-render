// SPDX-Licence-Identifier: EUPL-1.2

// Package cursor re-exports charmbracelet/bubbles/cursor through go-html's
// tui seam — the blinking text cursor textinput and textarea render inside
// themselves. Swap the import path (bubbles/cursor → html/tui/cursor) and
// keep every cursor.Model / cursor.New reference unchanged. Mode, SetMode,
// Focus, Blur, SetChar, View and Update are Model methods, not package
// functions, so they need no re-export here — they come along for free
// since Model is a genuine alias, and so do the exported Style, TextStyle,
// BlinkSpeed and IsBlinked fields. Blink is overloaded: the package-level
// Blink below is the Cmd that arms the very first blink (send it from Init,
// or take it straight off Focus/SetMode(CursorBlink), which call it for
// you), while Model's own Blink method re-arms each subsequent one and
// needs no re-export of its own. Both eventually deliver a BlinkMsg to
// Update, which is the loop that keeps a CursorBlink-mode cursor animating.
package cursor

import "charm.land/bubbles/v2/cursor"

// Model is the cursor itself. Mode selects its behaviour — CursorBlink
// animates, CursorStatic stays solid and visible, CursorHide is invisible —
// set via SetMode and read back via Model's own Mode method. BlinkMsg is the
// message a Blink Cmd (package-level or Model.Blink) delivers to Update to
// keep a CursorBlink-mode cursor animating.
type (
	Model    = cursor.Model
	Mode     = cursor.Mode
	BlinkMsg = cursor.BlinkMsg
)

// Cursor mode identities for SetMode, matching what Model's own Mode method
// returns.
const (
	CursorBlink  = cursor.CursorBlink
	CursorStatic = cursor.CursorStatic
	CursorHide   = cursor.CursorHide
)

var (
	New = cursor.New

	// Blink is the Cmd that arms a cursor's first blink — return it from
	// Init to start the loop. Focus and SetMode(CursorBlink) already call
	// it for you, so most consumers never name it directly.
	Blink = cursor.Blink
)
