// pkg/window/tiling.go
package window

import (
	core "dappco.re/go"

	"dappco.re/go/render/display/webkit/pkg/screen"
)

// TileMode defines how windows are arranged.
// Use: mode := window.TileModeLeftRight
type TileMode int

const (
	TileModeLeftHalf TileMode = iota
	TileModeRightHalf
	TileModeTopHalf
	TileModeBottomHalf
	TileModeTopLeft
	TileModeTopRight
	TileModeBottomLeft
	TileModeBottomRight
	TileModeLeftRight
	TileModeGrid
)

var tileModeNames = map[TileMode]string{
	TileModeLeftHalf: "left-half", TileModeRightHalf: "right-half",
	TileModeTopHalf: "top-half", TileModeBottomHalf: "bottom-half",
	TileModeTopLeft: "top-left", TileModeTopRight: "top-right",
	TileModeBottomLeft: "bottom-left", TileModeBottomRight: "bottom-right",
	TileModeLeftRight: "left-right", TileModeGrid: "grid",
}

// String returns the canonical layout name for the tile mode.
// Use: label := window.TileModeGrid.String()
func (m TileMode) String() string { return tileModeNames[m] }

// SnapPosition defines where a window snaps to.
// Use: pos := window.SnapRight
type SnapPosition int

const (
	SnapLeft SnapPosition = iota
	SnapRight
	SnapTop
	SnapBottom
	SnapTopLeft
	SnapTopRight
	SnapBottomLeft
	SnapBottomRight
	SnapCenter
)

var snapPositionNames = map[SnapPosition]string{
	SnapLeft: "left", SnapRight: "right",
	SnapTop: "top", SnapBottom: "bottom",
	SnapTopLeft: "top-left", SnapTopRight: "top-right",
	SnapBottomLeft: "bottom-left", SnapBottomRight: "bottom-right",
	SnapCenter: "center",
}

func (p SnapPosition) String() string { return snapPositionNames[p] }

// WorkflowLayout is a predefined arrangement for common tasks.
// Use: workflow := window.WorkflowCoding
type WorkflowLayout int

const (
	WorkflowCoding     WorkflowLayout = iota // 70/30 split
	WorkflowDebugging                        // 60/40 split
	WorkflowPresenting                       // maximised
	WorkflowSideBySide                       // 50/50 split
)

var workflowNames = map[WorkflowLayout]string{
	WorkflowCoding: "coding", WorkflowDebugging: "debugging",
	WorkflowPresenting: "presenting", WorkflowSideBySide: "side-by-side",
}

type Rect = screen.Rect
type Size = screen.Size

// WindowPlacement maps a window name to its suggested bounds.
type WindowPlacement struct {
	Name   string `json:"name"`
	Bounds Rect   `json:"bounds"`
}

const (
	goldenSplitNumerator   = 618
	standardSplitNumerator = 600
	splitDenominator       = 1000
)

// String returns the canonical workflow name.
// Use: label := window.WorkflowCoding.String()
func (w WorkflowLayout) String() string { return workflowNames[w] }

// BesideEditor returns the empty region beside an editor within a zero-origin screen.
//
// The returned rectangle uses the same vertical span as the visible portion of the
// editor and chooses the side with more remaining horizontal space. targetSize
// describes the available screen/work-area size.
func BesideEditor(editorBounds Rect, targetSize Size) Rect {
	screenBounds := Rect{Width: targetSize.Width, Height: targetSize.Height}
	editor := editorBounds.Intersect(screenBounds)
	if screenBounds.IsEmpty() || editor.IsEmpty() {
		return Rect{}
	}

	left := Rect{
		X:      screenBounds.X,
		Y:      editor.Y,
		Width:  max(editor.X-screenBounds.X, 0),
		Height: editor.Height,
	}
	right := Rect{
		X:      editor.X + editor.Width,
		Y:      editor.Y,
		Width:  max((screenBounds.X+screenBounds.Width)-(editor.X+editor.Width), 0),
		Height: editor.Height,
	}

	switch {
	case right.Width >= left.Width && right.Width > 0:
		return right
	case left.Width > 0:
		return left
	default:
		return Rect{}
	}
}

