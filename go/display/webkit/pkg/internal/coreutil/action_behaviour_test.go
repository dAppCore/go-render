package coreutil

import core "dappco.re/go"

// DispatchAction broadcasts to every registered handler; the handler runs.
//
//	DispatchAction(c, "gui.demo", payload{ID: 1})
func TestActionBehaviour_DispatchAction_Good(t *core.T) {
	c := core.New()
	var seen any
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		seen = msg
		return core.Result{OK: true}
	})

	DispatchAction(c, "gui.demo", "payload")

	core.AssertEqual(t, "payload", seen)
}

// DispatchAction tolerates a panicking handler — broadcast recovers and the
// dispatch returns without propagating the panic.
func TestActionBehaviour_DispatchAction_Bad(t *core.T) {
	c := core.New()
	c.RegisterAction(func(_ *core.Core, _ core.Message) core.Result {
		panic("handler blew up")
	})

	core.AssertNotPanics(t, func() {
		DispatchAction(c, "gui.demo", "payload")
	})
}

// DispatchAction is a safe no-op when the Core is nil.
func TestActionBehaviour_DispatchAction_Ugly(t *core.T) {
	core.AssertNotPanics(t, func() {
		DispatchAction(nil, "gui.demo", "payload")
	})
}

// ObserveResult returns silently when the Result is OK — nothing to log.
func TestActionBehaviour_ObserveResult_Good(t *core.T) {
	c := core.New()
	core.AssertNotPanics(t, func() {
		ObserveResult(c, "gui.demo", "should not log", core.Result{OK: true})
	})
}

// ObserveResult logs the carried error when the Result failed.
func TestActionBehaviour_ObserveResult_Bad(t *core.T) {
	c := core.New()
	failed := core.Fail(core.NewError("dispatch broke"))
	core.AssertNotPanics(t, func() {
		ObserveResult(c, "gui.demo", "dispatch failed", failed)
	})
}

// ObserveResult is a safe no-op when the Core is nil even on a failed Result,
// and synthesises an error when the failed Result carries a non-error Value.
func TestActionBehaviour_ObserveResult_Ugly(t *core.T) {
	core.AssertNotPanics(t, func() {
		ObserveResult(nil, "gui.demo", "msg", core.Result{OK: false})
	})

	c := core.New()
	stringFailure := core.Result{Value: "plain string failure", OK: false}
	core.AssertNotPanics(t, func() {
		ObserveResult(c, "gui.demo", "dispatch failed", stringFailure)
	})
}

// LogWarn forwards a real error to the Core warning log.
func TestActionBehaviour_LogWarn_Good(t *core.T) {
	c := core.New()
	core.AssertNotPanics(t, func() {
		LogWarn(c, core.NewError("config.host missing"), "config.Load", "using default host")
	})
}

// LogWarn is a no-op when the error is nil — there is nothing to warn about.
func TestActionBehaviour_LogWarn_Bad(t *core.T) {
	c := core.New()
	core.AssertNotPanics(t, func() {
		LogWarn(c, nil, "config.Load", "no error")
	})
}

// LogWarn is a safe no-op when the Core is nil.
func TestActionBehaviour_LogWarn_Ugly(t *core.T) {
	core.AssertNotPanics(t, func() {
		LogWarn(nil, core.NewError("orphan warning"), "config.Load", "no core")
	})
}
