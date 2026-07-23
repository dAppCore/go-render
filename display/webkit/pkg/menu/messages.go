// pkg/menu/messages.go
package menu

type QueryConfig struct{}

type QueryGetAppMenu struct{}

type TaskSetAppMenu struct{ Items []MenuItem }

type TaskSaveConfig struct{ Config map[string]any }
