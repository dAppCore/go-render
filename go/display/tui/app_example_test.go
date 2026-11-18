// SPDX-Licence-Identifier: EUPL-1.2

package tui_test

import (
	"fmt"

	tui "github.com/dappcore/go-render/go/display/tui"
	ctml "github.com/dappcore/go-render/go/engine/ctml"
)

// ExampleNewApp builds a manager around a .ctml document -- the shape a consumer
// follows to turn markup into a runnable terminal screen. tui.Run(app) launches
// it; here we just confirm it constructs and is ready to run, without importing
// charmbracelet.
func ExampleNewApp() {
	node, _ := ctml.Parse([]byte(`<h1>Welcome</h1>`))
	app := tui.NewApp(node)
	fmt.Println(app.Init() != nil) // Init requests the window size -> a Cmd
	// Output: true
}
