package tui_test

import (
	"testing"

	htui "dappco.re/go/html/tui"
	"dappco.re/go/html/tui/help"
	"dappco.re/go/html/tui/key"
	"dappco.re/go/html/tui/list"
	"dappco.re/go/html/tui/markdown"
	"dappco.re/go/html/tui/spinner"
	"dappco.re/go/html/tui/textarea"
	"dappco.re/go/html/tui/textinput"
	"dappco.re/go/html/tui/viewport"
)

// TestSeam_ConstructsEveryWidget is the acceptance that a consumer can drop
// charmbracelet and build the same widgets from html/tui: every re-exported
// constructor must resolve and work through the seam.
func TestSeam_ConstructsEveryWidget(t *testing.T) {
	_ = list.New(nil, list.NewDefaultDelegate(), 20, 10)
	_ = textinput.New()
	_ = textarea.New()
	_ = viewport.New(viewport.WithWidth(20), viewport.WithHeight(10))
	_ = spinner.New()
	_ = help.New()

	b := key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit"))
	if !b.Enabled() {
		t.Fatal("binding should be enabled")
	}

	if _, err := markdown.New(markdown.WithWordWrap(40)); err != nil {
		t.Fatalf("markdown.New: %v", err)
	}
}

// TestSeam_LoopPrimitivesResolve proves the bubbletea event-loop re-exports at
// the tui root are callable (Batch/Quit produce commands/messages).
func TestSeam_LoopPrimitivesResolve(t *testing.T) {
	if htui.Batch() != nil {
		t.Fatal("Batch() of no commands should be nil")
	}
	if htui.Quit() == nil {
		t.Fatal("Quit() should produce a QuitMsg")
	}
}
