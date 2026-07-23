// SPDX-Licence-Identifier: EUPL-1.2

package table_test

import (
	"bytes"
	"slices"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/exp/teatest/v2"

	"dappco.re/go/html/tui/style"
	"dappco.re/go/html/tui/table"
)

func cols() []table.Column {
	return []table.Column{{Title: "Name", Width: 10}}
}

func rows() []table.Row {
	return []table.Row{{"Alfa"}, {"Bravo"}, {"Charlie"}}
}

// TestNew constructs a bare table and confirms its zero state before any
// Option runs: unfocused, no columns, no rows.
func TestNew(t *testing.T) {
	m := table.New()
	if m.Focused() {
		t.Fatal("New().Focused() = true, want false (unfocused until WithFocused/Focus)")
	}
	if got := len(m.Columns()); got != 0 {
		t.Fatalf("len(New().Columns()) = %d, want 0", got)
	}
	if got := len(m.Rows()); got != 0 {
		t.Fatalf("len(New().Rows()) = %d, want 0", got)
	}
}

// TestColumn builds a Column and confirms both fields round-trip — proving
// Column is a genuine alias for the real struct, not a shadow type.
func TestColumn(t *testing.T) {
	c := table.Column{Title: "ID", Width: 4}
	if c.Title != "ID" || c.Width != 4 {
		t.Fatalf("Column{Title: %q, Width: %d} = %+v, want Title=%q Width=%d", "ID", 4, c, "ID", 4)
	}
}

// TestRow proves Row is a genuine []string alias: it supports slice
// operations (append, indexing, len) unchanged.
func TestRow(t *testing.T) {
	r := table.Row{"a", "b"}
	r = append(r, "c")
	if got := len(r); got != 3 {
		t.Fatalf("len(Row) after append = %d, want 3", got)
	}
	if r[2] != "c" {
		t.Fatalf("Row[2] = %q, want %q", r[2], "c")
	}
}

// TestWithColumns sets the table columns at construction time and confirms
// Columns() reflects them.
func TestWithColumns(t *testing.T) {
	m := table.New(table.WithColumns(cols()))
	got := m.Columns()
	if len(got) != 1 || got[0].Title != "Name" {
		t.Fatalf("Columns() = %+v, want %+v", got, cols())
	}
}

// TestWithRows sets the table rows at construction time and confirms both
// Rows() and the default cursor's SelectedRow() reflect them.
func TestWithRows(t *testing.T) {
	m := table.New(table.WithColumns(cols()), table.WithRows(rows()))
	if got := len(m.Rows()); got != 3 {
		t.Fatalf("len(Rows()) = %d, want 3", got)
	}
	if got, want := m.SelectedRow(), rows()[0]; !slices.Equal(got, want) {
		t.Fatalf("SelectedRow() = %v, want %v (cursor starts at row 0)", got, want)
	}
}

// TestWithHeight proves the option forwards through to the viewport: a
// taller request yields a taller Height(). (The table reserves one line of
// whatever height is requested for its header row, so Height() trails the
// requested value by one — a bubbles characteristic, not this wrap's doing.)
func TestWithHeight(t *testing.T) {
	small := table.New(table.WithHeight(3))
	big := table.New(table.WithHeight(30))
	if got := big.Height(); got <= small.Height() {
		t.Fatalf("Height() with WithHeight(30) = %d, want > Height() with WithHeight(3) = %d", got, small.Height())
	}
}

// TestWithWidth proves the option forwards through to the viewport exactly:
// unlike height, width carries no header offset.
func TestWithWidth(t *testing.T) {
	m := table.New(table.WithWidth(42))
	if got := m.Width(); got != 42 {
		t.Fatalf("Width() = %d, want 42 (WithWidth(42))", got)
	}
}

// TestWithFocused sets the initial focus state — the gate Update checks
// before it will move the cursor on a key press.
func TestWithFocused(t *testing.T) {
	m := table.New(table.WithFocused(true))
	if !m.Focused() {
		t.Fatal("Focused() = false, want true (WithFocused(true))")
	}
}

// TestDefaultStyles proves DefaultStyles is real: its Header style pads
// rendered content, so a header cell measures wider than its bare text.
func TestDefaultStyles(t *testing.T) {
	s := table.DefaultStyles()
	rendered := s.Header.Render("x")
	if got, want := style.Measure(rendered), style.Measure("x"); got <= want {
		t.Fatalf("DefaultStyles().Header.Render(%q) measured %d, want > %d (header padding)", "x", got, want)
	}
}

// TestWithStyles proves a custom Styles value takes effect: a wider header
// padding renders a wider table than the default styling does.
func TestWithStyles(t *testing.T) {
	def := table.New(table.WithColumns(cols()), table.WithRows(rows()))
	wide := table.New(table.WithColumns(cols()), table.WithRows(rows()), table.WithStyles(table.Styles{
		Header:   style.New().Padding(0, 4),
		Cell:     style.New(),
		Selected: style.New(),
	}))

	if got, want := style.Measure(wide.View()), style.Measure(def.View()); got <= want {
		t.Fatalf("View() width with WithStyles padding = %d, want > default width %d", got, want)
	}
}

// TestDefaultKeyMap proves the returned KeyMap carries live, enabled
// bindings — LineDown matches both the down arrow and the vi "j" key.
func TestDefaultKeyMap(t *testing.T) {
	km := table.DefaultKeyMap()
	if !km.LineDown.Enabled() {
		t.Fatal("DefaultKeyMap().LineDown.Enabled() = false, want true")
	}
	if got := km.LineDown.Keys(); !slices.Contains(got, "down") || !slices.Contains(got, "j") {
		t.Fatalf("DefaultKeyMap().LineDown.Keys() = %v, want to contain %q and %q", got, "down", "j")
	}
}

