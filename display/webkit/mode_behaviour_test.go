package webkit

import core "dappco.re/go"

// applyModeDefaults fills Mac / Windows platform fields per the chosen Mode,
// only touching zero-value fields so explicit caller settings survive.
//
//	cfg := &GuiConfig{Mode: ModeTray}
//	applyModeDefaults(cfg) // cfg.Mac.ActivationPolicy == ActivationPolicyAccessory
func TestModeBehaviour_applyModeDefaults_Tray(t *core.T) {
	cfg := &GuiConfig{Mode: ModeTray}
	applyModeDefaults(cfg)
	core.AssertEqual(t, ActivationPolicyAccessory, cfg.Mac.ActivationPolicy)
	core.AssertTrue(t, cfg.Windows.DisableQuitOnLastWindowClosed)
}

// ModeMultiWindow keeps the Regular activation policy but disables
// quit-on-last-window so the process survives the last close.
func TestModeBehaviour_applyModeDefaults_MultiWindow(t *core.T) {
	cfg := &GuiConfig{Mode: ModeMultiWindow}
	applyModeDefaults(cfg)
	core.AssertEqual(t, ActivationPolicyRegular, cfg.Mac.ActivationPolicy)
	core.AssertTrue(t, cfg.Windows.DisableQuitOnLastWindowClosed)
}

// ModeSingleWindow sets terminate-after-last-window and leaves Windows
// quit-on-last-window enabled (zero value).
func TestModeBehaviour_applyModeDefaults_SingleWindow(t *core.T) {
	cfg := &GuiConfig{Mode: ModeSingleWindow}
	applyModeDefaults(cfg)
	core.AssertTrue(t, cfg.Mac.ApplicationShouldTerminateAfterLastWindowClosed)
	core.AssertFalse(t, cfg.Windows.DisableQuitOnLastWindowClosed)
}

// applyModeDefaults preserves an explicitly-set activation policy in tray mode
// (only the Regular zero value is overridden).
func TestModeBehaviour_applyModeDefaults_PreservesExplicit(t *core.T) {
	cfg := &GuiConfig{Mode: ModeTray}
	cfg.Mac.ActivationPolicy = ActivationPolicyProhibited
	applyModeDefaults(cfg)
	core.AssertEqual(t, ActivationPolicyProhibited, cfg.Mac.ActivationPolicy)
}

// applyModeDefaults is a safe no-op for ModeDefault and a nil config.
func TestModeBehaviour_applyModeDefaults_Ugly(t *core.T) {
	cfg := &GuiConfig{Mode: ModeDefault}
	applyModeDefaults(cfg)
	core.AssertEqual(t, ActivationPolicyRegular, cfg.Mac.ActivationPolicy)
	core.AssertFalse(t, cfg.Windows.DisableQuitOnLastWindowClosed)

	core.AssertNotPanics(t, func() { applyModeDefaults(nil) })
}
