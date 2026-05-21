package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/git-town/git-town/v23/pkg/asserts"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

func main() {
	canonicalLine, hasRTAVersion := readCanonicalRTAVersionLine("Makefile").Get()
	if !hasRTAVersion {
		fmt.Println("No RTA_VERSION declaration found in Makefile")
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

// readCanonicalRTAVersionLine returns the first line in the given Makefile that
// starts with "RTA_VERSION =".
func readCanonicalRTAVersionLine(path string) Option[string] {
	content := asserts.NoError1(os.ReadFile(path))
	for line := range strings.SplitSeq(string(content), "\n") {
		if strings.HasPrefix(line, "RTA_VERSION =") {
			return Some(line)
		}
	}
	return None[string]()
}

// replaceRTAVersionAssignment replaces the first line in content that starts with
// optional whitespace followed by "RTA_VERSION =" with canonicalLine.
// It preserves the original line endings (LF or CRLF).
// Returns the updated content and whether any replacement was made.
func replaceRTAVersionAssignment(content string, canonicalLine string) (string, bool) {
	lines := strings.Split(content, "\n")
	modified := false
	for lineIndex, line := range lines {
		hasCR := strings.HasSuffix(line, "\r")
		bare := strings.TrimSuffix(line, "\r")
		bare = strings.TrimLeft(bare, " \t")
		if !strings.HasPrefix(bare, "RTA_VERSION =") {
			continue
		}
		if bare == canonicalLine {
			continue
		}
		if hasCR {
			lines[lineIndex] = canonicalLine + "\r"
		} else {
			lines[lineIndex] = canonicalLine
		}
		modified = true
	}
	return strings.Join(lines, "\n"), modified
}
