package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func shouldIgnore(path string) bool {
	return path == "main_test.go"
}

func isTestLine(line string) bool {
	return strings.HasPrefix(line, "func Test") && strings.HasSuffix(line, "(t *testing.T) {")
}

func isParallelLine(line string) bool {
	return line == "\tt.Parallel()"
}

func isEmptyLine(line string) bool {
	return line == ""
}

func formatFile(path string, perm os.FileMode) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	newContent := []string{}
	foundTestLine := false
	foundParallelLine := false
	for l, line := range strings.Split(string(content), "\n") {
		if isTestLine(line) {
			foundTestLine = true
			newContent = append(newContent, line)
			continue
		}
		if foundTestLine {
			if !isParallelLine(line) {
				// tests without a "t.Parallel()" line will not be formatted
				return nil
			}
			foundTestLine = false
			foundParallelLine = true
			newContent = append(newContent, line)
			continue
		}
		if foundParallelLine {
			if isEmptyLine(line) {
				newContent = append(newContent, line)
				continue
			}
			newContent = append(newContent, "")
			newContent = append(newContent, line)
		}
	}
	err = os.WriteFile(path, []byte(strings.Join(newContent, "\n")), perm)
	if err != nil {
		return err
	}
	fmt.Print(".")
	return nil
}

func main() {
	err := filepath.WalkDir(".", func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if shouldIgnore(path) {
			return nil
		}
		if dirEntry.IsDir() {
			return nil
		}
		if !strings.HasSuffix(dirEntry.Name(), "_test.go") {
			return nil
		}
		return formatFile(path, dirEntry.Type().Perm())
	})
	fmt.Println()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}
