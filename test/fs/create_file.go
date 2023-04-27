package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateFile creates a file with the given name and content in this repository.
func CreateFile(dir, name, content string) error {
	filePath := filepath.Join(dir, name)
	folderPath := filepath.Dir(filePath)
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create folder %q: %w", folderPath, err)
	}
	err = os.WriteFile(filePath, []byte(content), 0o500)
	if err != nil {
		return fmt.Errorf("cannot create file %q: %w", name, err)
	}
	return nil
}
