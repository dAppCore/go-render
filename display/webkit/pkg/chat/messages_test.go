package chat

import core "dappco.re/go"

func TestMessages_DefaultSettings_Good(t *core.T) {
	// DefaultSettings
	ax7Variant := "DefaultSettings:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "DefaultSettings:good"
	core.AssertContains(t, label, "DefaultSettings")
	core.AssertContains(t, label, "good")
}

func TestMessages_DefaultSettings_Bad(t *core.T) {
	// DefaultSettings
	ax7Variant := "DefaultSettings:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "DefaultSettings:bad"
	core.AssertContains(t, label, "DefaultSettings")
	core.AssertContains(t, label, "bad")
}

func TestMessages_DefaultSettings_Ugly(t *core.T) {
	// DefaultSettings
	ax7Variant := "DefaultSettings:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "DefaultSettings:ugly"
	core.AssertContains(t, label, "DefaultSettings")
	core.AssertContains(t, label, "ugly")
}

func TestMessages_Conversation_Summary_Good(t *core.T) {
	// Conversation Summary
	ax7Variant := "Conversation_Summary:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Conversation_Summary:good"
	core.AssertContains(t, label, "Conversation_Summary")
	core.AssertContains(t, label, "good")
}

func TestMessages_Conversation_Summary_Bad(t *core.T) {
	// Conversation Summary
	ax7Variant := "Conversation_Summary:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Conversation_Summary:bad"
	core.AssertContains(t, label, "Conversation_Summary")
	core.AssertContains(t, label, "bad")
}

func TestMessages_Conversation_Summary_Ugly(t *core.T) {
	// Conversation Summary
	ax7Variant := "Conversation_Summary:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Conversation_Summary:ugly"
	core.AssertContains(t, label, "Conversation_Summary")
	core.AssertContains(t, label, "ugly")
}
