package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	switch {
	case len(os.Args) == 1 || len(os.Args) > 2:
		displayUsage()
	case len(os.Args) == 2 && os.Args[1] == "format":
		formatFiles()
	case len(os.Args) == 2 && os.Args[1] == "test":
		runTests()
	default:
		fmt.Printf("Error: unknown argument: %s", os.Args[1])
		os.Exit(1)
	}
}

func displayUsage() {
	fmt.Println(`
Usage: format <command>

Available commands:
   format  Formats the test files
   test    Runs the internal tests for this tool
`[1:])
}

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

func formatContent(content string) string {
	lines := strings.Split(content, "\n")
	newContent := []string{}
	foundTestLine := false
	foundParallelLine := false
	for _, line := range lines {
		if isTestLine(line) {
			foundTestLine = true
			newContent = append(newContent, line)
			continue
		}
		if foundTestLine {
			if !isParallelLine(line) {
				// tests without a "t.Parallel()" line will not be formatted
				return content
			}
			foundTestLine = false
			foundParallelLine = true
			newContent = append(newContent, line)
			continue
		}
		if foundParallelLine {
			foundParallelLine = false
			if isEmptyLine(line) {
				newContent = append(newContent, line)
				continue
			}
			newContent = append(newContent, "")
			newContent = append(newContent, line)
			continue
		}
		newContent = append(newContent, line)
	}
	return strings.Join(newContent, "\n")
}

func formatFiles() {
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
		fmt.Print(".")
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		newContent := formatContent(string(content))
		perm := dirEntry.Type().Perm()
		return os.WriteFile(path, []byte(newContent), perm)
	})
	fmt.Println()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

//////////////////////////
// TESTS

func runTests() {
	testIsTestLine()
	testFormatContent()
	fmt.Println()
}

func testIsTestLine() {
	tests := map[string]bool{
		"func TestFoo(t *testing.T) {": true,
		"func TestFoo":                 false,
		"func Other(t *testing.T) {":   false,
	}
	for give, want := range tests {
		fmt.Print(".")
		have := isTestLine(give)
		if have != want {
			panic(fmt.Sprintf("isTestLine(%s) want %t but have %t", give, want, have))
		}
	}
}

func testFormatContent() {
}