// SuggestLayout returns suggested placements for the provided windows on screenBounds.
//
// One window fills the screen, two windows use a golden-ratio split, and three or
// more windows fall back to a simple grid.
func SuggestLayout(windows []Window, screenBounds Rect) []WindowPlacement {
	if len(windows) == 0 || screenBounds.IsEmpty() {
		return nil
	}

	placements := make([]WindowPlacement, 0, len(windows))
	switch len(windows) {
	case 1:
		return []WindowPlacement{{
			Name:   windows[0].Name,
			Bounds: screenBounds,
		}}
	case 2:
		if screenBounds.Width >= screenBounds.Height {
			primaryWidth := proportionalSplit(screenBounds.Width, goldenSplitNumerator, splitDenominator)
			placements = append(placements,
				WindowPlacement{
					Name: windows[0].Name,
					Bounds: Rect{
						X: screenBounds.X, Y: screenBounds.Y,
						Width: primaryWidth, Height: screenBounds.Height,
					},
				},
				WindowPlacement{
					Name: windows[1].Name,
					Bounds: Rect{
						X: screenBounds.X + primaryWidth, Y: screenBounds.Y,
						Width: screenBounds.Width - primaryWidth, Height: screenBounds.Height,
					},
				},
			)
			return placements
		}

		primaryHeight := proportionalSplit(screenBounds.Height, goldenSplitNumerator, splitDenominator)
		placements = append(placements,
			WindowPlacement{
				Name: windows[0].Name,
				Bounds: Rect{
					X: screenBounds.X, Y: screenBounds.Y,
					Width: screenBounds.Width, Height: primaryHeight,
				},
			},
			WindowPlacement{
				Name: windows[1].Name,
				Bounds: Rect{
					X: screenBounds.X, Y: screenBounds.Y + primaryHeight,
					Width: screenBounds.Width, Height: screenBounds.Height - primaryHeight,
				},
			},
		)
		return placements
	default:
		cols, rows := gridDimensions(len(windows))
		cellWidth := screenBounds.Width / cols
		cellHeight := screenBounds.Height / rows
		for i, window := range windows {
			col := i % cols
			row := i / cols
			x := screenBounds.X + col*cellWidth
			y := screenBounds.Y + row*cellHeight
			right := x + cellWidth
			if col == cols-1 {
				right = screenBounds.X + screenBounds.Width
			}
			bottom := y + cellHeight
			if row == rows-1 {
				bottom = screenBounds.Y + screenBounds.Height
			}
			placements = append(placements, WindowPlacement{
				Name: window.Name,
				Bounds: Rect{
					X: x, Y: y,
					Width: right - x, Height: bottom - y,
				},
			})
		}
		return placements
	}
}

// FindEmptySpace returns the largest empty rectangle within screenBounds that fits minSize.
func FindEmptySpace(screenBounds Rect, existingWindows []Window, minSize Size) (Rect, bool) {
	if screenBounds.IsEmpty() {
		return Rect{}, false
	}

	reqWidth := max(minSize.Width, 1)
	reqHeight := max(minSize.Height, 1)
	if screenBounds.Width < reqWidth || screenBounds.Height < reqHeight {
		return Rect{}, false
	}

	occupied := make([]Rect, 0, len(existingWindows))
	xEdges := []int{screenBounds.X, screenBounds.X + screenBounds.Width}
	yEdges := []int{screenBounds.Y, screenBounds.Y + screenBounds.Height}
	for _, window := range existingWindows {
		rect := Rect{X: window.X, Y: window.Y, Width: window.Width, Height: window.Height}.Intersect(screenBounds)
		if rect.IsEmpty() {
			continue
		}
		occupied = append(occupied, rect)
		xEdges = append(xEdges, rect.X, rect.X+rect.Width)
		yEdges = append(yEdges, rect.Y, rect.Y+rect.Height)
	}

	xEdges = uniqueSorted(xEdges)
	yEdges = uniqueSorted(yEdges)

	best := Rect{}
	bestArea := -1
	for i := 0; i < len(xEdges); i++ {
		for j := i + 1; j < len(xEdges); j++ {
			for k := 0; k < len(yEdges); k++ {
				for l := k + 1; l < len(yEdges); l++ {
					candidate := Rect{
						X:      xEdges[i],
						Y:      yEdges[k],
						Width:  xEdges[j] - xEdges[i],
						Height: yEdges[l] - yEdges[k],
					}
					if candidate.Width < reqWidth || candidate.Height < reqHeight {
						continue
					}
					if !rectContains(screenBounds, candidate) || intersectsAny(candidate, occupied) {
						continue
					}
					area := candidate.Width * candidate.Height
					if area > bestArea {
						best = candidate
						bestArea = area
					}
				}
			}
		}
	}

	if bestArea < 0 {
		return Rect{}, false
	}
	return best, true
}

