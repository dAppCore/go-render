// Package main provides the core-php binary — a standalone PHP/Laravel
// development tool with FrankenPHP embedding support.
package main

import (
	php "dappco.re/go/render/engine/php"

	"dappco.re/go/cli/pkg/cli"
)

func main() {
	cli.Main(
		cli.WithCommands("php", php.AddPHPRootCommands),
	)
}
