package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/git-town/git-town/v23/pkg/asserts"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

const RTAVersionDeclaration = "RTA_VERSION = "

func main() {
	canonicalLine, hasRTAVersion := readCanonicalRTAVersionLine("Makefile").Get()
	if !hasRTAVersion {
		fmt.Println("No RTA_VERSION declaration found in Makefile")
		return
	}
	asserts.NoError(filepath.WalkDir(".", func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() && (entry.Name() == "vendor" || entry.Name() == ".git" || entry.Name() == "node_modules") {
			return filepath.SkipDir
		}
		if entry.IsDir() || entry.Name() != "Makefile" || path == "Makefile" {
			return nil
		}
		info := asserts.NoError1(entry.Info())
		rawContent := asserts.NoError1(os.ReadFile(path))
		newContent, modified := replaceRTAVersionAssignment(string(rawContent), canonicalLine).Get()
		if !modified {
			return nil
		}
		return os.WriteFile(path, []byte(newContent), info.Mode())
	}))
}

// readCanonicalRTAVersionLine returns the first line in the given Makefile that
// starts with "RTA_VERSION =".
func readCanonicalRTAVersionLine(path string) Option[string] {
	content := asserts.NoError1(os.ReadFile(path))
	for line := range strings.SplitSeq(string(content), "\n") {
		if strings.HasPrefix(line, RTAVersionDeclaration) {
			return Some(line)
		}
	}
	return None[string]()
}

// replaceRTAVersionAssignment provides the given content with the
// RTA_VERSION assignment replaced with the given canonical line.
func replaceRTAVersionAssignment(content string, canonicalLine string) Option[string] {
	lines := strings.Split(content, "\n")
	modified := false
	for lineIndex, line := range lines {
		if !strings.HasPrefix(line, RTAVersionDeclaration) {
			continue
		}
		if line == canonicalLine {
			continue
		}
		lines[lineIndex] = canonicalLine
		modified = true
	}
	if !modified {
		return None[string]()
	}
	return Some(strings.Join(lines, "\n"))
}
