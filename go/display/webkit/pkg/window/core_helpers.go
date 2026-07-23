package window

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

func equalFold(left, right string) bool {
	return core.Lower(left) == core.Lower(right)
}
