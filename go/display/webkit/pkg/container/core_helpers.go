package container

import (
	"context"

	core "dappco.re/go"
)

func lookPath(file string) (string, resultFailure) {
	name := core.Trim(file)
	if name == "" {
		return "", core.NewError("executable name is empty")
	}
	if core.Contains(name, string(core.PathSeparator)) || core.PathIsAbs(name) {
		if executablePath(name) {
			return name, nil
		}
		return "", core.Errorf("executable file not found in path: %s", file)
	}
	for _, dir := range core.Split(core.Getenv("PATH"), string(core.PathListSeparator)) {
		if core.Trim(dir) == "" {
			continue
		}
		candidate := core.PathJoin(dir, name)
		if executablePath(candidate) {
			return candidate, nil
		}
	}
	return "", core.Errorf("executable file not found in path: %s", file)
}

func executablePath(path string) bool {
	result := core.Stat(path)
	if !result.OK {
		return false
	}
	info := result.Value.(core.FsFileInfo)
	return !info.IsDir() && info.Mode().Perm()&0o111 != 0
}

func commandContext(ctx context.Context, binary string, args ...string) *core.Cmd {
	cmd := &core.Cmd{Path: binary, Args: append([]string{binary}, args...)}
	if ctx != nil {
		go func() {
			<-ctx.Done()
			if cmd.Process != nil {
				if err := cmd.Process.Kill(); err != nil {
					core.Error("failed to kill container command", "err", err)
				}
			}
		}()
	}
	return cmd
}

func command(binary string, args ...string) *core.Cmd {
	return commandContext(nil, binary, args...)
}

func coreWriteMode(path, content string, mode core.FileMode) resultFailure {
	result := core.WriteFile(path, []byte(content), mode)
	if result.OK {
		return nil
	}
	if err, ok := result.Value.(error); ok {
		return err
	}
	return core.NewError(result.Error())
}

func cut(value, sep string) (string, string, bool) {
	parts := core.SplitN(value, sep, 2)
	if len(parts) != 2 {
		return value, "", false
	}
	return parts[0], parts[1], true
}

func cutPrefix(value, prefix string) (string, bool) {
	if !core.HasPrefix(value, prefix) {
		return value, false
	}
	return core.TrimPrefix(value, prefix), true
}
