package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	err := filepath.WalkDir(".", func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil || dirEntry.IsDir() || !IsGoTestFile(path) || shouldIgnorePath(path) {
			return err
		}
		fmt.Print(".")
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		newContent := FormatFileContent(string(content))
		return os.WriteFile(path, []byte(newContent), dirEntry.Type().Perm())
	})
	fmt.Println()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

// shouldIgnorePath indicates whether the file with the given path should be ignored (not formatted).
func shouldIgnorePath(path string) bool {
	return path == "main_test.go"
}

func IsTopLevelRunLine(line string) bool {
	return strings.HasPrefix(line, "\tt.Run(\"") && strings.HasSuffix(line, ", func(t *testing.T) {")
}

func isEmptyLine(line string) bool {
	return line == ""
}

func IsGoTestFile(path string) bool {
	return strings.HasSuffix(path, "_test.go")
}

func FormatFileContent(content string) string {
	lines := strings.Split(content, "\n")
	var newLines []string
	previousLineEmpty := false
	for _, line := range lines {
		if IsTopLevelRunLine(line) && !previousLineEmpty {
			newLines = append(newLines, "")
		}
		newLines = append(newLines, line)
		previousLineEmpty = isEmptyLine(line)
	}
	return strings.Join(newLines, "\n")
}
