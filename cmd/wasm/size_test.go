// SPDX-Licence-Identifier: EUPL-1.2
//go:build !js

package main

import (
	"bytes"
	"compress/gzip"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	coreio "dappco.re/go/core/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	wasmGzLimit  = 1_048_576 // 1 MB gzip transfer size limit
	wasmRawLimit = 3_670_016 // 3.5 MB raw size limit
)

func TestWASMBinarySize_Good(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping WASM build test in short mode")
	}

	dir := t.TempDir()
	out := filepath.Join(dir, "gohtml.wasm")

	cmd := exec.Command("go", "build", "-ldflags=-s -w", "-o", out, ".")
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "WASM build failed: %s", output)

	rawStr, err := coreio.Local.Read(out)
	require.NoError(t, err)
	raw := []byte(rawStr)

	var buf bytes.Buffer
	gz, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	require.NoError(t, err)
	_, err = gz.Write(raw)
	require.NoError(t, err)
	require.NoError(t, gz.Close())

	t.Logf("WASM size: %d bytes raw, %d bytes gzip", len(raw), buf.Len())

	assert.Less(t, buf.Len(), wasmGzLimit,
		"WASM gzip size %d exceeds 1MB limit", buf.Len())
	assert.Less(t, len(raw), wasmRawLimit,
		"WASM raw size %d exceeds 3MB limit", len(raw))
}
