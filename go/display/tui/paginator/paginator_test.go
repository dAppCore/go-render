// SPDX-Licence-Identifier: EUPL-1.2

package paginator_test

import (
	"bytes"
	"slices"
	"testing"
	"time"

	htui "dappco.re/go/html/display/tui"
	"dappco.re/go/html/display/tui/paginator"
	"github.com/charmbracelet/x/exp/teatest/v2"
)

// TestNew constructs a paginator through the wrap with no options and checks
// New's documented defaults, then re-constructs it with WithTotalPages and
// WithPerPage and confirms both options land on the Model.
func TestNew(t *testing.T) {
	m := paginator.New()
	if m.PerPage != 1 {
		t.Fatalf("New().PerPage = %d, want 1", m.PerPage)
	}
	if m.TotalPages != 1 {
		t.Fatalf("New().TotalPages = %d, want 1", m.TotalPages)
	}
	if m.Type != paginator.Arabic {
		t.Fatalf("New().Type = %v, want Arabic", m.Type)
	}

	m = paginator.New(paginator.WithTotalPages(7), paginator.WithPerPage(3))
	if m.TotalPages != 7 {
		t.Fatalf("New(WithTotalPages(7)).TotalPages = %d, want 7", m.TotalPages)
	}
	if m.PerPage != 3 {
		t.Fatalf("New(WithPerPage(3)).PerPage = %d, want 3", m.PerPage)
	}
}

// TestDefaultKeyMap proves the re-exported DefaultKeyMap resolves to the
// real bindings — page/arrow keys plus the vi h/l pair — and that New wires
// a fresh Model's KeyMap from it.
func TestDefaultKeyMap(t *testing.T) {
	km := paginator.DefaultKeyMap()

	if !slices.Equal(km.NextPage.Keys(), []string{"pgdown", "right", "l"}) {
		t.Fatalf("DefaultKeyMap().NextPage.Keys() = %v, want [pgdown right l]", km.NextPage.Keys())
	}
	if !slices.Equal(km.PrevPage.Keys(), []string{"pgup", "left", "h"}) {
		t.Fatalf("DefaultKeyMap().PrevPage.Keys() = %v, want [pgup left h]", km.PrevPage.Keys())
	}

	got := paginator.New().KeyMap
	if !slices.Equal(got.NextPage.Keys(), km.NextPage.Keys()) {
		t.Fatalf("New().KeyMap.NextPage.Keys() = %v, want DefaultKeyMap()'s %v", got.NextPage.Keys(), km.NextPage.Keys())
	}
}

// TestSetTotalPages proves the ceiling-division arithmetic that turns an
// item count into a page count, including the documented no-op for a
// non-positive item count.
func TestSetTotalPages(t *testing.T) {
	tests := []struct {
		name     string
		perPage  int
		items    int
		expected int
	}{
		{"exact multiple", 10, 20, 2},
		{"remainder rounds up", 10, 25, 3},
		{"fewer than one page", 10, 5, 1},
		{"non-positive items leaves TotalPages untouched", 10, 0, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := paginator.New(paginator.WithPerPage(tt.perPage), paginator.WithTotalPages(4))

			got := m.SetTotalPages(tt.items)
			if got != tt.expected {
				t.Fatalf("SetTotalPages(%d) = %d, want %d", tt.items, got, tt.expected)
			}
			if m.TotalPages != tt.expected {
				t.Fatalf("TotalPages after SetTotalPages(%d) = %d, want %d", tt.items, m.TotalPages, tt.expected)
			}
		})
	}
}

// TestNextPage proves NextPage advances the current page and refuses to
// page past the last one.
func TestNextPage(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		expected int
	}{
		{"advances one page", 0, 1},
		{"stops on the last page", 2, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := paginator.New(paginator.WithTotalPages(3))
			m.Page = tt.page

			m.NextPage()
			if m.Page != tt.expected {
				t.Fatalf("NextPage(): Page = %d, want %d", m.Page, tt.expected)
			}
		})
	}
}

// TestPrevPage proves PrevPage steps back one page and refuses to page
// before the first one.
func TestPrevPage(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		expected int
	}{
		{"steps back one page", 1, 0},
		{"stops on the first page", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := paginator.New(paginator.WithTotalPages(3))
			m.Page = tt.page

			m.PrevPage()
			if m.Page != tt.expected {
				t.Fatalf("PrevPage(): Page = %d, want %d", m.Page, tt.expected)
			}
		})
	}
}

