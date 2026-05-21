package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	canonicalLine, err := readCanonicalRTAVersionLine("Makefile")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	walkErr := filepath.WalkDir(".", func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || entry.Name() != "Makefile" || path == "Makefile" {
			return nil
		}
		info, statErr := entry.Info()
		if statErr != nil {
			return statErr
		}
		rawContent, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		newContent, modified := replaceRTAVersionAssignment(string(rawContent), canonicalLine)
		if !modified {
			return nil
		}
		return os.WriteFile(path, []byte(newContent), info.Mode())
	})
	if walkErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", walkErr)
		os.Exit(1)
	}
}