// ArrangePair places two windows side-by-side across screenBounds.
//
// Similar aspect ratios split the screen evenly. When one window is materially
// wider than the other, the wider one receives the larger 60% pane.
func ArrangePair(win1, win2 Window, screenBounds Rect) (Rect, Rect) {
	if screenBounds.IsEmpty() {
		return Rect{}, Rect{}
	}

	firstWidth := screenBounds.Width / 2
	aspect1 := windowAspect(win1)
	aspect2 := windowAspect(win2)
	delta := aspect1 - aspect2
	if delta < 0 {
		delta = -delta
	}
	if delta >= 0.35 {
		if aspect1 > aspect2 {
			firstWidth = proportionalSplit(screenBounds.Width, standardSplitNumerator, splitDenominator)
		} else {
			firstWidth = screenBounds.Width - proportionalSplit(screenBounds.Width, standardSplitNumerator, splitDenominator)
		}
	}

	first := Rect{
		X: screenBounds.X, Y: screenBounds.Y,
		Width: firstWidth, Height: screenBounds.Height,
	}
	second := Rect{
		X: screenBounds.X + firstWidth, Y: screenBounds.Y,
		Width: screenBounds.Width - firstWidth, Height: screenBounds.Height,
	}
	return first, second
}

func proportionalSplit(total, numerator, denominator int) int {
	if total <= 1 || denominator <= 0 {
		return total
	}
	split := (total*numerator + denominator/2) / denominator
	return min(max(split, 1), total-1)
}

func gridDimensions(count int) (int, int) {
	if count <= 0 {
		return 0, 0
	}
	cols := 1
	for cols*cols < count {
		cols++
	}
	rows := (count + cols - 1) / cols
	return cols, rows
}

func windowAspect(window Window) float64 {
	width := window.Width
	height := window.Height
	if width <= 0 || height <= 0 {
		return 1
	}
	return float64(width) / float64(height)
}

func layoutOrigin(origin []int) (int, int) {
	if len(origin) == 0 {
		return 0, 0
	}
	if len(origin) == 1 {
		return origin[0], 0
	}
	return origin[0], origin[1]
}

func (m *Manager) captureState(pw PlatformWindow) {
	if m.state == nil || pw == nil {
		return
	}
	m.state.CaptureState(pw)
}

func normalizeWindowForLayout(pw PlatformWindow) {
	if pw == nil {
		return
	}
	if pw.IsMaximised() || pw.IsMinimised() {
		pw.Restore()
	}
}

// TileWindows arranges the named windows in the given mode across the screen area.
func (m *Manager) TileWindows(mode TileMode, names []string, screenW, screenH int, origin ...int) resultFailure {
	originX, originY := layoutOrigin(origin)
	windows := make([]PlatformWindow, 0, len(names))
	for _, name := range names {
		pw, ok := m.Get(name)
		if !ok {
			return core.E("window.Manager.TileWindows", "window not found: "+name, nil)
		}
		windows = append(windows, pw)
	}
	if len(windows) == 0 {
		return core.E("window.Manager.TileWindows", "no windows to tile", nil)
	}
	for _, pw := range windows {
		normalizeWindowForLayout(pw)
	}

	halfW, halfH := screenW/2, screenH/2

	switch mode {
	case TileModeLeftRight:
		w := screenW / len(windows)
		for i, pw := range windows {
			pw.SetPosition(originX+i*w, originY)
			pw.SetSize(w, screenH)
			m.captureState(pw)
		}
	case TileModeGrid:
		cols := 2
		if len(windows) > 4 {
			cols = 3
		}
		cellW := screenW / cols
		for i, pw := range windows {
			row := i / cols
			col := i % cols
			rows := (len(windows) + cols - 1) / cols
			cellH := screenH / rows
			pw.SetPosition(originX+col*cellW, originY+row*cellH)
			pw.SetSize(cellW, cellH)
			m.captureState(pw)
		}
	case TileModeLeftHalf:
		for _, pw := range windows {
			pw.SetPosition(originX, originY)
			pw.SetSize(halfW, screenH)
			m.captureState(pw)
		}
	case TileModeRightHalf:
		for _, pw := range windows {
			pw.SetPosition(originX+halfW, originY)
			pw.SetSize(halfW, screenH)
			m.captureState(pw)
		}
	case TileModeTopHalf:
		for _, pw := range windows {
			pw.SetPosition(originX, originY)
			pw.SetSize(screenW, halfH)
			m.captureState(pw)
		}
	case TileModeBottomHalf:
		for _, pw := range windows {
			pw.SetPosition(originX, originY+halfH)
			pw.SetSize(screenW, halfH)
			m.captureState(pw)
		}
	case TileModeTopLeft:
		for _, pw := range windows {
			pw.SetPosition(originX, originY)
			pw.SetSize(halfW, halfH)
			m.captureState(pw)
		}
	case TileModeTopRight:
		for _, pw := range windows {
			pw.SetPosition(originX+halfW, originY)
			pw.SetSize(halfW, halfH)
			m.captureState(pw)
		}
	case TileModeBottomLeft:
		for _, pw := range windows {
			pw.SetPosition(originX, originY+halfH)
			pw.SetSize(halfW, halfH)
			m.captureState(pw)
		}
	case TileModeBottomRight:
		for _, pw := range windows {
			pw.SetPosition(originX+halfW, originY+halfH)
			pw.SetSize(halfW, halfH)
			m.captureState(pw)
		}
	}
	return nil
}

