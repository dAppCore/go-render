package lifecycle

import (
	core "dappco.re/go"
)

type AppMode string

const (
	ModeManager AppMode = "manager"
	ModeWorker  AppMode = "worker"
)

const (
	appModeEnv = "CORE_APP_MODE"
	ciEnv      = "CI"
)

// DetectMode returns ModeWorker when CORE_APP_MODE=worker.
//
//	mode := DetectMode()
func DetectMode() AppMode {
	if value, ok := core.LookupEnv(appModeEnv); ok {
		if mode, valid := parseAppMode(value); valid {
			return mode
		}
	}

	if value, ok := core.LookupEnv(ciEnv); ok && isTrue(value) {
		return ModeWorker
	}

	return ModeManager
}

func parseAppMode(value string) (AppMode, bool) {
	switch core.Lower(core.Trim(value)) {
	case string(ModeManager):
		return ModeManager, true
	case string(ModeWorker):
		return ModeWorker, true
	case "":
		return "", false
	default:
		return ModeManager, true
	}
}

func isTrue(value string) bool {
	return core.Lower(core.Trim(value)) == "true"
}
