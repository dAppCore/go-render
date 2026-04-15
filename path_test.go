package html

import (
	"testing"
)

func TestNestedLayout_PathChain_Good(t *testing.T) {
	inner := NewLayout("HCF").H(Raw("inner h")).C(Raw("inner c")).F(Raw("inner f"))
	outer := NewLayout("HLCRF").
		H(Raw("header")).L(inner).C(Raw("main")).R(Raw("right")).F(Raw("footer"))
	got := outer.Render(NewContext())

	// Inner layout paths must be prefixed with parent block ID
	for _, want := range []string{`data-block="L.0"`, `data-block="L.0.1"`, `data-block="L.0.2"`} {
		if !containsText(got, want) {
			t.Errorf("nested layout missing %q in:\n%s", want, got)
		}
	}

	// Outer layout must still have root-level paths
	for _, want := range []string{`data-block="H"`, `data-block="C"`, `data-block="F"`} {
		if !containsText(got, want) {
			t.Errorf("outer layout missing %q in:\n%s", want, got)
		}
	}
}

func TestNestedLayout_DeepNesting_Ugly(t *testing.T) {
	deepest := NewLayout("C").C(Raw("deep"))
	middle := NewLayout("C").C(deepest)
	outer := NewLayout("C").C(middle)
	got := outer.Render(NewContext())

	for _, want := range []string{`data-block="C"`, `data-block="C.0"`, `data-block="C.0.0"`} {
		if !containsText(got, want) {
			t.Errorf("deep nesting missing %q in:\n%s", want, got)
		}
	}
}

func TestBlockID_BuildsPath_Good(t *testing.T) {
	tests := []struct {
		path     string
		slot     byte
		rendered int
		want     string
	}{
		{"", 'H', 0, "H"},
		{"", 'H', 1, "H.1"},
		{"", 'F', 0, "F"},
		{"L.0", 'C', 0, "L.0"},
		{"L.0", 'C', 1, "L.0.1"},
		{"C.0.1", 'C', 0, "C.0.1"},
	}

	for _, tt := range tests {
		l := &Layout{path: tt.path}
		got := l.blockID(tt.slot, tt.rendered)
		if got != tt.want {
			t.Errorf("blockID(%q, %c, %d) = %q, want %q", tt.path, tt.slot, tt.rendered, got, tt.want)
		}
	}
}

func TestParseBlockID_ExtractsSlots_Good(t *testing.T) {
	tests := []struct {
		id   string
		want []byte
	}{
		{"L-0-C-0", []byte{'L', 'C'}},
		{"L.0.C.0", []byte{'L', 'C'}},
		{"L.0", []byte{'L'}},
		{"L.0.1", []byte{'L'}},
		{"C.0", []byte{'C'}},
		{"C.2.1", []byte{'C'}},
		{"C.0.1.2", []byte{'C'}},
		{"H", []byte{'H'}},
		{"C-0-C-0-C-0", []byte{'C', 'C', 'C'}},
		{"", nil},
	}

	for _, tt := range tests {
		got := ParseBlockID(tt.id)
		if len(got) != len(tt.want) {
			t.Errorf("ParseBlockID(%q) = %v, want %v", tt.id, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("ParseBlockID(%q)[%d] = %c, want %c", tt.id, i, got[i], tt.want[i])
			}
		}
	}
}

func TestParseBlockID_InvalidInput_Good(t *testing.T) {
	tests := []string{
		"L-1-C-0",
		"L-0-C",
		"L.0.",
		"L.0-C.0",
		"C.C.0",
		"C-0-0",
		"X",
	}

	for _, id := range tests {
		if got := ParseBlockID(id); got != nil {
			t.Errorf("ParseBlockID(%q) = %v, want nil", id, got)
		}
	}
}
