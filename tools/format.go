package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func isTestLine(line string) bool {
	return strings.HasPrefix(line, "func Test") && strings.HasSuffix(line, "(t *testing.T) {")
}

func formatFile(path string, perm os.FileMode) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	newContent := []string{}
	for _, line := range strings.Split(string(content), "\n") {
		if isTestLine(line) {
			newContent = append(newContent, line)
			newContent = append(newContent, "")
		} else {
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
