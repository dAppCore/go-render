// pkg/dialog/service_behaviour_test.go
package dialog

import core "dappco.re/go"

// hasDirectDialogOptions reports true when any option other than "task" is set.
//
//	hasDirectDialogOptions(core.NewOptions(core.Option{Key: "title", Value: "Hi"})) // true
func TestServiceBehaviour_hasDirectDialogOptions_Good(t *core.T) {
	core.AssertTrue(t, hasDirectDialogOptions(core.NewOptions(
		core.Option{Key: "title", Value: "Hi"},
	)))
}

// hasDirectDialogOptions ignores a sole "task" option.
func TestServiceBehaviour_hasDirectDialogOptions_Bad(t *core.T) {
	core.AssertFalse(t, hasDirectDialogOptions(core.NewOptions(
		core.Option{Key: "task", Value: TaskInfo{Title: "x"}},
	)))
}

// hasDirectDialogOptions reports false for empty options.
func TestServiceBehaviour_hasDirectDialogOptions_Ugly(t *core.T) {
	core.AssertFalse(t, hasDirectDialogOptions(core.NewOptions()))
}

// typedMessageDialogOptionsFrom decodes direct fields and stamps the type.
//
//	typedMessageDialogOptionsFrom(opts, DialogInfo, "op")
func TestServiceBehaviour_typedMessageDialogOptionsFrom_Good(t *core.T) {
	got, err := typedMessageDialogOptionsFrom(core.NewOptions(
		core.Option{Key: "title", Value: "Heads up"},
		core.Option{Key: "message", Value: "Check this"},
	), DialogWarning, "dialog.test")
	core.AssertNil(t, err)
	core.AssertEqual(t, DialogWarning, got.Type)
	core.AssertEqual(t, "Heads up", got.Title)
}

// typedMessageDialogOptionsFrom fails when no direct options are present.
func TestServiceBehaviour_typedMessageDialogOptionsFrom_Bad(t *core.T) {
	_, err := typedMessageDialogOptionsFrom(core.NewOptions(), DialogInfo, "dialog.test")
	core.AssertNotNil(t, err)
	core.AssertContains(t, err.Error(), "failed to decode")
}

// infoDialogOptionsFrom resolves a TaskInfo, a typed MessageDialogOptions, and
// direct field options — stamping DialogInfo in every case.
func TestServiceBehaviour_infoDialogOptionsFrom_Good(t *core.T) {
	fromTask, err := infoDialogOptionsFrom(core.NewOptions(
		core.Option{Key: "task", Value: TaskInfo{Title: "T", Message: "M"}},
	))
	core.AssertNil(t, err)
	core.AssertEqual(t, DialogInfo, fromTask.Type)
	core.AssertEqual(t, "T", fromTask.Title)

	fromOpts, err := infoDialogOptionsFrom(core.NewOptions(
		core.Option{Key: "task", Value: MessageDialogOptions{Title: "Opt"}},
	))
	core.AssertNil(t, err)
	core.AssertEqual(t, DialogInfo, fromOpts.Type)

	fromDirect, err := infoDialogOptionsFrom(core.NewOptions(
		core.Option{Key: "title", Value: "Direct"},
	))
	core.AssertNil(t, err)
	core.AssertEqual(t, DialogInfo, fromDirect.Type)
	core.AssertEqual(t, "Direct", fromDirect.Title)
}

// infoDialogOptionsFrom fails on an unrecognised task type.
func TestServiceBehaviour_infoDialogOptionsFrom_Bad(t *core.T) {
	_, err := infoDialogOptionsFrom(core.NewOptions(
		core.Option{Key: "task", Value: 42},
	))
	core.AssertNotNil(t, err)
}

// warningDialogOptionsFrom and errorDialogOptionsFrom stamp their own types.
func TestServiceBehaviour_warningErrorDialogOptionsFrom_Good(t *core.T) {
	warn, err := warningDialogOptionsFrom(core.NewOptions(
		core.Option{Key: "task", Value: TaskWarning{Title: "W"}},
	))
	core.AssertNil(t, err)
	core.AssertEqual(t, DialogWarning, warn.Type)

	derr, err := errorDialogOptionsFrom(core.NewOptions(
		core.Option{Key: "task", Value: TaskError{Title: "E"}},
	))
	core.AssertNil(t, err)
	core.AssertEqual(t, DialogError, derr.Type)
}

// questionDialogOptionsFrom resolves a TaskQuestion and stamps DialogQuestion.
func TestServiceBehaviour_questionDialogOptionsFrom_Good(t *core.T) {
	got, err := questionDialogOptionsFrom(core.NewOptions(
		core.Option{Key: "task", Value: TaskQuestion{Title: "Q", Message: "Delete?", Buttons: []string{"Yes", "No"}}},
	))
	core.AssertNil(t, err)
	core.AssertEqual(t, DialogQuestion, got.Type)
	core.AssertLen(t, got.Buttons, 2)
}

// openDirectoryOptionsFrom resolves both the TaskOpenDirectory and the bare
// OpenDirectoryOptions task shapes.
func TestServiceBehaviour_openDirectoryOptionsFrom_Good(t *core.T) {
	fromTask, err := openDirectoryOptionsFrom(core.NewOptions(
		core.Option{Key: "task", Value: TaskOpenDirectory{Options: OpenDirectoryOptions{Title: "Pick"}}},
	))
	core.AssertNil(t, err)
	core.AssertEqual(t, "Pick", fromTask.Title)

	fromOpts, err := openDirectoryOptionsFrom(core.NewOptions(
		core.Option{Key: "task", Value: OpenDirectoryOptions{Directory: "/tmp"}},
	))
	core.AssertNil(t, err)
	core.AssertEqual(t, "/tmp", fromOpts.Directory)
}

// messageDialogOptionsFrom resolves the TaskMessageDialog and direct shapes.
func TestServiceBehaviour_messageDialogOptionsFrom_Good(t *core.T) {
	got, err := messageDialogOptionsFrom(core.NewOptions(
		core.Option{Key: "task", Value: MessageDialogOptions{Title: "M", Type: DialogWarning}},
	))
	core.AssertNil(t, err)
	core.AssertEqual(t, "M", got.Title)
}
