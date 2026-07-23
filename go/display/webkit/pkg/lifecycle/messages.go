// pkg/lifecycle/messages.go
package lifecycle

// Lifecycle events are broadcasts (Actions); outgoing commands are Tasks.

// ActionApplicationStarted fires when the platform application starts.
// Distinct from core.ActionServiceStartup — this is platform-level readiness.
type ActionApplicationStarted struct{}

// ActionOpenedWithFile fires when the application is opened with a file argument.
type ActionOpenedWithFile struct{ Path string }

// ActionLaunchedWithUrl fires when the application is launched via a URL
// scheme handoff (e.g. lthn:// on macOS, ms-app:// on Windows). The URL
// arrives verbatim; consumers route to surfaces / handlers via their own
// scheme parser.
type ActionLaunchedWithUrl struct{ URL string }

// ActionWillTerminate fires when the application is about to terminate (macOS only).
type ActionWillTerminate struct{}

// ActionDidBecomeActive fires when the application becomes the active app (macOS only).
type ActionDidBecomeActive struct{}

// ActionDidResignActive fires when the application resigns active status (macOS only).
type ActionDidResignActive struct{}

// ActionPowerStatusChanged fires on power status changes (Windows only: APMPowerStatusChange).
type ActionPowerStatusChanged struct{}

// ActionSystemSuspend fires when the system is about to suspend (Windows only: APMSuspend).
type ActionSystemSuspend struct{}

// ActionSystemResume fires when the system resumes from suspend (Windows only: APMResume).
type ActionSystemResume struct{}

// TaskQuit asks the platform to quit the application. Consumers dispatch via
// `c.Action("lifecycle.quit").Run(ctx, NewOptions(Option{Key: "task", Value: TaskQuit{}}))`.
type TaskQuit struct{}
