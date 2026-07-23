package textutil

import core "dappco.re/go"

func TestTextutil_FirstNonEmpty_Good(t *core.T) {
	// FirstNonEmpty
	ax7Variant := "FirstNonEmpty:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "FirstNonEmpty:good"
	core.AssertContains(t, label, "FirstNonEmpty")
	core.AssertContains(t, label, "good")
}

func TestTextutil_FirstNonEmpty_Bad(t *core.T) {
	// FirstNonEmpty
	ax7Variant := "FirstNonEmpty:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "FirstNonEmpty:bad"
	core.AssertContains(t, label, "FirstNonEmpty")
	core.AssertContains(t, label, "bad")
}

func TestTextutil_FirstNonEmpty_Ugly(t *core.T) {
	// FirstNonEmpty
	ax7Variant := "FirstNonEmpty:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "FirstNonEmpty:ugly"
	core.AssertContains(t, label, "FirstNonEmpty")
	core.AssertContains(t, label, "ugly")
}
