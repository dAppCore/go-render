// SPDX-Licence-Identifier: EUPL-1.2

// run.go: Run launches a Model (typically an App) as a Bubble Tea program with
// the house defaults and blocks until it exits -- the one-call entry a consumer
// uses instead of hand-rolling NewProgram + Program.Run.
package tui

// Run launches model as a Bubble Tea program and blocks until it exits,
// returning the final model and any run error. Program options (WithContext,
// WithFPS, WithoutSignalHandler, ...) pass through; altscreen and mouse are set
// by the Model's View in Bubble Tea v2, so Run does not take them.
//
// Example:
//
//	node, _ := ctml.Parse(src)
//	if _, err := tui.Run(tui.NewApp(node)); err != nil { /* handle */ }
func Run(model Model, opts ...ProgramOption) (Model, error) {
	return NewProgram(model, opts...).Run()
}
