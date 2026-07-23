package deno

import (
	"context"

	core "dappco.re/go"
)

func commandContext(ctx context.Context, binary string, args ...string) *core.Cmd {
	cmd := &core.Cmd{Path: binary, Args: append([]string{binary}, args...)}
	if ctx != nil {
		go func() {
			<-ctx.Done()
			if cmd.Process != nil {
				if err := cmd.Process.Kill(); err != nil && !isProcessDone(err) {
					core.Error("failed to kill deno command", "err", err)
				}
			}
		}()
	}
	return cmd
}

func isProcessDone(err error) bool {
	if err == nil {
		return false
	}
	text := core.Lower(err.Error())
	return core.Contains(text, "process already finished") ||
		core.Contains(text, "no child processes") ||
		core.Contains(text, "os: process already finished")
}
