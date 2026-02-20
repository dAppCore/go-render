// Package main provides a build-time CLI for generating Web Component JS bundles.
// Reads a JSON slot map from stdin, writes the generated JS to stdout.
//
// Usage:
//
//	echo '{"H":"nav-bar","C":"main-content"}' | go run ./cmd/codegen/ > components.js
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"forge.lthn.ai/core/go-html/codegen"
)

func run(r io.Reader, w io.Writer) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("codegen: reading stdin: %w", err)
	}

	var slots map[string]string
	if err := json.Unmarshal(data, &slots); err != nil {
		return fmt.Errorf("codegen: invalid JSON: %w", err)
	}

	js, err := codegen.GenerateBundle(slots)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, js)
	return err
}

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
