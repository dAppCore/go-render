// SPDX-Licence-Identifier: EUPL-1.2
//go:build !js

package main

import (
	"compress/gzip"
	"context"
	"testing"

	core "dappco.re/go"
	corepkg "dappco.re/go/core"
	coreio "dappco.re/go/io"
	process "dappco.re/go/process"
)

const (
	wasmGzLimit  = 1_048_576 // 1 MB gzip transfer size limit
	wasmRawLimit = 3_670_016 // 3.5 MB raw size limit
)

func TestCmdWasm_WASMBinarySizeGood(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping WASM build test in short mode")
	}

	if err := process.Init(corepkg.New()); err != nil {
		t.Fatalf("process.Init: %v", err)
	}

	dir := t.TempDir()
	out := core.Path(dir, "gohtml.wasm")

	output, err := process.RunWithOptions(context.Background(), process.RunOptions{
		Command: "go",
		Args:    []string{"build", "-ldflags=-s -w", "-o", out, "."},
		Dir:     ".",
		Env:     []string{"GOOS=js", "GOARCH=wasm"},
	})
	if err != nil {
		t.Fatalf("WASM build failed: %v: %s", err, output)
	}

	rawStr, err := coreio.Local.Read(out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rawBytes := []byte(rawStr)

	buf := core.NewBuilder()
	gz, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := gz.Write(rawBytes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Logf("WASM size: %d bytes raw, %d bytes gzip", len(rawBytes), buf.Len())

	if buf.Len() >= wasmGzLimit {
		t.Fatalf("WASM gzip size %d exceeds 1MB limit", buf.Len())
	}
	if len(rawBytes) >= wasmRawLimit {
		t.Fatalf("WASM raw size %d exceeds 3MB limit", len(rawBytes))
	}
}
