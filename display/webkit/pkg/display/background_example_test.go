//go:build compliance

package display

import core "dappco.re/go"

func ExampleNewBackgroundRegistry() {
	core.Println("NewBackgroundRegistry")
	// Output:
	// NewBackgroundRegistry
}

func ExampleBackgroundRegistry_RegisterServiceWorker() {
	core.Println("BackgroundRegistry_RegisterServiceWorker")
	// Output:
	// BackgroundRegistry_RegisterServiceWorker
}

func ExampleBackgroundRegistry_AddFetch() {
	core.Println("BackgroundRegistry_AddFetch")
	// Output:
	// BackgroundRegistry_AddFetch
}

func ExampleBackgroundRegistry_AddSync() {
	core.Println("BackgroundRegistry_AddSync")
	// Output:
	// BackgroundRegistry_AddSync
}

func ExampleBackgroundRegistry_AddPush() {
	core.Println("BackgroundRegistry_AddPush")
	// Output:
	// BackgroundRegistry_AddPush
}

func ExampleBackgroundRegistry_SyncRegistrationsCount() {
	core.Println("BackgroundRegistry_SyncRegistrationsCount")
	// Output:
	// BackgroundRegistry_SyncRegistrationsCount
}

func ExampleBackgroundRegistry_PushSubscriptionsCount() {
	core.Println("BackgroundRegistry_PushSubscriptionsCount")
	// Output:
	// BackgroundRegistry_PushSubscriptionsCount
}

func ExampleBackgroundRegistry_SetPaymentInstrument() {
	core.Println("BackgroundRegistry_SetPaymentInstrument")
	// Output:
	// BackgroundRegistry_SetPaymentInstrument
}
