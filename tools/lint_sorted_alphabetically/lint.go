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
		lintFiles()
	case len(os.Args) == 2 && os.Args[1] == "test":
		runTests()
	default:
		fmt.Printf("Error: unknown argument: %s", os.Args[1])
		os.Exit(1)
	}
}

func displayUsage() {
	fmt.Println(`
This tool verifies that all Go struct definitions and usages list the struct properties sorted alphabetically.

Usage: lint <command>

Available commands:
   run   Lints the source code files
   test  Verifies that this tool works
`[1:])
}

// shouldIgnorePath indicates whether the file with the given path should be ignored (not formatted).
func shouldIgnorePath(path string) bool {
	return false // strings.HasPrefix(path, "vendor/") || path == "src/config/configdomain/push_hook.go" || path == "src/config/configdomain/offline.go" || path == "src/cli/print/logger.go"
}

func lintFiles() {
	err := filepath.WalkDir(".", func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil || dirEntry.IsDir() || !isGoFile(path) || shouldIgnorePath(path) {
			return err
		}
		fmt.Print(".")
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		issues := lintFileContent(string(content))
		for _, issue := range issues {
			fmt.Println("%s:%d  %s", path, issue.line, issue.msg)
		}
		return nil
	})
	fmt.Println()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

type issue struct {
	line int
	msg  string
}

var structDefRE regexp.Regexp = *regexp.MustCompile(`^\s*type (.+) struct \{\n.*?\n\}`)

func lintFileContent(content string) []string {

}

func findStructDefinitions(code string) []string {
	matches
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
