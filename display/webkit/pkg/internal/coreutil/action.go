package coreutil

import core "dappco.re/go"

func DispatchAction(c *core.Core, operation string, message any) {
	if c == nil {
		return
	}
	ObserveResult(c, operation, "action dispatch failed", c.ACTION(message))
}

func ObserveResult(c *core.Core, operation, message string, result core.Result) {
	if result.OK || c == nil {
		return
	}
	err, ok := result.Value.(error)
	if !ok {
		err = core.NewError(result.Error())
	}
	if warn := c.LogWarn(err, operation, message); !warn.OK {
		return
	}
}

func LogWarn(c *core.Core, err error, operation, message string) {
	if c == nil || err == nil {
		return
	}
	if warn := c.LogWarn(err, operation, message); !warn.OK {
		return
	}
}
