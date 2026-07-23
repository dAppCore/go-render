//go:build windows

package php

import (
	"syscall"

	core "dappco.re/go"
)

// setSysProcAttr sets Windows-specific process attributes.
// Windows doesn't support Setpgid, so this is a no-op.
func setSysProcAttr(cmd *core.Cmd) {
	// No-op on Windows - process groups work differently
}

// signalProcessGroup sends a termination signal to the process.
// On Windows, we can only signal the main process, not a group.
func signalProcessGroup(cmd *core.Cmd, sig syscall.Signal) error { // Result boundary
	if cmd.Process == nil {
		return nil
	}

	return cmd.Process.Signal(sig)
}

// termSignal returns SIGINT for Windows (closest to SIGTERM).
func termSignal() syscall.Signal {
	return syscall.SIGINT
}

// killSignal returns SIGKILL for Windows.
func killSignal() syscall.Signal {
	return syscall.SIGKILL
}
