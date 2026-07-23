// Package motion re-exports charmbracelet/harmonica through go-html's tui
// seam — spring physics and projectile motion, the smooth-animation ("feels
// alive") layer beneath panel transitions, value tweening, and smooth
// scroll. Swap the import path (charmbracelet/harmonica → html/tui/motion)
// and keep every motion.Spring / motion.NewSpring / motion.Projectile
// reference unchanged. Update (on both Spring and Projectile) and Position /
// Velocity / Acceleration (on Projectile) are methods, not package
// functions, so they need no re-export here — they come along for free
// since Spring and Projectile are genuine aliases.
//
// Spring usage — natural easing of a value toward a moving target, e.g. a
// panel sliding to its resting position or a scroll offset settling:
//
//	spring := motion.NewSpring(motion.FPS(60), 6.0, 0.2)
//
//	pos, vel := 0.0, 0.0
//	targetPos := 100.0
//	someUpdateLoop(func() {
//		pos, vel = spring.Update(pos, vel, targetPos)
//	})
//
// Projectile usage — ballistic motion under gravity, e.g. a particle or a
// toast falling/arcing across the terminal:
//
//	projectile := motion.NewProjectile(
//		motion.FPS(60),
//		motion.Point{X: 6.0, Y: 100.0, Z: 0.0},
//		motion.Vector{X: 2.0, Y: 0.0, Z: 0.0},
//		motion.TerminalGravity,
//	)
//
//	someUpdateLoop(func() {
//		pos := projectile.Update()
//	})
package motion

import "github.com/charmbracelet/harmonica"

type (
	// Spring is a damped harmonic oscillator: cached motion coefficients for
	// a given time step, frequency and damping ratio, replayed every frame
	// via Update to ease a value toward a target.
	Spring = harmonica.Spring
	// Projectile is a position, velocity and acceleration under simple
	// physics — advanced one frame at a time via Update.
	Projectile = harmonica.Projectile
	// Point is an X, Y, Z coordinate on a plane.
	Point = harmonica.Point
	// Vector is a magnitude-and-direction pair, represented the same way as
	// Point (X, Y, Z from the origin); used for velocity and acceleration.
	Vector = harmonica.Vector
)

var (
	// NewSpring computes a Spring's coefficients from a time delta (see
	// FPS), an angular frequency and a damping ratio. Damping < 1 is
	// under-damped (fastest to equilibrium, but overshoots and oscillates),
	// = 1 is critically damped (fastest without oscillating), > 1 is
	// over-damped (no oscillation, slower).
	NewSpring = harmonica.NewSpring
	// FPS converts a frame rate into the time-delta argument NewSpring and
	// NewProjectile expect. Prefer a real engine/render-loop delta over this
	// when one is available.
	FPS = harmonica.FPS
	// NewProjectile builds a Projectile from a time delta (see FPS) and
	// initial position, velocity and acceleration.
	NewProjectile = harmonica.NewProjectile
	// Gravity is a preset acceleration Vector for a coordinate plane whose
	// origin sits bottom-left (Y increases upward) — the conventional maths
	// orientation.
	Gravity = harmonica.Gravity
	// TerminalGravity is a preset acceleration Vector for a coordinate plane
	// whose origin sits top-left (Y increases downward) — the terminal's
	// own orientation, so a falling projectile's Y grows frame over frame.
	TerminalGravity = harmonica.TerminalGravity
)
