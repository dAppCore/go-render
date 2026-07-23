// SPDX-License-Identifier: EUPL-1.2

package webkit

// Mode is the meta-decision that drives GuiConfig defaults. Apps pick
// the boot shape (Tray-only, SingleWindow, or MultiWindow) and
// webkit.Service auto-applies the corresponding platform-policy defaults
// (Mac.ActivationPolicy, Windows.DisableQuitOnLastWindowClosed, etc.)
// before the consumer's own field assignments override them.
//
// Mode is purely a defaults provider. Explicit values in
// GuiConfig.Mac / GuiConfig.Windows always win — Mode never overrides
// a non-zero field the caller set.
type Mode int

const (
	// ModeDefault leaves all platform fields at their zero values.
	// Equivalent to setting no Mode at all; useful as a sentinel for
	// callers that want to opt out of auto-defaults explicitly.
	ModeDefault Mode = iota

	// ModeTray is the menu-bar / system-tray app shape. Defaults:
	//
	//   Mac.ApplicationShouldTerminateAfterLastWindowClosed = false
	//   Mac.ActivationPolicy                                 = Accessory
	//   Windows.DisableQuitOnLastWindowClosed                = true
	//
	// The tray IS the process lifetime anchor. Closing every window
	// hides them; the process survives until the user picks Quit from
	// the tray menu. No Dock icon (macOS) / no taskbar entry (Win).
	ModeTray

	// ModeSingleWindow is the traditional one-window app shape.
	// Defaults:
	//
	//   Mac.ApplicationShouldTerminateAfterLastWindowClosed = true
	//   Mac.ActivationPolicy                                 = Regular
	//   Windows.DisableQuitOnLastWindowClosed                = false
	//
	// Closing the window quits the process. Standard Dock icon,
	// Cmd+Tab presence.
	ModeSingleWindow

	// ModeMultiWindow is the multi-window app shape (editor, IDE,
	// browser). Defaults:
	//
	//   Mac.ApplicationShouldTerminateAfterLastWindowClosed = false
	//   Mac.ActivationPolicy                                 = Regular
	//   Windows.DisableQuitOnLastWindowClosed                = false
	//
	// Process survives the last window closing (macOS pattern — File >
	// New Window re-opens). Standard Dock icon. WindowStatePath /
	// WindowLayoutPath are most useful in this mode.
	ModeMultiWindow
)

// applyModeDefaults fills in the Mac / Windows fields per the chosen
// Mode. Only fills zero-value fields — explicit caller settings are
// preserved. Called by webkit.Service.start() BEFORE buildWailsOptions
// reads the config.
func applyModeDefaults(cfg *GuiConfig) {
	if cfg == nil || cfg.Mode == ModeDefault {
		return
	}
	switch cfg.Mode {
	case ModeTray:
		// Tray = process survives last window; no Dock icon.
		// Mac defaults: terminate-after-last-window=false (already
		// the zero value, so nothing to set), ActivationPolicy=Accessory.
		if cfg.Mac.ActivationPolicy == ActivationPolicyRegular {
			// Regular IS the zero value of the int enum — only override
			// when the caller hasn't picked something else.
			cfg.Mac.ActivationPolicy = ActivationPolicyAccessory
		}
		// Windows: DisableQuitOnLastWindowClosed=true. The bool zero
		// is false; flip to true unless the caller already set it.
		// (Can't distinguish "explicit false" from "unset" on a bool,
		// so this is best-effort — apps wanting explicit quit-on-last
		// in tray mode would override.)
		if !cfg.Windows.DisableQuitOnLastWindowClosed {
			cfg.Windows.DisableQuitOnLastWindowClosed = true
		}

	case ModeMultiWindow:
		// MultiWindow = process survives last window; standard Dock.
		// Mac defaults: terminate-after-last-window=false (zero value),
		// ActivationPolicy=Regular (zero value). Nothing to set.
		// Windows defaults: DisableQuitOnLastWindowClosed=true so the
		// process survives last close (apps re-open windows from menu).
		if !cfg.Windows.DisableQuitOnLastWindowClosed {
			cfg.Windows.DisableQuitOnLastWindowClosed = true
		}

	case ModeSingleWindow:
		// SingleWindow = process quits with last window; standard Dock.
		// Mac defaults: terminate-after-last-window=true. The wails
		// default is true (the zero of a bool we DON'T set on the wails
		// side is "use platform default"), but consumers passing through
		// our GuiConfig get the explicit signal.
		cfg.Mac.ApplicationShouldTerminateAfterLastWindowClosed = true
		// Windows: DisableQuitOnLastWindowClosed stays false (zero).
	}
}
