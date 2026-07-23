package style

import "charm.land/lipgloss/v2"

// Blend1D blends a series of colour stops together along one dimension,
// across the given number of steps, in the CIE L*a*b* colour space —
// perceptually-even interpolation rather than a flat RGB lerp. This is the
// primitive behind a linear gradient, e.g.
//
//	ramp := style.Blend1D(10, style.Color("#FF6B6B"), style.Color("#4ECDC4"))
//	// ramp[0] is the first stop, ramp[len(ramp)-1] the last, with
//	// perceptually-even colours between.
//
// It can back the gradients tui/anim currently builds by hand via a direct
// go-colorful BlendLuv call.
func Blend1D(steps int, stops ...Paint) []Paint {
	return lipgloss.Blend1D(steps, stops...)
}

// Blend2D blends a series of colour stops together across a width×height
// grid at the given angle in degrees (0 = left-to-right), in the CIE
// L*a*b* colour space, returned in row-major order (index = y*width+x) —
// the primitive behind a two-dimensional colour wash across a panel.
func Blend2D(width, height int, angle float64, stops ...Paint) []Paint {
	return lipgloss.Blend2D(width, height, angle, stops...)
}
