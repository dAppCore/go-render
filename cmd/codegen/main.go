//go:build !js

// Package main provides a build-time CLI for generating Web Component JS bundles.
// Reads a JSON slot map from stdin, writes the generated JS to stdout.
//
// Usage:
//
//	echo '{"H":"nav-bar","C":"main-content"}' | go run ./cmd/codegen/ > components.js
package main

import (
	goio "io"

	core "dappco.re/go/core"
	"dappco.re/go/core/html/codegen"
	coreio "dappco.re/go/core/io"
	log "dappco.re/go/core/log"
)

func run(r goio.Reader, w goio.Writer) error {
	data, err := goio.ReadAll(r)
	if err != nil {
		return log.E("codegen", "reading stdin", err)
	}

	var slots map[string]string
	if result := core.JSONUnmarshal(data, &slots); !result.OK {
		err, _ := result.Value.(error)
		return log.E("codegen", "invalid JSON", err)
	}

	js, err := codegen.GenerateBundle(slots)
	if err != nil {
		return log.E("codegen", "generate bundle", err)
	}

	_, err = goio.WriteString(w, js)
	if err != nil {
		return log.E("codegen", "writing bundle", err)
	}
	return nil
}

func main() {
	stdin, err := coreio.Local.Open("/dev/stdin")
	if err != nil {
		panic(log.E("codegen.main", "open stdin", err))
	}

	stdout, err := coreio.Local.Create("/dev/stdout")
	if err != nil {
		panic(log.E("codegen.main", "open stdout", err))
	}
	defer func() {
		_ = stdin.Close()
		_ = stdout.Close()
	}()

	if err := run(stdin, stdout); err != nil {
		log.Error("codegen failed", "err", err)
		panic(err)
	}
}
