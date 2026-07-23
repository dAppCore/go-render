package php

import (
	"io/fs"

	core "dappco.re/go"
)

// Extract copies an embedded Laravel app (from embed.FS) to a temporary directory.
// FrankenPHP needs real filesystem paths — it cannot serve from embed.FS.
// The prefix is the embed directory name (e.g. "laravel").
// Returns the path to the extracted Laravel root. Caller must core.RemoveAll on cleanup.
func Extract(fsys fs.FS, prefix string) (string, error) { // Result boundary
	tmpResult := core.MkdirTemp("", "go-php-laravel-*")
	if !tmpResult.OK {
		err, _ := tmpResult.Value.(error)
		return "", core.Errorf("create temp dir: %w", err)
	}
	tmpDir, _ := tmpResult.Value.(string)

	err := fs.WalkDir(fsys, prefix, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relResult := core.PathRel(prefix, path)
		if !relResult.OK {
			relErr, _ := relResult.Value.(error)
			return relErr
		}
		relPath, _ := relResult.Value.(string)
		targetPath := core.PathJoin(tmpDir, relPath)

		if d.IsDir() {
			r := core.MkdirAll(targetPath, 0o755)
			if !r.OK {
				mkErr, _ := r.Value.(error)
				return mkErr
			}
			return nil
		}

		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return core.Errorf("read embedded %s: %w", path, err)
		}

		r := core.WriteFile(targetPath, data, 0o644)
		if !r.OK {
			wrErr, _ := r.Value.(error)
			return wrErr
		}
		return nil
	})

	if err != nil {
		if cleanupResult := core.RemoveAll(tmpDir); !cleanupResult.OK {
			cleanupErr, _ := cleanupResult.Value.(error)
			return "", cleanupErr
		}
		return "", core.Errorf("extract Laravel: %w", err)
	}

	return tmpDir, nil
}
