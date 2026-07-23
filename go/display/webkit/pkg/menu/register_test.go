package menu

import core "dappco.re/go"

func TestRegister_Register_Good(t *core.T) {
	// Register
	ax7Variant := "Register:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Register:good"
	core.AssertContains(t, label, "Register")
	core.AssertContains(t, label, "good")
}

func TestRegister_Register_Bad(t *core.T) {
	// Register
	ax7Variant := "Register:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Register:bad"
	core.AssertContains(t, label, "Register")
	core.AssertContains(t, label, "bad")
}

func TestRegister_Register_Ugly(t *core.T) {
	// Register
	ax7Variant := "Register:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Register:ugly"
	core.AssertContains(t, label, "Register")
	core.AssertContains(t, label, "ugly")
}
