// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go"

func ExampleRender() {
	core.Println(Render(El("p", Text("hello")), NewContext()))
	// Output: <p>hello</p>
}
