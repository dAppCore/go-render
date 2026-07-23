package style_test

import (
	"testing"

	"dappco.re/go/render/display/tui/style"
)

func TestAdaptiveColor_Resolve_PicksDarkOnADarkTerminal(t *testing.T) {
	c := style.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}

	r, g, b, _ := c.Resolve(true).RGBA()
	if r != 0xffff || g != 0xffff || b != 0xffff {
		t.Fatalf("Resolve(isDark=true) = %d,%d,%d, want the Dark colour white (65535,65535,65535)", r, g, b)
	}
}

func TestAdaptiveColor_Resolve_PicksLightOnALightTerminal(t *testing.T) {
	c := style.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}

	r, g, b, _ := c.Resolve(false).RGBA()
	if r != 0 || g != 0 || b != 0 {
		t.Fatalf("Resolve(isDark=false) = %d,%d,%d, want the Light colour black (0,0,0)", r, g, b)
	}
}

func TestAdaptiveColor_Resolve_ZeroValueResolvesWithoutPanicking(t *testing.T) {
	var c style.AdaptiveColor // both fields empty — an unset theme role

	if got := c.Resolve(true); got == nil {
		t.Fatal("Resolve on a zero AdaptiveColor returned a nil Paint, want a non-nil no-colour")
	}
}
