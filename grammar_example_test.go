// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go"

func ExampleGrammarImprint_Imprint() {
	stamp := (&GrammarImprint{}).Imprint(El("section", Text("body")), *NewContext())
	core.Println(stamp.Path, stamp.Tags[0])
	// Output: 0 branch
}
