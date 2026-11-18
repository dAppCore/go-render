// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "github.com/dappcore/go"

func ExampleParseBlockID() {
	core.Println(string(ParseBlockID("C.0.H.1")))
	// Output: CH
}
