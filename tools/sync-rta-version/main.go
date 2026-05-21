package main

import (
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v23/pkg/asserts"
)

func main() {
	canonicalLine := asserts.NoError1(readCanonicalRTAVersionLine("Makefile"))
	walkErr := filepath.WalkDir(".", func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || entry.Name() != "Makefile" || path == "Makefile" {
			return nil
		}
		info := asserts.NoError1(entry.Info())
		rawContent := asserts.NoError1(os.ReadFile(path))
		newContent, modified := replaceRTAVersionAssignment(string(rawContent), canonicalLine)
		if !modified {
			return nil
		}
		return os.WriteFile(path, []byte(newContent), info.Mode())
	})
	asserts.NoError(walkErr)
}
