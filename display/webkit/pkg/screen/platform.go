// pkg/screen/platform.go
package screen

// Platform abstracts the screen/display backend.
//
//	core.WithService(screen.Register(wailsPlatform))
type Platform interface {
	GetAll() []Screen
	GetPrimary() *Screen
	// GetCurrent returns the most recently active screen, or the primary if unset.
	// current := platform.GetCurrent()
	GetCurrent() *Screen
}

// Screen describes a display/monitor.
type Screen struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	ScaleFactor      float64 `json:"scaleFactor"`
	Size             Size    `json:"size"`
	Bounds           Rect    `json:"bounds"`
	PhysicalBounds   Rect    `json:"physicalBounds"`
	WorkArea         Rect    `json:"workArea"`
	PhysicalWorkArea Rect    `json:"physicalWorkArea"`
	IsPrimary        bool    `json:"isPrimary"`
	Rotation         float64 `json:"rotation"`
}

// Rect represents a rectangle with position and dimensions.
//
//	if bounds.Contains(Point{X: cursor.X, Y: cursor.Y}) { highlightWindow() }
type Rect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Origin returns the top-left corner of the rectangle.
//
//	pt := bounds.Origin() // Point{X: bounds.X, Y: bounds.Y}
func (r Rect) Origin() Point {
	return Point{X: r.X, Y: r.Y}
}

// Corner returns the exclusive bottom-right corner (X+Width, Y+Height).
//
//	end := bounds.Corner() // Point{X: bounds.X+bounds.Width, Y: bounds.Y+bounds.Height}
func (r Rect) Corner() Point {
	return Point{X: r.X + r.Width, Y: r.Y + r.Height}
}

// InsideCorner returns the inclusive bottom-right corner (X+Width-1, Y+Height-1).
//
//	last := bounds.InsideCorner()
func (r Rect) InsideCorner() Point {
	return Point{X: r.X + r.Width - 1, Y: r.Y + r.Height - 1}
}

// IsEmpty reports whether the rectangle has non-positive area.
//
//	if r.IsEmpty() { return }
func (r Rect) IsEmpty() bool {
	return r.Width <= 0 || r.Height <= 0
}

// Contains reports whether point pt lies within the rectangle.
//
//	if workArea.Contains(windowOrigin) { snapToScreen() }
func (r Rect) Contains(pt Point) bool {
	return pt.X >= r.X && pt.X < r.X+r.Width && pt.Y >= r.Y && pt.Y < r.Y+r.Height
}

// RectSize returns the dimensions of the rectangle as a Size value.
//
//	sz := bounds.RectSize() // Size{Width: bounds.Width, Height: bounds.Height}
func (r Rect) RectSize() Size {
	return Size{Width: r.Width, Height: r.Height}
}

// Intersect returns the overlapping region of r and other, or an empty Rect if they do not overlap.
//
//	overlap := a.Intersect(b)
//	if !overlap.IsEmpty() { handleOverlap(overlap) }
func (r Rect) Intersect(other Rect) Rect {
	if r.IsEmpty() || other.IsEmpty() {
		return Rect{}
	}
	maxLeft := max(r.X, other.X)
	maxTop := max(r.Y, other.Y)
	minRight := min(r.X+r.Width, other.X+other.Width)
	minBottom := min(r.Y+r.Height, other.Y+other.Height)
	if minRight > maxLeft && minBottom > maxTop {
		return Rect{X: maxLeft, Y: maxTop, Width: minRight - maxLeft, Height: minBottom - maxTop}
	}
	return Rect{}
}

// Point is a two-dimensional coordinate.
//
//	centre := Point{X: bounds.X + bounds.Width/2, Y: bounds.Y + bounds.Height/2}
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Size represents dimensions.
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Alignment describes which edge of a parent screen a child screen is placed against.
type Alignment int

const (
	AlignTop    Alignment = iota // child is above parent
	AlignRight                   // child is to the right of parent
	AlignBottom                  // child is below parent
	AlignLeft                    // child is to the left of parent
)

// OffsetReference specifies whether the placement offset is measured from the
// beginning (top/left) or end (bottom/right) of the parent edge.
type OffsetReference int

const (
	OffsetBegin OffsetReference = iota // offset from top or left
	OffsetEnd                          // offset from bottom or right
)

// ScreenPlacement positions a screen relative to a parent screen.
//
//	placement := screen.NewPlacement(parent, AlignRight, 0, OffsetBegin)
//	placement.Apply()
type ScreenPlacement struct {
	screen          *Screen
	parent          *Screen
	alignment       Alignment
	offset          int
	offsetReference OffsetReference
}

// NewPlacement creates a ScreenPlacement that positions screen relative to parent.
//
//	p := NewPlacement(secondary, primary, AlignRight, 0, OffsetBegin)
//	p.Apply()
func NewPlacement(screen, parent *Screen, alignment Alignment, offset int, reference OffsetReference) ScreenPlacement {
	return ScreenPlacement{
		screen:          screen,
		parent:          parent,
		alignment:       alignment,
		offset:          offset,
		offsetReference: reference,
	}
}

// Apply moves screen.Bounds so that it sits against the specified edge of parent.
//
//	NewPlacement(s, p, AlignRight, 0, OffsetBegin).Apply()
func (p ScreenPlacement) Apply() {
	parentBounds := p.parent.Bounds
	screenBounds := p.screen.Bounds

	newX := parentBounds.X
	newY := parentBounds.Y
	offset := p.offset

	if p.alignment == AlignTop || p.alignment == AlignBottom {
		if p.offsetReference == OffsetEnd {
			offset = parentBounds.Width - offset - screenBounds.Width
		}
		offset = min(offset, parentBounds.Width)
		offset = max(offset, -screenBounds.Width)
		newX += offset
		if p.alignment == AlignTop {
			newY -= screenBounds.Height
		} else {
			newY += parentBounds.Height
		}
	} else {
		if p.offsetReference == OffsetEnd {
			offset = parentBounds.Height - offset - screenBounds.Height
		}
		offset = min(offset, parentBounds.Height)
		offset = max(offset, -screenBounds.Height)
		newY += offset
		if p.alignment == AlignLeft {
			newX -= screenBounds.Width
		} else {
			newX += parentBounds.Width
		}
	}

	workAreaOffsetX := p.screen.WorkArea.X - p.screen.Bounds.X
	workAreaOffsetY := p.screen.WorkArea.Y - p.screen.Bounds.Y
	p.screen.Bounds.X = newX
	p.screen.Bounds.Y = newY
	p.screen.WorkArea.X = newX + workAreaOffsetX
	p.screen.WorkArea.Y = newY + workAreaOffsetY
}
