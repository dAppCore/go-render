//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package pipe

import (
	core "dappco.re/go"
	tsengine "dappco.re/go/render/engine/ts"
)

// Emit renders entry once with a nil input and writes the exact output bytes
// to writer.
func Emit(engine tsengine.Renderer, entry string, writer core.Writer) error {
	if engine == nil {
		return core.E("pipe.Emit", "render engine is nil", nil)
	}
	if writer == nil {
		return core.E("pipe.Emit", "writer is nil", nil)
	}

	output, err := engine.Render(core.Background(), entry, nil)
	if err != nil {
		return core.E("pipe.Emit", "render TypeScript entry", err)
	}
	written, err := writer.Write(output)
	if err != nil {
		return core.E("pipe.Emit", "write rendered output", err)
	}
	if written != len(output) {
		return core.E("pipe.Emit", "short write", nil)
	}
	return nil
}
