package php

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"runtime"

	core "dappco.re/go"
)

// RuntimeEnvironment holds the resolved paths for the running application.
type RuntimeEnvironment struct {
	// DataDir is the persistent data directory (survives app updates).
	DataDir string
	// LaravelRoot is the extracted Laravel app in the temp directory.
	LaravelRoot string
	// DatabasePath is the full path to the SQLite database file.
	DatabasePath string
}

// PrepareRuntimeEnvironment creates data directories, generates .env, and symlinks
// storage so Laravel can write to persistent locations.
// The appName is used for the data directory name (e.g. "bugseti").
func PrepareRuntimeEnvironment(laravelRoot, appName string) (*RuntimeEnvironment, error) { // Result boundary
	dataDir, err := resolveDataDir(appName)
	if err != nil {
		return nil, core.E("php.PrepareRuntimeEnvironment", "resolve data dir", err)
	}

	env := &RuntimeEnvironment{
		DataDir:      dataDir,
		LaravelRoot:  laravelRoot,
		DatabasePath: core.PathJoin(dataDir, appName+".sqlite"),
	}

	// Create persistent directories
	dirs := []string{
		dataDir,
		core.PathJoin(dataDir, "storage", "app"),
		core.PathJoin(dataDir, "storage", "framework", "cache", "data"),
		core.PathJoin(dataDir, "storage", "framework", "sessions"),
		core.PathJoin(dataDir, "storage", "framework", "views"),
		core.PathJoin(dataDir, "storage", "logs"),
	}
	for _, dir := range dirs {
		if r := core.MkdirAll(dir, 0o755); !r.OK {
			return nil, core.E("php.PrepareRuntimeEnvironment", core.Sprintf("create dir %s", dir), r.Value.(error))
		}
	}

	// Create empty SQLite database if it doesn't exist
	if r := core.Stat(env.DatabasePath); !r.OK && core.IsNotExist(r.Value.(error)) {
		if r := core.WriteFile(env.DatabasePath, nil, 0o644); !r.OK {
			return nil, core.E("php.PrepareRuntimeEnvironment", "create database", r.Value.(error))
		}
		core.Println(core.Sprintf("go-php: created new database: %s", env.DatabasePath))
	}

	// Replace the extracted storage/ with a symlink to the persistent one
	extractedStorage := core.PathJoin(laravelRoot, "storage")
	if r := core.RemoveAll(extractedStorage); !r.OK {
		return nil, core.E("php.PrepareRuntimeEnvironment", "remove extracted storage", r.Value.(error))
	}
	persistentStorage := core.PathJoin(dataDir, "storage")
	// os.Symlink is retained — no core.Symlink wrapper exists in dappco.re/go v0.9.0.
	// Surface as wrapper gap when bumping core/go module dep.
	if err := os.Symlink(persistentStorage, extractedStorage); err != nil {
		return nil, core.E("php.PrepareRuntimeEnvironment", "symlink storage", err)
	}

	// Generate .env file with resolved paths
	if err := writeEnvFile(laravelRoot, appName, env); err != nil {
		return nil, core.E("php.PrepareRuntimeEnvironment", "write .env", err)
	}

	return env, nil
}

// AppendEnv appends a key=value pair to the Laravel .env file.
func AppendEnv(laravelRoot, key, value string) error { // Result boundary
	envFile := core.PathJoin(laravelRoot, ".env")
	r := core.OpenFile(envFile, core.O_APPEND|core.O_WRONLY, 0o644)
	if !r.OK {
		return r.Value.(error)
	}
	f := r.Value.(*core.OSFile)
	defer func() { _ = f.Close() }()
	_, err := f.WriteString(core.Sprintf("%s=\"%s\"\n", key, value))
	return err
}

func resolveDataDir(appName string) (string, error) { // Result boundary
	var base string
	switch runtime.GOOS {
	case "darwin":
		homeR := core.UserHomeDir()
		if !homeR.OK {
			return "", homeR.Value.(error)
		}
		base = core.PathJoin(homeR.Value.(string), "Library", "Application Support", appName)
	case "linux":
		if xdg := core.Getenv("XDG_DATA_HOME"); xdg != "" {
			base = core.PathJoin(xdg, appName)
		} else {
			homeR := core.UserHomeDir()
			if !homeR.OK {
				return "", homeR.Value.(error)
			}
			base = core.PathJoin(homeR.Value.(string), ".local", "share", appName)
		}
	default:
		homeR := core.UserHomeDir()
		if !homeR.OK {
			return "", homeR.Value.(error)
		}
		base = core.PathJoin(homeR.Value.(string), "."+appName)
	}
	return base, nil
}

func writeEnvFile(laravelRoot, appName string, env *RuntimeEnvironment) error { // Result boundary
	appKey, err := loadOrGenerateAppKey(env.DataDir)
	if err != nil {
		return core.E("php.writeEnvFile", "app key", err)
	}

	content := core.Sprintf(`APP_NAME="%s"
APP_ENV=production
APP_KEY=%s
APP_DEBUG=false
APP_URL=http://localhost

DB_CONNECTION=sqlite
DB_DATABASE="%s"

CACHE_STORE=file
SESSION_DRIVER=file
LOG_CHANNEL=single
LOG_LEVEL=warning

`, appName, appKey, env.DatabasePath)

	if r := core.WriteFile(core.PathJoin(laravelRoot, ".env"), []byte(content), 0o644); !r.OK {
		return r.Value.(error)
	}
	return nil
}

func loadOrGenerateAppKey(dataDir string) (string, error) { // Result boundary
	keyFile := core.PathJoin(dataDir, ".app-key")

	if r := core.ReadFile(keyFile); r.OK {
		data := r.Value.([]byte)
		if len(data) > 0 {
			return string(data), nil
		}
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", core.E("php.loadOrGenerateAppKey", "generate key", err)
	}
	appKey := "base64:" + base64.StdEncoding.EncodeToString(key)

	if r := core.WriteFile(keyFile, []byte(appKey), 0o600); !r.OK {
		return "", core.E("php.loadOrGenerateAppKey", "save key", r.Value.(error))
	}

	core.Println(core.Sprintf("go-php: generated new APP_KEY (saved to %s)", keyFile))
	return appKey, nil
}
