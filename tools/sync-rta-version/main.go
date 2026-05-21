package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v23/pkg/asserts"
)

func main() {
	canonicalLine, hasRTAVersion := readCanonicalRTAVersionLine("Makefile").Get()
	if !hasRTAVersion {
		fmt.Println("No RTA_VERSION assignment found in Makefile")
		return
	}
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
