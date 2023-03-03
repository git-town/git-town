package runstate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/git-town/git-town/v7/src/git"
)

// Load loads the run state for the given Git repo from disk. Can return nil if there is no saved runstate.
func Load(repo *git.ProdRepo) (*RunState, error) {
	filename, err := PersistenceFilePath(repo)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil //nolint:nilnil
		}
		return nil, fmt.Errorf("cannot check file %q: %w", filename, err)
	}
	var runState RunState
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %q: %w", filename, err)
	}
	err = json.Unmarshal(content, &runState)
	if err != nil {
		return nil, fmt.Errorf("cannot parse content of file %q: %w", filename, err)
	}
	return &runState, nil
}

// Delete removes the stored run state from disk.
func Delete(repo *git.ProdRepo) error {
	filename, err := PersistenceFilePath(repo)
	if err != nil {
		return err
	}
	_, err = os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("cannot check file %q: %w", filename, err)
	}
	err = os.Remove(filename)
	if err != nil {
		return fmt.Errorf("cannot delete file %q: %w", filename, err)
	}
	return nil
}

// Save stores the given run state for the given Git repo to disk.
func Save(runState *RunState, repo *git.ProdRepo) error {
	content, err := json.MarshalIndent(runState, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot encode run-state: %w", err)
	}
	persistencePath, err := PersistenceFilePath(repo)
	if err != nil {
		return err
	}
	persistenceDir := filepath.Dir(persistencePath)
	err = os.MkdirAll(persistenceDir, 0o700)
	if err != nil {
		return err
	}
	err = os.WriteFile(persistencePath, content, 0o600)
	if err != nil {
		return fmt.Errorf("cannot write file %q: %w", persistencePath, err)
	}
	return nil
}

func PersistenceFilePath(repo *git.ProdRepo) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	persistenceDir := filepath.Join(configDir, "git-town", "runstate")
	repoDir, err := repo.Silent.RootDirectory()
	if err != nil {
		return "", err
	}
	filename := SanitizePath(repoDir)
	if err != nil {
		return "", err
	}
	return filepath.Join(persistenceDir, filename+".json"), nil
}

func SanitizePath(dir string) string {
	replaceCharacterRE := regexp.MustCompile("[[:^alnum:]]")
	sanitized := replaceCharacterRE.ReplaceAllString(dir, "-")
	sanitized = strings.ToLower(sanitized)
	replaceDoubleMinusRE := regexp.MustCompile("--+") // two or more dashes
	sanitized = replaceDoubleMinusRE.ReplaceAllString(sanitized, "-")
	for strings.HasPrefix(sanitized, "-") {
		sanitized = sanitized[1:]
	}
	return sanitized
}
