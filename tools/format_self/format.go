package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	switch {
	case len(os.Args) == 1 || len(os.Args) > 2:
		displayUsage()
	case len(os.Args) == 2 && os.Args[1] == "run":
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
This tool makes all instance variables have the name "self".
See https://github.com/git-town/git-town/issues/2589 for details.

Usage: format <command>

Available commands:
   run   Formats the source code files
   test  Verifies that this tool works
`[1:])
}

// shouldIgnorePath indicates whether the file with the given path should be ignored (not formatted).
func shouldIgnorePath(path string) bool {
	return strings.HasPrefix(path, "vendor/") || path == "src/config/configdomain/push_hook.go" || path == "src/config/configdomain/offline.go"
}

func formatFiles() {
	err := filepath.WalkDir(".", func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil || dirEntry.IsDir() || !isGoFile(path) || shouldIgnorePath(path) {
			return err
		}
		fmt.Print(".")
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		newContent := formatFileContent(string(content))
		return os.WriteFile(path, []byte(newContent), dirEntry.Type().Perm())
	})
	fmt.Println()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

func formatFileContent(content string) string {
	lines := strings.Split(content, "\n")
	for l, line := range lines {
		lines[l] = formatLine(line)
	}
	return strings.Join(lines, "\n")
}

func formatLine(line string) string {
	if !strings.HasPrefix(line, "func (") {
		return line
	}
	instanceRE := regexp.MustCompile(`func \((\w+) (\*?\w+\).*)$`)
	matches := instanceRE.FindStringSubmatch(line)
	if len(matches) < 2 {
		return line
	}
	return strings.Replace(line, "("+matches[1], "(self", 1)
}

func isGoFile(path string) bool {
	if strings.HasSuffix(path, "_test.go") {
		return false
	}
	return strings.HasSuffix(path, ".go")
}

/************************************************************************************
 * TESTS
 */

func runTests() {
	testFormatLine()
	fmt.Println()
}

func testFormatLine() {
	tests := map[string]string{
		"func (bcs *BackendCommands) CommentOutSquashCommitMessage(prefix string) error {": "func (self *BackendCommands) CommentOutSquashCommitMessage(prefix string) error {",
		"func (c *Counter) Count() int {":                                                  "func (self *Counter) Count() int {",
		"	if err != nil {":                                                                 "	if err != nil {",
	}
	for give, want := range tests {
		have := formatLine(give)
		assertEqual(want, have, "testFormatLine")
	}
}

func assertEqual[T comparable](want, have T, testName string) {
	fmt.Print(".")
	if have != want {
		fmt.Printf("\nTEST FAILURE in %q\n", testName)
		fmt.Println("\n\nWANT")
		fmt.Println("--------------------------------------------------------")
		fmt.Println(want)
		fmt.Println("\n\nHAVE")
		fmt.Println("--------------------------------------------------------")
		fmt.Println(have)
		os.Exit(1)
	}
}
