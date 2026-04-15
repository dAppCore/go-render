// SPDX-Licence-Identifier: EUPL-1.2
//go:build !js

package main

import (
	"compress/gzip"
	"context"
	"testing"

	core "dappco.re/go/core"
	coreio "dappco.re/go/core/io"
	process "dappco.re/go/core/process"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	wasmGzLimit  = 1_048_576 // 1 MB gzip transfer size limit
	wasmRawLimit = 3_670_016 // 3.5 MB raw size limit
)

func TestCmdWasm_WASMBinarySize_Good(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping WASM build test in short mode")
	}

	dir := t.TempDir()
	out := core.Path(dir, "gohtml.wasm")

	factory := process.NewService(process.Options{})
	serviceValue, err := factory(core.New())
	require.NoError(t, err)

	svc, ok := serviceValue.(*process.Service)
	require.True(t, ok, "process service factory returned %T", serviceValue)

	output, err := svc.RunWithOptions(context.Background(), process.RunOptions{
		Command: "go",
		Args:    []string{"build", "-ldflags=-s -w", "-o", out, "."},
		Dir:     ".",
		Env:     []string{"GOOS=js", "GOARCH=wasm"},
	})
	require.NoError(t, err, "WASM build failed: %s", output)

	rawStr, err := coreio.Local.Read(out)
	require.NoError(t, err)
	rawBytes := []byte(rawStr)

	buf := core.NewBuilder()
	gz, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
	require.NoError(t, err)
	_, err = gz.Write(rawBytes)
	require.NoError(t, err)
	require.NoError(t, gz.Close())

	t.Logf("WASM size: %d bytes raw, %d bytes gzip", len(rawBytes), buf.Len())

	assert.Less(t, buf.Len(), wasmGzLimit,
		"WASM gzip size %d exceeds 1MB limit", buf.Len())
	assert.Less(t, len(rawBytes), wasmRawLimit,
		"WASM raw size %d exceeds 3MB limit", len(rawBytes))
}
