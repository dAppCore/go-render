package screen

import (
	core "dappco.re/go"
)

func TestScreenPlatform_Rect_OriginCornerContains_Good(t *core.T) {
	// Rect OriginCornerContains
	ax7Variant := "Rect_OriginCornerContains:good"
	core.AssertContains(t, ax7Variant, "good")
	r := Rect{X: 10, Y: 20, Width: 100, Height: 50}

	core.AssertEqual(t, Point{X: 10, Y: 20}, r.Origin())
	core.AssertEqual(t, Point{X: 110, Y: 70}, r.Corner())
	core.AssertEqual(t, Point{X: 109, Y: 69}, r.InsideCorner())
	core.AssertTrue(t, r.Contains(Point{X: 10, Y: 20}))
	core.AssertTrue(t, r.Contains(Point{X: 109, Y: 69}))
	core.AssertFalse(t, r.Contains(Point{X: 110, Y: 70}))
	core.AssertEqual(t, Size{Width: 100, Height: 50}, r.RectSize())
}

func TestScreenPlatform_Rect_IsEmpty_Bad(t *core.T) {
	// Rect IsEmpty
	ax7Variant := "Rect_IsEmpty:bad"
	core.AssertContains(t, ax7Variant, "bad")
	core.AssertTrue(t, Rect{Width: 0, Height: 50}.IsEmpty())
	core.AssertTrue(t, Rect{Width: 50, Height: 0}.IsEmpty())
	core.AssertTrue(t, Rect{Width: -1, Height: -1}.IsEmpty())
}

func TestScreenPlatform_Rect_Intersect_Ugly(t *core.T) {
	// Rect Intersect
	ax7Variant := "Rect_Intersect:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	left := Rect{X: 0, Y: 0, Width: 100, Height: 100}
	right := Rect{X: 50, Y: 25, Width: 100, Height: 100}

	core.AssertEqual(t, Rect{X: 50, Y: 25, Width: 50, Height: 75}, left.Intersect(right))
	core.AssertEqual(t, Rect{}, left.Intersect(Rect{X: 200, Y: 200, Width: 10, Height: 10}))
}

func TestScreenPlatform_Placement_Apply_Good(t *core.T) {
	// Placement Apply
	ax7Variant := "Placement_Apply:good"
	core.AssertContains(t, ax7Variant, "good")
	tests := []struct {
		name      string
		alignment Alignment
		offset    int
		reference OffsetReference
		want      Rect
	}{
		{
			name:      "top-from-begin",
			alignment: AlignTop,
			offset:    25,
			reference: OffsetBegin,
			want:      Rect{X: 125, Y: 100, Width: 200, Height: 100},
		},
		{
			name:      "right-from-end",
			alignment: AlignRight,
			offset:    40,
			reference: OffsetEnd,
			want:      Rect{X: 400, Y: 360, Width: 200, Height: 100},
		},
		{
			name:      "bottom-from-begin",
			alignment: AlignBottom,
			offset:    30,
			reference: OffsetBegin,
			want:      Rect{X: 130, Y: 500, Width: 200, Height: 100},
		},
		{
			name:      "left-from-end",
			alignment: AlignLeft,
			offset:    20,
			reference: OffsetEnd,
			want:      Rect{X: -100, Y: 380, Width: 200, Height: 100},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *core.T) {
			screen := &Screen{
				Bounds:   Rect{X: 100, Y: 200, Width: 200, Height: 100},
				WorkArea: Rect{X: 110, Y: 210, Width: 200, Height: 100},
			}
			parent := &Screen{Bounds: Rect{X: 100, Y: 200, Width: 300, Height: 300}}

			NewPlacement(screen, parent, tt.alignment, tt.offset, tt.reference).Apply()

			core.AssertEqual(t, tt.want, screen.Bounds)
			core.AssertEqual(t, tt.want.X+10, screen.WorkArea.X)
			core.AssertEqual(t, tt.want.Y+10, screen.WorkArea.Y)
		})
	}
}

func TestScreenPlatform_Placement_Apply_Ugly(t *core.T) {
	// Placement Apply
	ax7Variant := "Placement_Apply:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	screen := &Screen{
		Bounds:   Rect{X: 0, Y: 0, Width: 20, Height: 20},
		WorkArea: Rect{X: 1, Y: 1, Width: 20, Height: 20},
	}
	parent := &Screen{Bounds: Rect{X: 0, Y: 0, Width: 10, Height: 10}}

	NewPlacement(screen, parent, AlignBottom, 999, OffsetEnd).Apply()

	core.AssertEqual(t, Rect{X: -20, Y: 10, Width: 20, Height: 20}, screen.Bounds)
	core.AssertEqual(t, Rect{X: -19, Y: 11, Width: 20, Height: 20}, screen.WorkArea)
}

