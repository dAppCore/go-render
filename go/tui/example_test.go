package tui_test

import (
	"fmt"

	htui "dappco.re/go/html/tui"
)

// houseModel is the smallest tea.Model that turns on the two capabilities a
// consumer most often reaches for: the alternate screen and cell-motion
// mouse tracking. Both are View fields rather than Program options — the
// mechanism AltScreen and MouseMode use now (see this file's package note)
// — so they are set here, in View, rather than passed to NewProgram.
type houseModel struct{}

func (houseModel) Init() htui.Cmd { return nil }

func (houseModel) Update(msg htui.Msg) (htui.Model, htui.Cmd) { return houseModel{}, nil }

func (houseModel) View() htui.View {
	v := htui.NewView("house")
	v.AltScreen = true
	v.MouseMode = htui.MouseModeCellMotion
	return v
}

// ExampleNewProgram builds a Program with a couple of this file's house
// options — WithoutSignalHandler and WithFPS — and a Model whose View turns
// on the alternate screen and mouse tracking, entirely through
// dappco.re/go/html/tui: no charmbracelet import anywhere in this file.
func ExampleNewProgram() {
	p := htui.NewProgram(houseModel{}, htui.WithoutSignalHandler(), htui.WithFPS(30))
	fmt.Println(p != nil)

	v := houseModel{}.View()
	fmt.Println(v.AltScreen, v.MouseMode == htui.MouseModeCellMotion)

	// Output:
	// true
	// true true
}
