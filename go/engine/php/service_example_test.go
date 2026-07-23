package php

func ExampleNewService() {
	_ = NewService
}

func ExampleRegister() {
	_ = Register
}

func ExampleCoreService_OnStartup() {
	_ = (*CoreService).OnStartup
}

func ExampleCoreService_OnShutdown() {
	_ = (*CoreService).OnShutdown
}
