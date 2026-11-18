// Package main provides the core-php binary — a standalone PHP/Laravel
// development tool with FrankenPHP embedding support.
package main

import (
	php "github.com/dappcore/go-render/go/engine/php"

	"github.com/dappcore/cli/go/pkg/cli"
)

func main() {
	cli.Main(
		cli.WithCommands("php", php.AddPHPRootCommands),
	)
}
