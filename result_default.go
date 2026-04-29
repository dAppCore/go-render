//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go"

type result = core.Result

func okResult(value any) result {
	return core.Ok(value)
}
