//go:build compliance

package window

import core "dappco.re/go"

func ExampleNewLayoutManager() {
	core.Println("NewLayoutManager")
	// Output:
	// NewLayoutManager
}

func ExampleNewLayoutManagerWithDir() {
	core.Println("NewLayoutManagerWithDir")
	// Output:
	// NewLayoutManagerWithDir
}

func ExampleNewLayoutManagerWithPath() {
	core.Println("NewLayoutManagerWithPath")
	// Output:
	// NewLayoutManagerWithPath
}

func ExampleLayoutManager_SetPath() {
	core.Println("LayoutManager_SetPath")
	// Output:
	// LayoutManager_SetPath
}

func ExampleLayoutManager_SaveLayout() {
	core.Println("LayoutManager_SaveLayout")
	// Output:
	// LayoutManager_SaveLayout
}

func ExampleLayoutManager_GetLayout() {
	core.Println("LayoutManager_GetLayout")
	// Output:
	// LayoutManager_GetLayout
}

func ExampleLayoutManager_ListLayouts() {
	core.Println("LayoutManager_ListLayouts")
	// Output:
	// LayoutManager_ListLayouts
}

func ExampleLayoutManager_DeleteLayout() {
	core.Println("LayoutManager_DeleteLayout")
	// Output:
	// LayoutManager_DeleteLayout
}
