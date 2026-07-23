package motion_test

import (
	"fmt"
	"math"

	"dappco.re/go/render/display/tui/motion"
)

// ExampleNewSpring runs a Spring toward a target position over a run of
// frames and reports whether it settled — the shape a consumer follows to
// ease a panel, a scroll offset, or any tweened value toward a target
// without ever importing charmbracelet.
func ExampleNewSpring() {
	const target = 100.0

	spring := motion.NewSpring(motion.FPS(60), 6.0, 1.0)
	pos, vel := 0.0, 0.0
	for range 300 {
		pos, vel = spring.Update(pos, vel, target)
	}

	converged := math.Abs(pos-target) < 0.01 && math.Abs(vel) < 0.01
	fmt.Println("converged:", converged)
	// Output: converged: true
}
