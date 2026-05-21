package main

import (
	"fmt"
	"os"
	"strings"
)

// readCanonicalRTAVersionLine returns the first line in the given Makefile that
// starts with "RTA_VERSION =".
func readCanonicalRTAVersionLine(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSuffix(line, "\r")
		if strings.HasPrefix(trimmed, "RTA_VERSION =") {
			return trimmed, nil
		}
	}
	return "", fmt.Errorf("no RTA_VERSION assignment found in %s", path)
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
