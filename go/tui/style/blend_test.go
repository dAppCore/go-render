package style_test

import (
	"testing"

	"dappco.re/go/html/tui/style"
)

func TestBlend1D_MidpointLiesBetweenTheStops(t *testing.T) {
	steps := style.Blend1D(3, style.Color("#000000"), style.Color("#ffffff"))
	if len(steps) != 3 {
		t.Fatalf("Blend1D(3, black, white) returned %d colours, want 3", len(steps))
	}

	r0, g0, b0, _ := steps[0].RGBA()
	r1, _, _, _ := steps[1].RGBA()
	r2, g2, b2, _ := steps[2].RGBA()

	if r0 != 0 || g0 != 0 || b0 != 0 {
		t.Fatalf("first stop = %d,%d,%d, want black (0,0,0)", r0, g0, b0)
	}
	if r2 != 0xffff || g2 != 0xffff || b2 != 0xffff {
		t.Fatalf("last stop = %d,%d,%d, want white (65535,65535,65535)", r2, g2, b2)
	}
	if !(r0 < r1 && r1 < r2) {
		t.Fatalf("midpoint red=%d is not strictly between %d and %d", r1, r0, r2)
	}
}

func TestBlend2D_RowIsMonotonicTowardTheAngle(t *testing.T) {
	grid := style.Blend2D(3, 1, 0, style.Color("#000000"), style.Color("#ffffff"))
	if len(grid) != 3 {
		t.Fatalf("Blend2D(3, 1, 0, black, white) returned %d colours, want 3 (width*height)", len(grid))
	}

	r0, _, _, _ := grid[0].RGBA()
	r1, _, _, _ := grid[1].RGBA()
	r2, _, _, _ := grid[2].RGBA()

	if !(r0 <= r1 && r1 <= r2) {
		t.Fatalf("row is not monotonic left-to-right at angle 0: %d, %d, %d", r0, r1, r2)
	}
	if r0 == r2 {
		t.Fatalf("first and last cell are equal (%d) — angle 0 should shade left-to-right", r0)
	}
}
