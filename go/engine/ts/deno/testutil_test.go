package deno

import (
	"os"
	"testing"
)

func shortSocketDir(t *testing.T) string {
	t.Helper()

	// /tmp, not the working directory: socket paths must stay under macOS's
	// 104-byte sun_path limit however deep the checkout is.
	dir, err := os.MkdirTemp("/tmp", "cts")
	if err != nil {
		t.Fatalf("mktemp: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	return dir
}
