package php

import (
	"os"
	"os/exec"

	core "dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
)

// LinkedPackage represents a linked local package.
type LinkedPackage struct {
	Name    string `json:"name"`
	Path    string
	Version string `json:"version"`
}

// composerRepository represents a composer repository entry.
type composerRepository struct {
	Type    string         `json:"type"`
	URL     string         `json:"url,omitempty"`
	Options map[string]any `json:"options,omitempty"`
}

// readComposerJSON reads and parses composer.json from the given directory.
func readComposerJSON(dir string) (map[string][]byte, error) { // Result boundary
	m := getMedium()
	composerPath := core.PathJoin(dir, composerJSONFile)
	content, err := m.Read(composerPath)
	if err != nil {
		return nil, phpWrapAction(err, "read", composerJSONFile)
	}

	var raw map[string][]byte
	if r := core.JSONUnmarshal([]byte(content), &raw); !r.OK {
		return nil, phpWrapAction(r.Value.(error), "parse", composerJSONFile)
	}

	return raw, nil
}

// writeComposerJSON writes the composer.json to the given directory.
func writeComposerJSON(dir string, raw map[string][]byte) error { // Result boundary
	m := getMedium()
	composerPath := core.PathJoin(dir, composerJSONFile)

	r := core.JSONMarshalIndent(raw, "", "    ")
	if !r.OK {
		return phpWrapAction(r.Value.(error), "marshal", composerJSONFile)
	}
	data := r.Value.([]byte)

	// Add trailing newline
	content := string(data) + "\n"

	if err := m.Write(composerPath, content); err != nil {
		return phpWrapAction(err, "write", composerJSONFile)
	}

	return nil
}

// getRepositories extracts repositories from raw composer.json.
func getRepositories(raw map[string][]byte) ([]composerRepository, error) { // Result boundary
	reposRaw, ok := raw["repositories"]
	if !ok {
		return []composerRepository{}, nil
	}

	var repos []composerRepository
	if r := core.JSONUnmarshal(reposRaw, &repos); !r.OK {
		return nil, phpWrapAction(r.Value.(error), "parse", "repositories")
	}

	return repos, nil
}

// setRepositories sets repositories in raw composer.json.
func setRepositories(raw map[string][]byte, repos []composerRepository) error { // Result boundary
	if len(repos) == 0 {
		delete(raw, "repositories")
		return nil
	}

	r := core.JSONMarshal(repos)
	if !r.OK {
		return phpWrapAction(r.Value.(error), "marshal", "repositories")
	}

	raw["repositories"] = r.Value.([]byte)
	return nil
}

// getPackageInfo reads package name and version from a composer.json in the given path.
func getPackageInfo(packagePath string) (name, version string, err error) { // Result boundary
	m := getMedium()
	composerPath := core.PathJoin(packagePath, composerJSONFile)
	content, err := m.Read(composerPath)
	if err != nil {
		return "", "", phpWrapAction(err, "read", "package composer.json")
	}

	var pkg struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	if r := core.JSONUnmarshal([]byte(content), &pkg); !r.OK {
		return "", "", phpWrapAction(r.Value.(error), "parse", "package composer.json")
	}

	if pkg.Name == "" {
		return "", "", phpFailure("package name not found in composer.json")
	}

	return pkg.Name, pkg.Version, nil
}

// LinkPackages adds path repositories to composer.json for local package development.
func LinkPackages(dir string, packages []string) error { // Result boundary
	if !IsPHPProject(dir) {
		return phpFailure(notPHPProjectComposerMessage)
	}

	raw, err := readComposerJSON(dir)
	if err != nil {
		return err
	}

	repos, err := getRepositories(raw)
	if err != nil {
		return err
	}

	for _, packagePath := range packages {
		absPath, pkgName, err := validateLinkPackage(packagePath)
		if err != nil {
			return err
		}

		if isPackageLinked(repos, absPath) {
			continue
		}

		repos = append(repos, pathComposerRepository(absPath))
		cli.Print("Linked: %s -> %s\n", pkgName, absPath)
	}

	if err := setRepositories(raw, repos); err != nil {
		return err
	}

	return writeComposerJSON(dir, raw)
}

