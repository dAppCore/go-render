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
	for _, want := range []string{`data-block="L-0-H-0"`, `data-block="L-0-C-0"`, `data-block="L-0-F-0"`} {
		if !containsText(got, want) {
			t.Errorf("nested layout missing %q in:\n%s", want, got)
		}
	}

	// Outer layout must still have root-level paths
	for _, want := range []string{`data-block="H-0"`, `data-block="C-0"`, `data-block="F-0"`} {
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

	for _, want := range []string{`data-block="C-0"`, `data-block="C-0-C-0"`, `data-block="C-0-C-0-C-0"`} {
		if !containsText(got, want) {
			t.Errorf("deep nesting missing %q in:\n%s", want, got)
		}
	}
}

func TestBlockID_BuildsPath_Good(t *testing.T) {
	tests := []struct {
		path string
		slot byte
		want string
	}{
		{"", 'H', "H-0"},
		{"L-0-", 'C', "L-0-C-0"},
		{"C-0-C-0-", 'C', "C-0-C-0-C-0"},
		{"", 'F', "F-0"},
	}

	for _, tt := range tests {
		l := &Layout{path: tt.path}
		got := l.blockID(tt.slot)
		if got != tt.want {
			t.Errorf("blockID(%q, %c) = %q, want %q", tt.path, tt.slot, got, tt.want)
		}
	}
}

func TestParseBlockID_ExtractsSlots_Good(t *testing.T) {
	tests := []struct {
		id   string
		want []byte
	}{
		{"L-0-C-0", []byte{'L', 'C'}},
		{"H-0", []byte{'H'}},
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
