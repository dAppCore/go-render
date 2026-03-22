// Package main provides a build-time CLI for generating Web Component JS bundles.
// Reads a JSON slot map from stdin, writes the generated JS to stdout.
//
// Usage:
//
//	echo '{"H":"nav-bar","C":"main-content"}' | go run ./cmd/codegen/ > components.js
package main

import (
	"encoding/json"
	goio "io"
	"os"

	"dappco.re/go/core/html/codegen"
	log "dappco.re/go/core/log"
)

func run(r goio.Reader, w goio.Writer) error {
	data, err := goio.ReadAll(r)
	if err != nil {
		return log.E("codegen", "reading stdin", err)
	}

	var slots map[string]string
	if err := json.Unmarshal(data, &slots); err != nil {
		return log.E("codegen", "invalid JSON", err)
	}

	js, err := codegen.GenerateBundle(slots)
	if err != nil {
		return err
	}

	_, err = goio.WriteString(w, js)
	return err
}

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		log.Error("codegen failed", "err", err)
		os.Exit(1)
	}
}