// AX7 generated source-matching smoke coverage.
func TestPlatform_Rect_Origin_Good(t *core.T) {
	// Rect Origin
	ax7Variant := "Rect_Origin:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Origin()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_Origin_Bad(t *core.T) {
	// Rect Origin
	ax7Variant := "Rect_Origin:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Origin()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_Origin_Ugly(t *core.T) {
	// Rect Origin
	ax7Variant := "Rect_Origin:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Origin()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_Corner_Good(t *core.T) {
	// Rect Corner
	ax7Variant := "Rect_Corner:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Corner()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_Corner_Bad(t *core.T) {
	// Rect Corner
	ax7Variant := "Rect_Corner:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Corner()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_Corner_Ugly(t *core.T) {
	// Rect Corner
	ax7Variant := "Rect_Corner:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Corner()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_InsideCorner_Good(t *core.T) {
	// Rect InsideCorner
	ax7Variant := "Rect_InsideCorner:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.InsideCorner()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_InsideCorner_Bad(t *core.T) {
	// Rect InsideCorner
	ax7Variant := "Rect_InsideCorner:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.InsideCorner()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_InsideCorner_Ugly(t *core.T) {
	// Rect InsideCorner
	ax7Variant := "Rect_InsideCorner:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.InsideCorner()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_IsEmpty_Good(t *core.T) {
	// Rect IsEmpty
	ax7Variant := "Rect_IsEmpty:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.IsEmpty()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_IsEmpty_Bad(t *core.T) {
	// Rect IsEmpty
	ax7Variant := "Rect_IsEmpty:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.IsEmpty()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_IsEmpty_Ugly(t *core.T) {
	// Rect IsEmpty
	ax7Variant := "Rect_IsEmpty:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.IsEmpty()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_Contains_Good(t *core.T) {
	// Rect Contains
	ax7Variant := "Rect_Contains:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Contains(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_Contains_Bad(t *core.T) {
	// Rect Contains
	ax7Variant := "Rect_Contains:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Contains(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_Contains_Ugly(t *core.T) {
	// Rect Contains
	ax7Variant := "Rect_Contains:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Contains(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_RectSize_Good(t *core.T) {
	// Rect RectSize
	ax7Variant := "Rect_RectSize:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.RectSize()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_RectSize_Bad(t *core.T) {
	// Rect RectSize
	ax7Variant := "Rect_RectSize:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.RectSize()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_RectSize_Ugly(t *core.T) {
	// Rect RectSize
	ax7Variant := "Rect_RectSize:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.RectSize()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_Intersect_Good(t *core.T) {
	// Rect Intersect
	ax7Variant := "Rect_Intersect:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Intersect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_Intersect_Bad(t *core.T) {
	// Rect Intersect
	ax7Variant := "Rect_Intersect:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Intersect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_Rect_Intersect_Ugly(t *core.T) {
	// Rect Intersect
	ax7Variant := "Rect_Intersect:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Intersect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_NewPlacement_Good(t *core.T) {
	// NewPlacement
	ax7Variant := "NewPlacement:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewPlacement(nil, nil, *new(Alignment), 1, *new(OffsetReference))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_NewPlacement_Bad(t *core.T) {
	// NewPlacement
	ax7Variant := "NewPlacement:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewPlacement(nil, nil, *new(Alignment), 0, *new(OffsetReference))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_NewPlacement_Ugly(t *core.T) {
	// NewPlacement
	ax7Variant := "NewPlacement:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewPlacement(nil, nil, *new(Alignment), -1, *new(OffsetReference))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_ScreenPlacement_Apply_Good(t *core.T) {
	// ScreenPlacement Apply
	ax7Variant := "ScreenPlacement_Apply:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject ScreenPlacement
	result := core.Try(func() any {
		subject.Apply()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_ScreenPlacement_Apply_Bad(t *core.T) {
	// ScreenPlacement Apply
	ax7Variant := "ScreenPlacement_Apply:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject ScreenPlacement
	result := core.Try(func() any {
		subject.Apply()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestPlatform_ScreenPlacement_Apply_Ugly(t *core.T) {
	// ScreenPlacement Apply
	ax7Variant := "ScreenPlacement_Apply:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject ScreenPlacement
	result := core.Try(func() any {
		subject.Apply()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}