// TestWithKeyMap installs a custom KeyMap at construction time and confirms
// the Model's exported KeyMap field carries it through.
func TestWithKeyMap(t *testing.T) {
	km := table.DefaultKeyMap()
	km.LineDown.SetKeys("n")

	m := table.New(table.WithKeyMap(km))
	if got := m.KeyMap.LineDown.Keys(); !slices.Contains(got, "n") {
		t.Fatalf("KeyMap.LineDown.Keys() = %v, want to contain %q", got, "n")
	}
}

// TestOption proves Option is the exact func(*Model) type table.New accepts,
// so a consumer can hand-write one instead of only using the With* helpers.
func TestOption(t *testing.T) {
	var custom table.Option = func(m *table.Model) {
		m.KeyMap = table.DefaultKeyMap()
		m.KeyMap.LineDown.SetKeys("n")
	}

	m := table.New(custom)
	if got := m.KeyMap.LineDown.Keys(); !slices.Contains(got, "n") {
		t.Fatalf("KeyMap.LineDown.Keys() after a hand-written Option = %v, want to contain %q", got, "n")
	}
}

// TestModel_SelectedRow proves SelectedRow tracks the cursor set by
// SetCursor.
func TestModel_SelectedRow(t *testing.T) {
	m := table.New(table.WithColumns(cols()), table.WithRows(rows()))
	m.SetCursor(2)
	if got, want := m.SelectedRow(), rows()[2]; !slices.Equal(got, want) {
		t.Fatalf("SelectedRow() after SetCursor(2) = %v, want %v", got, want)
	}
}

// TestModel_SetRows replaces the row data after construction and confirms
// Rows()/SelectedRow reflect the new data.
func TestModel_SetRows(t *testing.T) {
	m := table.New(table.WithColumns(cols()), table.WithRows(rows()))

	replacement := []table.Row{{"Zulu"}}
	m.SetRows(replacement)

	if got := len(m.Rows()); got != 1 {
		t.Fatalf("len(Rows()) after SetRows = %d, want 1", got)
	}
	if got, want := m.SelectedRow(), replacement[0]; !slices.Equal(got, want) {
		t.Fatalf("SelectedRow() after SetRows = %v, want %v", got, want)
	}
}

// TestModel_FocusBlur proves Focus/Blur toggle Focused() and gate whether a
// key press is allowed to move the cursor.
func TestModel_FocusBlur(t *testing.T) {
	m := table.New(table.WithColumns(cols()), table.WithRows(rows()))

	m.Focus()
	if !m.Focused() {
		t.Fatal("Focused() = false after Focus(), want true")
	}

	updated, _ := m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	if got, want := updated.SelectedRow(), rows()[1]; !slices.Equal(got, want) {
		t.Fatalf("SelectedRow() after a %q key press while focused = %v, want %v", "j", got, want)
	}

	m.Blur()
	if m.Focused() {
		t.Fatal("Focused() = true after Blur(), want false")
	}

	updated, _ = m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	if got, want := updated.SelectedRow(), rows()[0]; !slices.Equal(got, want) {
		t.Fatalf("SelectedRow() after a %q key press while blurred = %v, want %v (Update is a no-op unfocused)", "j", got, want)
	}
}

// harness is the minimal tea.Model a consumer wraps a Model in to drive it
// under a real bubbletea program: Model.Update returns a concrete
// table.Model, not the tea.Model interface, so the surrounding app supplies
// the outer Init/Update/View that satisfies tea.Model.
type harness struct {
	table.Model
}

func (h harness) Init() tea.Cmd { return nil }

func (h harness) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	h.Model, cmd = h.Model.Update(msg)
	return h, cmd
}

func (h harness) View() tea.View { return tea.NewView(h.Model.View()) }

// TestModel_NavigatesWithArrowKeys drives a focused table through a real
// bubbletea program via teatest, proving the wrap is a genuine key-driven
// widget rather than a plain struct. The viewport is exactly one data row
// tall, so each down/up press must scroll a different row into view before
// the cursor (SelectedRow) reflects the move.
func TestModel_NavigatesWithArrowKeys(t *testing.T) {
	allRows := []table.Row{{"Alfa"}, {"Bravo"}, {"Charlie"}, {"Delta"}, {"Echo"}}

	h := harness{Model: table.New(
		table.WithColumns(cols()),
		table.WithRows(allRows),
		table.WithFocused(true),
		table.WithHeight(2), // 1 header line + 1 data line
		table.WithWidth(20),
	)}

	tm := teatest.NewTestModel(t, h)

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Alfa"))
	})

	// Type drives the vi-style "j" binding — plain runes, the common case.
	tm.Type("j")
	tm.Type("j")
	tm.Type("j")
	tm.Type("j")

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Echo"))
	})

	// Send drives a real special-key message — the arrow-key case Type can't
	// produce.
	tm.Send(tea.KeyPressMsg{Code: tea.KeyUp})

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Delta"))
	})

	tm.Quit()

	final, ok := tm.FinalModel(t).(harness)
	if !ok {
		t.Fatalf("FinalModel() = %T, want harness", tm.FinalModel(t))
	}
	if got, want := final.SelectedRow(), (table.Row{"Delta"}); !slices.Equal(got, want) {
		t.Fatalf("SelectedRow() after 4 down presses + 1 up = %v, want %v", got, want)
	}
}
