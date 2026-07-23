//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package pipe

import (
	"testing"

	core "dappco.re/go"
)

type pipeRenderer struct {
	output []byte
	err    error
	entry  string
	input  any
}

func (r *pipeRenderer) Render(_ core.Context, entry string, input any) ([]byte, error) {
	r.entry = entry
	r.input = input
	return r.output, r.err
}

type shortWriter struct{}

func (shortWriter) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}
	return len(data) - 1, nil
}

func TestPipe_Emit_Good(t *testing.T) {
	renderer := &pipeRenderer{output: []byte("<p>ready</p>")}
	output := core.NewBuffer()
	err := Emit(renderer, "server.ts", output)

	core.AssertNoError(t, err)
	core.AssertEqual(t, "<p>ready</p>", output.String())
	core.AssertEqual(t, "server.ts", renderer.entry)
	core.AssertNil(t, renderer.input)
}

func TestPipe_Emit_Bad(t *testing.T) {
	renderer := &pipeRenderer{err: core.E("test.render", "failed", nil)}
	output := core.NewBuffer()
	err := Emit(renderer, "broken.ts", output)

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "render TypeScript entry")
	core.AssertEqual(t, 0, output.Len())
	core.AssertEqual(t, "broken.ts", renderer.entry)
}

func TestPipe_Emit_Ugly(t *testing.T) {
	renderer := &pipeRenderer{output: []byte("partial")}
	err := Emit(renderer, "server.ts", shortWriter{})

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "short write")
	core.AssertEqual(t, "server.ts", renderer.entry)
	core.AssertNil(t, renderer.input)
}