// SnapWindow snaps a window to a screen edge/corner/centre.
func (m *Manager) SnapWindow(name string, pos SnapPosition, screenW, screenH int, origin ...int) resultFailure {
	originX, originY := layoutOrigin(origin)
	pw, ok := m.Get(name)
	if !ok {
		return core.E("window.Manager.SnapWindow", "window not found: "+name, nil)
	}

	halfW, halfH := screenW/2, screenH/2

	switch pos {
	case SnapLeft:
		pw.SetPosition(originX, originY)
		pw.SetSize(halfW, screenH)
	case SnapRight:
		pw.SetPosition(originX+halfW, originY)
		pw.SetSize(halfW, screenH)
	case SnapTop:
		pw.SetPosition(originX, originY)
		pw.SetSize(screenW, halfH)
	case SnapBottom:
		pw.SetPosition(originX, originY+halfH)
		pw.SetSize(screenW, halfH)
	case SnapTopLeft:
		pw.SetPosition(originX, originY)
		pw.SetSize(halfW, halfH)
	case SnapTopRight:
		pw.SetPosition(originX+halfW, originY)
		pw.SetSize(halfW, halfH)
	case SnapBottomLeft:
		pw.SetPosition(originX, originY+halfH)
		pw.SetSize(halfW, halfH)
	case SnapBottomRight:
		pw.SetPosition(originX+halfW, originY+halfH)
		pw.SetSize(halfW, halfH)
	case SnapCenter:
		normalizeWindowForLayout(pw)
		cw, ch := pw.Size()
		pw.SetPosition(originX+(screenW-cw)/2, originY+(screenH-ch)/2)
	}
	m.captureState(pw)
	return nil
}

// StackWindows cascades windows with an offset.
func (m *Manager) StackWindows(names []string, offsetX, offsetY int, origin ...int) resultFailure {
	originX, originY := layoutOrigin(origin)
	for i, name := range names {
		pw, ok := m.Get(name)
		if !ok {
			return core.E("window.Manager.StackWindows", "window not found: "+name, nil)
		}
		pw.SetPosition(originX+i*offsetX, originY+i*offsetY)
		m.captureState(pw)
	}
	return nil
}

// ApplyWorkflow arranges windows in a predefined workflow layout.
func (m *Manager) ApplyWorkflow(workflow WorkflowLayout, names []string, screenW, screenH int, origin ...int) resultFailure {
	originX, originY := layoutOrigin(origin)
	if len(names) == 0 {
		return core.E("window.Manager.ApplyWorkflow", "no windows for workflow", nil)
	}

	switch workflow {
	case WorkflowCoding:
		// 70/30 split — main editor + terminal
		mainW := screenW * 70 / 100
		if pw, ok := m.Get(names[0]); ok {
			pw.SetPosition(originX, originY)
			pw.SetSize(mainW, screenH)
			m.captureState(pw)
		}
		if len(names) > 1 {
			if pw, ok := m.Get(names[1]); ok {
				pw.SetPosition(originX+mainW, originY)
				pw.SetSize(screenW-mainW, screenH)
				m.captureState(pw)
			}
		}
	case WorkflowDebugging:
		// 60/40 split
		mainW := screenW * 60 / 100
		if pw, ok := m.Get(names[0]); ok {
			pw.SetPosition(originX, originY)
			pw.SetSize(mainW, screenH)
			m.captureState(pw)
		}
		if len(names) > 1 {
			if pw, ok := m.Get(names[1]); ok {
				pw.SetPosition(originX+mainW, originY)
				pw.SetSize(screenW-mainW, screenH)
				m.captureState(pw)
			}
		}
	case WorkflowPresenting:
		// Maximise first window
		if pw, ok := m.Get(names[0]); ok {
			pw.SetPosition(originX, originY)
			pw.SetSize(screenW, screenH)
			m.captureState(pw)
		}
	case WorkflowSideBySide:
		return m.TileWindows(TileModeLeftRight, names, screenW, screenH, originX, originY)
	}
	return nil
}