func validateLinkPackage(packagePath string) (string, string, error) { // Result boundary
	absR := core.PathAbs(packagePath)
	if !absR.OK {
		return "", "", core.E("php", core.Sprintf("failed to resolve path %s", packagePath), absR.Value.(error))
	}
	absPath := absR.Value.(string)

	if !IsPHPProject(absPath) {
		return "", "", phpFailure("not a PHP package (missing composer.json): %s", absPath)
	}

	pkgName, _, err := getPackageInfo(absPath)
	if err != nil {
		return "", "", core.E("php", core.Sprintf("failed to get package info from %s", absPath), err)
	}

	return absPath, pkgName, nil
}

func isPackageLinked(repos []composerRepository, absPath string) bool {
	for _, repo := range repos {
		if repo.Type == `path` && repo.URL == absPath {
			return true
		}
	}
	return false
}

func pathComposerRepository(absPath string) composerRepository {
	return composerRepository{
		Type: `path`,
		URL:  absPath,
		Options: map[string]any{
			"symlink": true,
		},
	}
}

// UnlinkPackages removes path repositories from composer.json.
func UnlinkPackages(dir string, packages []string) error { // Result boundary
	if !IsPHPProject(dir) {
		return phpFailure(notPHPProjectComposerMessage)
	}

	raw, err := readComposerJSON(dir)
	if err != nil {
		return err
	}

	repos, err := getRepositories(raw)
	if err != nil {
		return err
	}

	// Build set of packages to unlink
	toUnlink := make(map[string]bool)
	for _, pkg := range packages {
		toUnlink[pkg] = true
	}

	// Filter out unlinked packages
	filtered := make([]composerRepository, 0, len(repos))
	for _, repo := range repos {
		if !shouldUnlinkRepository(repo, toUnlink) {
			filtered = append(filtered, repo)
		}
	}

	if err := setRepositories(raw, filtered); err != nil {
		return err
	}

	return writeComposerJSON(dir, raw)
}

func shouldUnlinkRepository(repo composerRepository, toUnlink map[string]bool) bool {
	if repo.Type != `path` {
		return false
	}

	shouldUnlink := false
	if IsPHPProject(repo.URL) {
		pkgName, _, err := getPackageInfo(repo.URL)
		if err == nil && toUnlink[pkgName] {
			shouldUnlink = true
			cli.Print("Unlinked: %s\n", pkgName)
		}
	}

	for pkg := range toUnlink {
		if repo.URL == pkg || core.PathBase(repo.URL) == pkg {
			shouldUnlink = true
			cli.Print("Unlinked: %s\n", repo.URL)
			break
		}
	}

	return shouldUnlink
}

// UpdatePackages runs composer update for specific packages.
func UpdatePackages(dir string, packages []string) error { // Result boundary
	if !IsPHPProject(dir) {
		return phpFailure(notPHPProjectComposerMessage)
	}

	args := []string{"update"}
	args = append(args, packages...)

	cmd := exec.Command("composer", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ListLinkedPackages returns all path repositories from composer.json.
func ListLinkedPackages(dir string) ([]LinkedPackage, error) { // Result boundary
	if !IsPHPProject(dir) {
		return nil, phpFailure(notPHPProjectComposerMessage)
	}

	raw, err := readComposerJSON(dir)
	if err != nil {
		return nil, err
	}

	repos, err := getRepositories(raw)
	if err != nil {
		return nil, err
	}

	linked := make([]LinkedPackage, 0)
	for _, repo := range repos {
		if repo.Type != `path` {
			continue
		}

		pkg := LinkedPackage{
			Path: repo.URL,
		}

		// Try to get package info
		if IsPHPProject(repo.URL) {
			name, version, err := getPackageInfo(repo.URL)
			if err == nil {
				pkg.Name = name
				pkg.Version = version
			}
		}

		if pkg.Name == "" {
			pkg.Name = core.PathBase(repo.URL)
		}

		linked = append(linked, pkg)
	}

	return linked, nil
}
