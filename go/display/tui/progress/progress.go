// SPDX-Licence-Identifier: EUPL-1.2

// Package progress re-exports charmbracelet/bubbles/progress through
// go-html's tui seam — an animated progress bar. Swap the import path
// (bubbles/progress → html/tui/progress) and keep every progress.Model /
// progress.New reference unchanged. ViewAs, SetPercent, IncrPercent,
// DecrPercent, Percent, Width, SetWidth and IsAnimating are Model methods,
// not package functions, so they need no re-export here — they come along
// for free since Model is a genuine alias. SetPercent/IncrPercent/DecrPercent
// return a Cmd that must flow back into Update to animate the bar; FrameMsg
// is the message that Cmd carries.
package progress

import "charm.land/bubbles/v2/progress"

// Model is the bar itself. Option configures New. ColorFunc backs
// WithColorFunc, painting the bar from the total/current percentage rather
// than a fixed colour or blend. FrameMsg is what SetPercent/IncrPercent/
// DecrPercent's Cmd delivers to Update while the bar animates towards its
// target percentage.
type (
	Model     = progress.Model
	Option    = progress.Option
	ColorFunc = progress.ColorFunc
	FrameMsg  = progress.FrameMsg
)

// Fill-rune identities for WithFillCharacters (or set directly on
// Model.Full / Model.Empty).
const (
	DefaultFullCharHalfBlock = progress.DefaultFullCharHalfBlock
	DefaultFullCharFullBlock = progress.DefaultFullCharFullBlock
	DefaultEmptyCharBlock    = progress.DefaultEmptyCharBlock
)

var (
	New = progress.New

	WithDefaultBlend   = progress.WithDefaultBlend
	WithColors         = progress.WithColors
	WithColorFunc      = progress.WithColorFunc
	WithFillCharacters = progress.WithFillCharacters
	WithoutPercentage  = progress.WithoutPercentage
	WithWidth          = progress.WithWidth
	WithSpringOptions  = progress.WithSpringOptions
	WithScaled         = progress.WithScaled
)
