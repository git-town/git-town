package git

import (
	"fmt"
	"os"
	"path/filepath"
)

// RemoveUnnecessaryFiles trims all files that aren't necessary in this repo.
func RemoveUnnecessaryFiles(dir string) error {
	fullPath := filepath.Join(dir, ".git", "hooks")
	err := os.RemoveAll(fullPath)
	if err != nil {
		return fmt.Errorf("cannot remove unnecessary files in %q: %w", fullPath, err)
	}
	_ = os.Remove(filepath.Join(dir, ".git", "COMMIT_EDITMSG"))
	_ = os.Remove(filepath.Join(dir, ".git", "description"))
	return nil
}
