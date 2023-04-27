package fs

import (
	"os"
	"path/filepath"
)

// FileContent provides the current content of a file.
func FileContent(dir, filename string) (string, error) {
	content, err := os.ReadFile(filepath.Join(dir, filename))
	return string(content), err
}
