package php

import (
	core "dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
)

var phpCommandValueFlags = map[string]bool{
	"container":  true,
	"dockerfile": true,
	"domain":     true,
	"env-file":   true,
	"fail-on":    true,
	"format":     true,
	"https-port": true,
	"id":         true,
	"limit":      true,
	"name":       true,
	"output":     true,
	`path`:       true,
	"platform":   true,
	"port":       true,
	"service":    true,
	"tag":        true,
	"template":   true,
	"threads":    true,
	"type":       true,
	"workers":    true,
}

type phpCommandLine struct {
	flags map[string]string
	args  []string
}

func phpFailure(format string, args ...any) error {
	return core.E("php", core.Sprintf(format, args...), nil)
}

func phpWrapMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return core.E("php", message, err)
}

func phpWrapAction(err error, verb, subject string) error {
	if err == nil {
		return nil
	}
	return core.E("php", core.Sprintf("failed to %s %s", verb, subject), err)
}

func phpExit(code int, err error) error { // Result boundary
	if err == nil {
		return nil
	}
	return &cli.ExitError{Code: code, Err: err}
}

func phpFailureorResult(err error) core.Result {
	if err != nil {
		return core.Fail(err)
	}
	return core.Ok(nil)
}

func phpCommand(c *core.Core, path, description string, action core.CommandAction) {
	c.Command(path, core.Command{
		Description: description,
		Action: func(opts core.Options) core.Result {
			if err := activateWorkspacePackage(); err != nil {
				return core.Fail(err)
			}
			return action(opts)
		},
	})
}

func phpFailureorCommand(c *core.Core, path, description string, action func(core.Options) error) { // Result boundary
	phpCommand(c, path, description, func(opts core.Options) core.Result {
		return phpFailureorResult(action(opts))
	})
}

func phpHelpCommand(c *core.Core, path, description string) {
	phpCommand(c, path, description, func(core.Options) core.Result {
		if cl := c.Cli(); cl != nil {
			cl.PrintHelp()
		}
		return core.Ok(nil)
	})
}

func phpCommandLineFor(path string, opts core.Options) phpCommandLine {
	if raw, ok := phpRawArgsForCommand(path); ok {
		return phpParseCommandLine(raw)
	}
	return phpCommandLineFromOptions(opts)
}

func phpRawArgsForCommand(path string) ([]string, bool) {
	parts := core.Split(path, "/")
	if len(parts) == 0 {
		return nil, false
	}

	args := core.FilterArgs(core.Args()[1:])
	if len(args) < len(parts) {
		return nil, false
	}
	for i, part := range parts {
		if args[i] != part {
			return nil, false
		}
	}
	return args[len(parts):], true
}

func phpCommandLineFromOptions(opts core.Options) phpCommandLine {
	line := phpCommandLine{flags: map[string]string{}}
	for _, item := range opts.Items() {
		switch item.Key {
		case "_arg":
			if value, ok := item.Value.(string); ok && value != "" {
				line.args = []string{value}
			}
		case "_args":
			if values, ok := item.Value.([]string); ok {
				line.args = append([]string(nil), values...)
			}
		default:
			switch value := item.Value.(type) {
			case string:
				line.flags[item.Key] = value
			case bool:
				if value {
					line.flags[item.Key] = "true"
				}
			case int:
				line.flags[item.Key] = core.Itoa(value)
			}
		}
	}
	return line
}

func phpParseCommandLine(args []string) phpCommandLine {
	line := phpCommandLine{flags: map[string]string{}}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			line.args = append(line.args, args[i+1:]...)
			break
		}

		key, value, hasValue, ok := phpSplitFlag(arg)
		if ok {
			if !hasValue && phpCommandValueFlags[key] && i+1 < len(args) && !core.IsFlag(args[i+1]) {
				value = args[i+1]
				hasValue = true
				i++
			}
			if !hasValue {
				value = "true"
			}
			line.flags[key] = value
			continue
		}

		line.args = append(line.args, arg)
	}
	return line
}

func phpSplitFlag(arg string) (key, value string, hasValue bool, ok bool) {
	if core.HasPrefix(arg, "--") {
		body := core.TrimPrefix(arg, "--")
		if body == "" {
			return "", "", false, false
		}
		key, value, hasValue = phpCutOnEquals(body)
		return key, value, hasValue, key != ""
	}

	if core.HasPrefix(arg, "-") {
		body := core.TrimPrefix(arg, "-")
		if body == "" {
			return "", "", false, false
		}
		key, value, hasValue = phpCutOnEquals(body)
		return key, value, hasValue, key != ""
	}

	return "", "", false, false
}

// phpCutOnEquals splits "key=value" into (key, value, true) or returns
// (s, "", false) when no `=` is present. Equivalent of strings.Cut(s, "=")
// without importing strings.
func phpCutOnEquals(s string) (key, value string, hasValue bool) {
	parts := core.SplitN(s, "=", 2)
	if len(parts) == 2 {
		return parts[0], parts[1], true
	}
	return s, "", false
}

func (line phpCommandLine) Bool(name string, aliases ...string) bool {
	keys := append([]string{name}, aliases...)
	for _, key := range keys {
		value, ok := line.flags[key]
		if !ok {
			continue
		}
		switch core.Lower(value) {
		case "", "1", "true", "yes", "on":
			return true
		default:
			return false
		}
	}
	return false
}

func (line phpCommandLine) String(name string, fallback string) string {
	if value, ok := line.flags[name]; ok {
		return value
	}
	return fallback
}

func (line phpCommandLine) Int(name string, fallback int) int {
	value, ok := line.flags[name]
	if !ok || value == "" {
		return fallback
	}
	if parsed := core.Atoi(value); parsed.OK {
		return parsed.Value.(int)
	}
	return fallback
}

func (line phpCommandLine) Args() []string {
	return append([]string(nil), line.args...)
}