// TestOnFirstPage proves the first-page boundary check.
func TestOnFirstPage(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		expected bool
	}{
		{"page zero is first", 0, true},
		{"page one is not first", 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := paginator.New(paginator.WithTotalPages(2))
			m.Page = tt.page

			if got := m.OnFirstPage(); got != tt.expected {
				t.Fatalf("OnFirstPage() = %t, want %t", got, tt.expected)
			}
		})
	}
}

// TestOnLastPage proves the last-page boundary check.
func TestOnLastPage(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		expected bool
	}{
		{"last index is last page", 1, true},
		{"page zero is not last", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := paginator.New(paginator.WithTotalPages(2))
			m.Page = tt.page

			if got := m.OnLastPage(); got != tt.expected {
				t.Fatalf("OnLastPage() = %t, want %t", got, tt.expected)
			}
		})
	}
}

// TestGetSliceBounds proves the start/end arithmetic used to slice a
// consumer's own backing data to the current page.
func TestGetSliceBounds(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		length    int
		wantStart int
		wantEnd   int
	}{
		{"first page, full page", 0, 25, 0, 10},
		{"middle page, full page", 1, 25, 10, 20},
		{"last page, partial page", 2, 25, 20, 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := paginator.New(paginator.WithPerPage(10))
			m.Page = tt.page

			start, end := m.GetSliceBounds(tt.length)
			if start != tt.wantStart || end != tt.wantEnd {
				t.Fatalf("GetSliceBounds(%d) = (%d, %d), want (%d, %d)", tt.length, start, end, tt.wantStart, tt.wantEnd)
			}
		})
	}
}

// TestItemsOnPage proves the item count derived for the current page,
// including the short final page.
func TestItemsOnPage(t *testing.T) {
	tests := []struct {
		name       string
		page       int
		totalItems int
		expected   int
	}{
		{"full page", 0, 25, 10},
		{"partial last page", 2, 25, 5},
		{"no items", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := paginator.New(paginator.WithPerPage(10))
			m.Page = tt.page

			if got := m.ItemsOnPage(tt.totalItems); got != tt.expected {
				t.Fatalf("ItemsOnPage(%d) = %d, want %d", tt.totalItems, got, tt.expected)
			}
		})
	}
}

// TestView renders both pagination styles directly, proving the Type
// constants (Arabic, Dots) select the documented output shapes.
func TestView(t *testing.T) {
	m := paginator.New(paginator.WithTotalPages(3))
	m.Page = 1

	m.Type = paginator.Arabic
	if got, want := m.View(), "2/3"; got != want {
		t.Fatalf("Arabic View() = %q, want %q", got, want)
	}

	m.Type = paginator.Dots
	if got, want := m.View(), "○•○"; got != want {
		t.Fatalf("Dots View() = %q, want %q", got, want)
	}
}

// harnessModel wraps a paginator.Model in the minimal tea.Model shape
// teatest.NewTestModel needs, so TestModel_TeatestKeyDriven exercises Update
// and View through a real bubbletea event loop rather than calling them
// directly.
type harnessModel struct {
	p paginator.Model
}

func (m harnessModel) Init() htui.Cmd { return nil }

func (m harnessModel) Update(msg htui.Msg) (htui.Model, htui.Cmd) {
	var cmd htui.Cmd
	m.p, cmd = m.p.Update(msg)
	return m, cmd
}

func (m harnessModel) View() htui.View {
	return htui.NewView(m.p.View())
}

// TestModel_TeatestKeyDriven drives a Dots-style paginator with real
// tea.KeyPressMsg values — right then left, the arrow-key half of
// DefaultKeyMap — through teatest's Program harness and confirms the
// rendered dots move to the second page and NextPage/PrevPage land back on
// the first.
func TestModel_TeatestKeyDriven(t *testing.T) {
	p := paginator.New(paginator.WithTotalPages(3))
	p.Type = paginator.Dots

	tm := teatest.NewTestModel(t, harnessModel{p: p})

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("•○○"))
	}, teatest.WithDuration(3*time.Second))

	tm.Send(htui.KeyPressMsg{Code: htui.KeyRight})
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("○•○"))
	}, teatest.WithDuration(3*time.Second))

	tm.Send(htui.KeyPressMsg{Code: htui.KeyLeft})
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("•○○"))
	}, teatest.WithDuration(3*time.Second))

	if err := tm.Quit(); err != nil {
		t.Fatalf("Quit() error = %v", err)
	}

	final, ok := tm.FinalModel(t, teatest.WithFinalTimeout(3*time.Second)).(harnessModel)
	if !ok {
		t.Fatalf("FinalModel() type = %T, want harnessModel", tm.FinalModel(t))
	}
	if final.p.Page != 0 {
		t.Fatalf("after right,left: Page = %d, want 0 (back on the first page)", final.p.Page)
	}
}
