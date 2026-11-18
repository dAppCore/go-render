// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "github.com/dappcore/go"

func ExampleRender() {
	core.Println(Render(El("p", Text("hello")), NewContext()))
	// Output: <p>hello</p>
}
