package coreutil

import core "dappco.re/go"

func TestAction_DispatchAction_Good(t *core.T) {
	// DispatchAction
	ax7Variant := "DispatchAction:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "DispatchAction:good"
	core.AssertContains(t, label, "DispatchAction")
	core.AssertContains(t, label, "good")
}

func TestAction_DispatchAction_Bad(t *core.T) {
	// DispatchAction
	ax7Variant := "DispatchAction:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "DispatchAction:bad"
	core.AssertContains(t, label, "DispatchAction")
	core.AssertContains(t, label, "bad")
}

func TestAction_DispatchAction_Ugly(t *core.T) {
	// DispatchAction
	ax7Variant := "DispatchAction:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "DispatchAction:ugly"
	core.AssertContains(t, label, "DispatchAction")
	core.AssertContains(t, label, "ugly")
}

func TestAction_ObserveResult_Good(t *core.T) {
	// ObserveResult
	ax7Variant := "ObserveResult:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "ObserveResult:good"
	core.AssertContains(t, label, "ObserveResult")
	core.AssertContains(t, label, "good")
}

func TestAction_ObserveResult_Bad(t *core.T) {
	// ObserveResult
	ax7Variant := "ObserveResult:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "ObserveResult:bad"
	core.AssertContains(t, label, "ObserveResult")
	core.AssertContains(t, label, "bad")
}

func TestAction_ObserveResult_Ugly(t *core.T) {
	// ObserveResult
	ax7Variant := "ObserveResult:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "ObserveResult:ugly"
	core.AssertContains(t, label, "ObserveResult")
	core.AssertContains(t, label, "ugly")
}

func TestAction_LogWarn_Good(t *core.T) {
	// LogWarn
	ax7Variant := "LogWarn:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "LogWarn:good"
	core.AssertContains(t, label, "LogWarn")
	core.AssertContains(t, label, "good")
}

func TestAction_LogWarn_Bad(t *core.T) {
	// LogWarn
	ax7Variant := "LogWarn:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "LogWarn:bad"
	core.AssertContains(t, label, "LogWarn")
	core.AssertContains(t, label, "bad")
}

func TestAction_LogWarn_Ugly(t *core.T) {
	// LogWarn
	ax7Variant := "LogWarn:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "LogWarn:ugly"
	core.AssertContains(t, label, "LogWarn")
	core.AssertContains(t, label, "ugly")
}
