package table_test

import (
	"strings"
	"testing"

	"dappco.re/go/render/display/tui/style"
	"dappco.re/go/render/display/tui/style/table"
)

// StringData must satisfy Data -- the row/column source a Table reads from.
var _ table.Data = table.NewStringData()

// TestNew builds a table with headers and rows and confirms the rendered
// string carries every cell.
func TestNew(t *testing.T) {
	tb := table.New().Headers("NAME", "AGE").Row("Ada", "30").Row("Grace", "85")

	out := tb.String()
	if out == "" {
		t.Fatal("String() of a table with headers and rows must not be empty")
	}
	for _, want := range []string{"NAME", "AGE", "Ada", "30", "Grace", "85"} {
		if !strings.Contains(out, want) {
			t.Fatalf("String() = %q, missing %q", out, want)
		}
	}
}

// TestDefaultStyles proves DefaultStyles is the no-attribute StyleFunc: it
// renders text back unchanged regardless of which cell it is asked about.
func TestDefaultStyles(t *testing.T) {
	if got := table.DefaultStyles(table.HeaderRow, 0).Render("x"); got != "x" {
		t.Fatalf("DefaultStyles(HeaderRow, 0).Render(%q) = %q, want unchanged", "x", got)
	}
	if got := table.DefaultStyles(3, 1).Render("y"); got != "y" {
		t.Fatalf("DefaultStyles(3, 1).Render(%q) = %q, want unchanged", "y", got)
	}
}

// TestStyleFunc proves the exported StyleFunc type accepts a hand-written
// closure: HeaderRow distinguishes the header from the data rows the closure
// is invoked for, and a Style it returns for the header alone widens that
// table's rendered column relative to an unstyled one.
func TestStyleFunc(t *testing.T) {
	var sawHeader, sawDataRow bool
	var fn table.StyleFunc = func(row, _ int) style.Style {
		switch row {
		case table.HeaderRow:
			sawHeader = true
			return style.New().PaddingRight(4)
		default:
			sawDataRow = true
			return style.New()
		}
	}

	plain := table.New().Headers("A").Row("x")
	padded := table.New().Headers("A").Row("x").StyleFunc(fn)

	plainTop := strings.SplitN(plain.String(), "\n", 2)[0]
	paddedTop := strings.SplitN(padded.String(), "\n", 2)[0]

	if !sawHeader || !sawDataRow {
		t.Fatalf("StyleFunc should be invoked for the header (HeaderRow) and a data row, sawHeader=%v sawDataRow=%v", sawHeader, sawDataRow)
	}
	if len(paddedTop) <= len(plainTop) {
		t.Fatalf("padded table's top border (%q) should be wider than the plain one's (%q)", paddedTop, plainTop)
	}
}

// TestNewStringData builds a StringData from both the variadic constructor
// and the chained Item appender, and reads it back through At/Rows/Columns.
func TestNewStringData(t *testing.T) {
	d := table.NewStringData([]string{"a", "b"}).Item("c", "d", "e")

	if got := d.Rows(); got != 2 {
		t.Fatalf("Rows() = %d, want 2", got)
	}
	if got := d.Columns(); got != 3 {
		t.Fatalf("Columns() = %d, want 3 (the widest row)", got)
	}
	if got := d.At(1, 2); got != "e" {
		t.Fatalf("At(1, 2) = %q, want %q", got, "e")
	}
	if got := d.At(0, 2); got != "" {
		t.Fatalf("At(0, 2) = %q, want empty (row 0 only has 2 cells)", got)
	}
}

// TestNewFilter narrows a StringData to its even-indexed rows and confirms
// both Rows and At are reindexed relative to the filtered set, not the
// source.
func TestNewFilter(t *testing.T) {
	data := table.NewStringData([]string{"0"}, []string{"1"}, []string{"2"}, []string{"3"})
	f := table.NewFilter(data).Filter(func(row int) bool { return row%2 == 0 })

	if got := f.Rows(); got != 2 {
		t.Fatalf("Rows() = %d, want 2 (source rows 0 and 2)", got)
	}
	if got := f.At(0, 0); got != "0" {
		t.Fatalf("At(0, 0) = %q, want %q", got, "0")
	}
	if got := f.At(1, 0); got != "2" {
		t.Fatalf("At(1, 0) = %q, want %q (second filtered row is source row 2)", got, "2")
	}
	if got := f.Columns(); got != data.Columns() {
		t.Fatalf("Columns() = %d, want %d (Filter delegates column count to the wrapped Data)", got, data.Columns())
	}
}

// TestDataToMatrix reads a Data implementation back into a plain [][]string.
func TestDataToMatrix(t *testing.T) {
	data := table.NewStringData([]string{"a", "b"}, []string{"c", "d"})

	got := table.DataToMatrix(data)
	want := [][]string{{"a", "b"}, {"c", "d"}}

	if len(got) != len(want) {
		t.Fatalf("DataToMatrix rows = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if len(got[i]) != len(want[i]) {
			t.Fatalf("DataToMatrix[%d] columns = %d, want %d", i, len(got[i]), len(want[i]))
		}
		for j := range want[i] {
			if got[i][j] != want[i][j] {
				t.Fatalf("DataToMatrix[%d][%d] = %q, want %q", i, j, got[i][j], want[i][j])
			}
		}
	}
}
