// SPDX-Licence-Identifier: EUPL-1.2

// Package form is go-html's terminal-form-building API, wrapping
// charmbracelet/huh (MIT, actively maintained — published today at the
// charm.land/huh/v2 vanity path). A consumer builds/runs/reads a form —
// a Form of Groups of Fields — through this seam and never imports charm.land
// or charmbracelet directly.
//
// huh is a finished, self-driving product (Form.Run spins up its own
// bubbletea v2 program internally) rather than a widget wired through
// go-html's own tui.Program, so this package is a WRAP (re-export): charm
// stays the implementation, and go-html only swaps the import path if huh
// ever moves.
//
// Usage example:
//
//	var name string
//	f := form.NewForm(
//		form.NewGroup(
//			form.NewInput().Title("Name").Value(&name),
//			form.NewSelect[string]().Title("Colour").
//				Options(form.NewOptions("red", "green", "blue")...),
//		),
//	)
//	if err := f.Run(); err != nil {
//		// form.ErrUserAborted means the user cancelled; anything else is real.
//	}
package form

import "charm.land/huh/v2"

// Form-building types. A Form holds Groups (pages shown one at a time); a
// Group holds Fields (the inputs on that page). Theme, KeyMap and
// FieldPosition are the customisation seams; FormState and the Err* sentinels
// below are how a consumer reads what Form.Run returned.
type (
	Form          = huh.Form
	Group         = huh.Group
	Field         = huh.Field
	Theme         = huh.Theme
	ThemeFunc     = huh.ThemeFunc
	KeyMap        = huh.KeyMap
	FieldPosition = huh.FieldPosition
	FormState     = huh.FormState
)

// Field types a consumer names directly when building a Group.
type (
	Input      = huh.Input
	Text       = huh.Text
	Confirm    = huh.Confirm
	Note       = huh.Note
	FilePicker = huh.FilePicker
	EchoMode   = huh.EchoMode
)

// Select, MultiSelect and Option are generic. huh.NewSelect[T] and friends
// can't be var-forwarded — a generic func has no value until it's
// instantiated — so these are generic type aliases instead; the matching
// generic constructor funcs follow below.
type (
	Select[T comparable]      = huh.Select[T]
	MultiSelect[T comparable] = huh.MultiSelect[T]
	Option[T comparable]      = huh.Option[T]
)

// FormState identities, readable from Form.State once Run returns.
const (
	StateNormal    = huh.StateNormal
	StateCompleted = huh.StateCompleted
	StateAborted   = huh.StateAborted
)

// EchoMode identities for Input.EchoMode: plain text, password masking, or no
// echo at all.
const (
	EchoModeNormal   = huh.EchoModeNormal
	EchoModePassword = huh.EchoModePassword
	EchoModeNone     = huh.EchoModeNone
)

// Sentinel errors Form.Run / Form.RunWithContext can return — check with
// errors.Is.
var (
	ErrUserAborted        = huh.ErrUserAborted
	ErrTimeout            = huh.ErrTimeout
	ErrTimeoutUnsupported = huh.ErrTimeoutUnsupported
)

// Non-generic constructors.
var (
	NewForm       = huh.NewForm
	NewGroup      = huh.NewGroup
	NewInput      = huh.NewInput
	NewText       = huh.NewText
	NewConfirm    = huh.NewConfirm
	NewNote       = huh.NewNote
	NewFilePicker = huh.NewFilePicker

	// NewDefaultKeyMap returns the default KeyMap, the only way to obtain one
	// to customise: its bindings are charm.land/bubbles/v2/key.Binding values,
	// not wrapped here, but their SetKeys/SetEnabled/SetHelp methods (called
	// on the fields of the returned KeyMap, e.g. km.Input.Next.SetKeys(...))
	// cover the common rebind/disable cases without ever naming that package.
	NewDefaultKeyMap = huh.NewDefaultKeyMap
)

// Generic constructors. huh.NewSelect[T], NewMultiSelect[T], NewOption[T] and
// NewOptions[T] can't be var-forwarded for the same reason the types can't be
// plain aliases of a concrete instantiation — each gets its own generic
// wrapper that instantiates and calls straight through.
func NewSelect[T comparable]() *Select[T]                   { return huh.NewSelect[T]() }
func NewMultiSelect[T comparable]() *MultiSelect[T]         { return huh.NewMultiSelect[T]() }
func NewOption[T comparable](key string, value T) Option[T] { return huh.NewOption(key, value) }
func NewOptions[T comparable](values ...T) []Option[T]      { return huh.NewOptions(values...) }

// Themes. Pass one to Form.WithTheme wrapped in ThemeFunc, e.g.
// form.NewForm(...).WithTheme(form.ThemeFunc(form.ThemeDracula)).
var (
	ThemeCharm      = huh.ThemeCharm
	ThemeBase       = huh.ThemeBase
	ThemeBase16     = huh.ThemeBase16
	ThemeDracula    = huh.ThemeDracula
	ThemeCatppuccin = huh.ThemeCatppuccin
)

// Validation helpers for Field.Validate.
var (
	ValidateNotEmpty  = huh.ValidateNotEmpty
	ValidateMinLength = huh.ValidateMinLength
	ValidateMaxLength = huh.ValidateMaxLength
	ValidateLength    = huh.ValidateLength
	ValidateOneOf     = huh.ValidateOneOf
)
