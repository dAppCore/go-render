// SPDX-Licence-Identifier: EUPL-1.2

package form_test

import (
	"testing"

	"dappco.re/go/html/display/tui/form"
)

// TestNewForm_Good builds a real form through the wrap — a generic Select and
// a required text Input sharing one Group — and proves it constructs,
// initialises, and round-trips a value both by Value pointer and by
// form.Get*.
func TestNewForm_Good(t *testing.T) {
	var name string

	sel := form.NewSelect[string]().Title("Colour").Key("colour").
		Options(form.NewOptions("a", "b")...)
	input := form.NewInput().Title("Name").Key("name").Value(&name)

	f := form.NewForm(form.NewGroup(sel, input))
	if f == nil {
		t.Fatal("NewForm returned nil")
	}

	if cmd := f.Init(); cmd == nil {
		t.Fatal("Init() returned a nil Cmd, want the group/window-size batch")
	}

	// Value(&name) binds the pointer directly: mutating it flows straight
	// through to GetValue without touching the running widget.
	name = "Ada"
	if got := input.GetValue(); got != "Ada" {
		t.Fatalf("Input.GetValue() = %v, want %q (Value pointer round-trip)", got, "Ada")
	}

	// Select defaults its accessor to the first option once Options is set;
	// advancing the field saves that value into the form's results map,
	// readable back through form.Get*.
	f.NextField()
	if got := f.GetString("colour"); got != "a" {
		t.Fatalf(`GetString("colour") = %q, want "a" (form.Get* round-trip)`, got)
	}
}

// TestNewSelect constructs a generic Select field for a non-string type
// parameter, proving the generic wrapper instantiates cleanly and that the
// cursor defaults to the first option.
func TestNewSelect(t *testing.T) {
	s := form.NewSelect[int]().Options(form.NewOptions(1, 2, 3)...)
	if got := s.GetValue(); got != 1 {
		t.Fatalf("Select[int].GetValue() = %v, want 1 (first option)", got)
	}
}

// TestNewMultiSelect constructs a generic MultiSelect field and proves its
// value is a same-typed, empty slice until an option is toggled.
func TestNewMultiSelect(t *testing.T) {
	m := form.NewMultiSelect[string]().Options(form.NewOptions("x", "y")...)
	if m == nil {
		t.Fatal("NewMultiSelect returned nil")
	}

	got, ok := m.GetValue().([]string)
	if !ok {
		t.Fatalf("MultiSelect.GetValue() = %T, want []string", m.GetValue())
	}
	if len(got) != 0 {
		t.Fatalf("MultiSelect.GetValue() = %v, want empty (nothing toggled yet)", got)
	}
}

// TestNewOption builds a single keyed option and confirms both halves.
func TestNewOption(t *testing.T) {
	o := form.NewOption("One", 1)
	if o.Key != "One" || o.Value != 1 {
		t.Fatalf("NewOption(%q, %d) = %+v, want Key=%q Value=%d", "One", 1, o, "One", 1)
	}
}

// TestNewOptions derives options from bare values, keying each by its
// fmt.Sprint form.
func TestNewOptions(t *testing.T) {
	opts := form.NewOptions(1, 2, 3)
	if got := len(opts); got != 3 {
		t.Fatalf("len(NewOptions(1,2,3)) = %d, want 3", got)
	}
	if opts[1].Key != "2" || opts[1].Value != 2 {
		t.Fatalf("NewOptions(1,2,3)[1] = %+v, want Key=%q Value=2", opts[1], "2")
	}
}
