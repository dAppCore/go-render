package php

import (
	core "dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
)

func addPHPPackagesCommands(c *core.Core, prefix string) {
	phpHelpCommand(c, phpCommandPath(prefix, "packages"), phpT("cmd.php.packages.short"))
	addPHPPackagesLinkCommand(c, prefix)
	addPHPPackagesUnlinkCommand(c, prefix)
	addPHPPackagesUpdateCommand(c, prefix)
	addPHPPackagesListCommand(c, prefix)
}

func addPHPPackagesLinkCommand(c *core.Core, prefix string) {
	path := phpCommandPath(prefix, "packages/link")
	phpFailureorCommand(c, path, phpT("cmd.php.packages.link.short"), func(opts core.Options) error {
		args := phpCommandLineFor(path, opts).Args()
		if len(args) < 1 {
			return phpFailure("requires at least 1 arg(s), only received %d", len(args))
		}

		cwdResult := core.Getwd()
		if !cwdResult.OK {
			err, _ := cwdResult.Value.(error)
			return core.E("php", phpT(i18nFailGetKey, workingDirectorySubject), err)
		}
		cwd, _ := cwdResult.Value.(string)

		cli.Print(cliLabelValueBlankFormat, dimStyle.Render(phpT(cmdPHPLabelKey)), phpT("cmd.php.packages.link.linking"))

		if err := LinkPackages(cwd, args); err != nil {
			return core.E("php", phpT("i18n.fail.link", "packages"), err)
		}

		cli.Print(cliSectionLabelValueFormat, successStyle.Render(phpLabel("done")), phpT("cmd.php.packages.link.done"))
		return nil
	})
}

func addPHPPackagesUnlinkCommand(c *core.Core, prefix string) {
	path := phpCommandPath(prefix, "packages/unlink")
	phpFailureorCommand(c, path, phpT("cmd.php.packages.unlink.short"), func(opts core.Options) error {
		args := phpCommandLineFor(path, opts).Args()
		if len(args) < 1 {
			return phpFailure("requires at least 1 arg(s), only received %d", len(args))
		}

		cwdResult := core.Getwd()
		if !cwdResult.OK {
			err, _ := cwdResult.Value.(error)
			return core.E("php", phpT(i18nFailGetKey, workingDirectorySubject), err)
		}
		cwd, _ := cwdResult.Value.(string)

		cli.Print(cliLabelValueBlankFormat, dimStyle.Render(phpT(cmdPHPLabelKey)), phpT("cmd.php.packages.unlink.unlinking"))

		if err := UnlinkPackages(cwd, args); err != nil {
			return core.E("php", phpT("i18n.fail.unlink", "packages"), err)
		}

		cli.Print(cliSectionLabelValueFormat, successStyle.Render(phpLabel("done")), phpT("cmd.php.packages.unlink.done"))
		return nil
	})
}

func addPHPPackagesUpdateCommand(c *core.Core, prefix string) {
	path := phpCommandPath(prefix, "packages/update")
	phpFailureorCommand(c, path, phpT("cmd.php.packages.update.short"), func(opts core.Options) error {
		args := phpCommandLineFor(path, opts).Args()
		cwdResult := core.Getwd()
		if !cwdResult.OK {
			err, _ := cwdResult.Value.(error)
			return core.E("php", phpT(i18nFailGetKey, workingDirectorySubject), err)
		}
		cwd, _ := cwdResult.Value.(string)

		cli.Print(cliLabelValueBlankFormat, dimStyle.Render(phpT(cmdPHPLabelKey)), phpT("cmd.php.packages.update.updating"))

		if err := UpdatePackages(cwd, args); err != nil {
			return core.E("php", phpT("cmd.php.error.update_packages"), err)
		}

		cli.Print(cliSectionLabelValueFormat, successStyle.Render(phpLabel("done")), phpT("cmd.php.packages.update.done"))
		return nil
	})
}

func addPHPPackagesListCommand(c *core.Core, prefix string) {
	path := phpCommandPath(prefix, "packages/list")
	phpFailureorCommand(c, path, phpT("cmd.php.packages.list.short"), func(opts core.Options) error {
		cwdResult := core.Getwd()
		if !cwdResult.OK {
			err, _ := cwdResult.Value.(error)
			return core.E("php", phpT(i18nFailGetKey, workingDirectorySubject), err)
		}
		cwd, _ := cwdResult.Value.(string)

		packages, err := ListLinkedPackages(cwd)
		if err != nil {
			return core.E("php", phpT("i18n.fail.list", "packages"), err)
		}

		if len(packages) == 0 {
			cli.Print(cliLabelValueFormat, dimStyle.Render(phpT(cmdPHPLabelKey)), phpT("cmd.php.packages.list.none_found"))
			return nil
		}

		cli.Print(cliLabelValueBlankFormat, dimStyle.Render(phpT(cmdPHPLabelKey)), phpT("cmd.php.packages.list.linked"))

		for _, pkg := range packages {
			name := pkg.Name
			if name == "" {
				name = phpT("cmd.php.packages.list.unknown")
			}
			version := pkg.Version
			if version == "" {
				version = "dev"
			}

			cli.Print("  %s %s\n", successStyle.Render("*"), name)
			cli.Print("    %s %s\n", dimStyle.Render(phpLabel(`path`)), pkg.Path)
			cli.Print("    %s %s\n", dimStyle.Render(phpLabel("version")), version)
			cli.Blank()
		}

		return nil
	})
}
