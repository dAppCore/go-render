package preload

import core "dappco.re/go"

func coreResultError(result core.Result, fallback string) resultFailure {
	if result.OK {
		return nil
	}
	if err, ok := result.Value.(error); ok {
		return err
	}
	if text := core.Trim(result.Error()); text != "" {
		return core.NewError(text)
	}
	return core.NewError(fallback)
}

func coreOpen(path string) (*core.OSFile, resultFailure) {
	result := core.Open(path)
	if !result.OK {
		return nil, coreResultError(result, "failed to open file")
	}
	return result.Value.(*core.OSFile), nil
}

func coreStat(path string) (core.FsFileInfo, resultFailure) {
	result := core.Stat(path)
	if !result.OK {
		return nil, coreResultError(result, "failed to stat path")
	}
	return result.Value.(core.FsFileInfo), nil
}

func coreLstat(path string) (core.FsFileInfo, resultFailure) {
	result := core.Lstat(path)
	if !result.OK {
		return nil, coreResultError(result, "failed to stat path")
	}
	return result.Value.(core.FsFileInfo), nil
}

func coreReadFile(path string) ([]byte, resultFailure) {
	result := core.ReadFile(path)
	if !result.OK {
		return nil, coreResultError(result, "failed to read file")
	}
	return result.Value.([]byte), nil
}

func coreWriteFile(path string, data []byte, mode core.FileMode) resultFailure {
	return coreResultError(core.WriteFile(path, data, mode), "failed to write file")
}

func coreMkdirAll(path string, mode core.FileMode) resultFailure {
	return coreResultError(core.MkdirAll(path, mode), "failed to create directory")
}

func pathAbs(path string) (string, resultFailure) {
	result := core.PathAbs(path)
	if !result.OK {
		return "", coreResultError(result, "failed to make path absolute")
	}
	return result.Value.(string), nil
}

func pathEvalSymlinks(path string) (string, resultFailure) {
	result := core.PathEvalSymlinks(path)
	if !result.OK {
		return "", coreResultError(result, "failed to resolve symlinks")
	}
	return result.Value.(string), nil
}

func pathRel(base, target string) (string, resultFailure) {
	result := core.PathRel(base, target)
	if !result.OK {
		return "", coreResultError(result, "failed to compare paths")
	}
	return result.Value.(string), nil
}

func pathFromSlash(path string) string {
	if core.PathSeparator == '/' {
		return path
	}
	return core.Replace(path, "/", string(core.PathSeparator))
}

func trimRight(value, cutset string) string {
	for value != "" {
		trimmed := false
		for _, char := range cutset {
			part := string(char)
			if core.HasSuffix(value, part) {
				value = core.TrimSuffix(value, part)
				trimmed = true
			}
		}
		if !trimmed {
			return value
		}
	}
	return value
}

func trimRunes(value, cutset string) string {
	for value != "" {
		trimmed := false
		for _, char := range cutset {
			part := string(char)
			if core.HasPrefix(value, part) {
				value = core.TrimPrefix(value, part)
				trimmed = true
			}
			if core.HasSuffix(value, part) {
				value = core.TrimSuffix(value, part)
				trimmed = true
			}
		}
		if !trimmed {
			return value
		}
	}
	return value
}
