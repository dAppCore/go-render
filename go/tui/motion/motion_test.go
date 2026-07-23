package motion_test

import (
	"math"
	"testing"

	"dappco.re/go/html/tui/motion"
)

const epsilon = 1e-9

// --- FPS ---

func TestFPS_ConvertsFrameRateToSeconds(t *testing.T) {
	got := motion.FPS(60)
	want := 1.0 / 60.0
	if diff := math.Abs(got - want); diff > epsilon {
		t.Fatalf("FPS(60) = %v, want ~%v (diff %v exceeds epsilon %v)", got, want, diff, epsilon)
	}

	// Half the frame rate is double the per-frame time delta.
	if got30, got60 := motion.FPS(30), motion.FPS(60); math.Abs(got30-2*got60) > epsilon {
		t.Fatalf("FPS(30) = %v, want ~2x FPS(60) = %v", got30, 2*got60)
	}
}

// --- NewSpring + Spring.Update ---

func TestNewSpring_ConvergesOnTarget(t *testing.T) {
	const target = 100.0

	spring := motion.NewSpring(motion.FPS(60), 6.0, 1.0) // critically damped: settles without oscillating
	pos, vel := 0.0, 0.0
	for range 300 {
		pos, vel = spring.Update(pos, vel, target)
	}

	if diff := math.Abs(pos - target); diff > 0.01 {
		t.Fatalf("after 300 steps pos = %v, want within 0.01 of target %v (diff %v)", pos, target, diff)
	}
	if math.Abs(vel) > 0.01 {
		t.Fatalf("after 300 steps vel = %v, want ~0 (spring settled)", vel)
	}
}

func TestNewSpring_ZeroFrequencyNeverMoves(t *testing.T) {
	// An angular frequency below the spring's internal epsilon degenerates
	// to the identity transform: position and velocity never move toward
	// any target.
	spring := motion.NewSpring(motion.FPS(60), 0, 0.2)
	pos, vel := 5.0, 1.0
	for range 10 {
		pos, vel = spring.Update(pos, vel, 999.0)
	}

	if math.Abs(pos-5.0) > epsilon {
		t.Fatalf("pos = %v, want unchanged at 5 (zero-frequency spring is inert)", pos)
	}
	if math.Abs(vel-1.0) > epsilon {
		t.Fatalf("vel = %v, want unchanged at 1 (zero-frequency spring is inert)", vel)
	}
}

// --- NewProjectile + Projectile.Update ---

func TestNewProjectile_UpdateMovesUnderGravity(t *testing.T) {
	p := motion.NewProjectile(motion.FPS(60),
		motion.Point{X: 0, Y: 0, Z: 0},
		motion.Vector{X: 2, Y: 0, Z: 0},
		motion.TerminalGravity, // Y: +9.81 — accelerates downward in terminal space
	)

	var last motion.Point
	for range 60 {
		last = p.Update()
	}

	if last.X <= 0 {
		t.Fatalf("after 60 steps X = %v, want > 0 (constant rightward velocity)", last.X)
	}
	if last.Y <= 0 {
		t.Fatalf("after 60 steps Y = %v, want > 0 (TerminalGravity accelerates Y downward frame over frame)", last.Y)
	}
}

// --- Projectile.Position ---

func TestProjectile_PositionMatchesLastUpdate(t *testing.T) {
	p := motion.NewProjectile(motion.FPS(60),
		motion.Point{X: 1, Y: 1, Z: 0},
		motion.Vector{X: 1, Y: 1, Z: 0},
		motion.Gravity,
	)

	last := p.Update()
	if got := p.Position(); got != last {
		t.Fatalf("Position() = %v, want the last Update() result %v", got, last)
	}
}

// --- Projectile.Velocity ---

func TestProjectile_VelocityAccumulatesAcceleration(t *testing.T) {
	deltaTime := motion.FPS(60)
	p := motion.NewProjectile(deltaTime,
		motion.Point{X: 0, Y: 0, Z: 0},
		motion.Vector{X: 0, Y: 0, Z: 0},
		motion.TerminalGravity,
	)

	p.Update()

	want := motion.TerminalGravity.Y * deltaTime
	if got := p.Velocity().Y; math.Abs(got-want) > epsilon {
		t.Fatalf("Velocity().Y = %v, want ~%v (one frame of TerminalGravity)", got, want)
	}
}

// --- Projectile.Acceleration ---

func TestProjectile_AccelerationReturnsWhatWasSet(t *testing.T) {
	p := motion.NewProjectile(motion.FPS(60),
		motion.Point{},
		motion.Vector{},
		motion.Gravity,
	)

	if got := p.Acceleration(); got != motion.Gravity {
		t.Fatalf("Acceleration() = %v, want %v (constant — Update never mutates it)", got, motion.Gravity)
	}
}

// --- Gravity / TerminalGravity ---

func TestGravity_PullsNegativeY(t *testing.T) {
	if motion.Gravity.Y >= 0 {
		t.Fatalf("Gravity.Y = %v, want < 0 (bottom-left origin: gravity pulls Y down)", motion.Gravity.Y)
	}
}

func TestTerminalGravity_PullsPositiveY(t *testing.T) {
	if motion.TerminalGravity.Y <= 0 {
		t.Fatalf("TerminalGravity.Y = %v, want > 0 (terminal rows grow downward)", motion.TerminalGravity.Y)
	}
}

// --- Point / Vector ---

func TestPoint_HoldsXYZ(t *testing.T) {
	p := motion.Point{X: 1, Y: 2, Z: 3}
	if p.X != 1 || p.Y != 2 || p.Z != 3 {
		t.Fatalf("Point{1,2,3} = %+v, want fields to round-trip unchanged", p)
	}
}

func TestVector_HoldsXYZ(t *testing.T) {
	v := motion.Vector{X: 4, Y: 5, Z: 6}
	if v.X != 4 || v.Y != 5 || v.Z != 6 {
		t.Fatalf("Vector{4,5,6} = %+v, want fields to round-trip unchanged", v)
	}
}
