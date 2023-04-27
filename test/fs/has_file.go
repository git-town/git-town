package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// HasFile indicates whether this repository contains a file with the given name and content.
func HasFile(dir, name, content string) (bool, error) {
	rawContent, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return false, fmt.Errorf("repo doesn't have file %q: %w", name, err)
	}
	actualContent := string(rawContent)
	if actualContent != content {
		return false, fmt.Errorf("file %q should have content %q but has %q", name, content, actualContent)
	}
	return true, nil
}
