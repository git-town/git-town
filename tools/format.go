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
This tool formats Go unit tests to have an empty line before top-level subtests.

Usage: format <command>

Available commands:
   format  Formats the test files
   test    Runs the internal tests for this tool
`[1:])
}

func shouldIgnorePath(path string) bool {
	return path == "main_test.go"
}

func isTopLevelRunLine(line string) bool {
	return strings.HasPrefix(line, "\tt.Run(\"") && strings.HasSuffix(line, ", func(t *testing.T) {")
}

func isEmptyLine(line string) bool {
	return line == ""
}

func isGoTestFile(path string) bool {
	return strings.HasSuffix(path, "_test.go")
}

func formatFileContent(content string) string {
	lines := strings.Split(content, "\n")
	newLines := []string{}
	previousLineEmpty := false
	for _, line := range lines {
		if isTopLevelRunLine(line) && !previousLineEmpty {
			newLines = append(newLines, "")
		}
		newLines = append(newLines, line)
		previousLineEmpty = isEmptyLine(line)
	}
	return strings.Join(newLines, "\n")
}

func formatFiles() {
	err := filepath.WalkDir(".", func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil || shouldIgnorePath(path) || dirEntry.IsDir() || !isGoTestFile(path) {
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

//////////////////////////
// TESTS

func runTests() {
	testIsGoTestFile()
	testIsTopLevelRunLine()
	testFormatContentWithoutSubTests()
	testFormatContentWithSubtests()
	testFormatContentWithNestedSubtests()
	fmt.Println()
}

func testIsTopLevelRunLine() {
	tests := map[string]bool{
		"\tt.Run(\"HasLocalBranch\", func(t *testing.T) {":   true,
		"\t\tt.Run(\"HasLocalBranch\", func(t *testing.T) {": false,
	}
	for give, want := range tests {
		fmt.Print(".")
		have := isTopLevelRunLine(give)
		if have != want {
			fmt.Printf("isTestLine(%s) want %t but have %t\n", give, want, have)
			os.Exit(1)
		}
	}
}

func testFormatContentWithSubtests() {
	give := `
package hosting_test

import (
	"code.gitea.io/sdk/gitea"
)

func TestNewGiteaConnector(t *testing.T) {
	t.Parallel()
	t.Run("top-level test 1", func(t *testing.T) {
		t.Parallel()
		give := 123
	})
	t.Run("top-level test 2", func(t *testing.T) {
		t.Parallel()
		give := 123
	})
}`
	want := `
package hosting_test

import (
	"code.gitea.io/sdk/gitea"
)

func TestNewGiteaConnector(t *testing.T) {
	t.Parallel()

	t.Run("top-level test 1", func(t *testing.T) {
		t.Parallel()
		give := 123
	})

	t.Run("top-level test 2", func(t *testing.T) {
		t.Parallel()
		give := 123
	})
}`
	have := formatFileContent(give)
	assertEqual(want, have, "formatContent with subtests")
}

func testFormatContentWithNestedSubtests() {
	give := `
package hosting_test

import (
	"code.gitea.io/sdk/gitea"
)

func TestNewGiteaConnector(t *testing.T) {
	t.Parallel()
	t.Run("top-level test 1", func(t *testing.T) {
		t.Parallel()
		t.Run("nested test 1a", func(t *testing.T) {
			t.Parallel()
			give := 123
		})
		t.Run("nested test 1b", func(t *testing.T) {
			t.Parallel()
			give := 123
		})
	})
	t.Run("top-level test 2", func(t *testing.T) {
		t.Parallel()
		t.Run("nested test 2a", func(t *testing.T) {
			t.Parallel()
			give := 123
		})
		t.Run("nested test 2b", func(t *testing.T) {
			t.Parallel()
			give := 123
		})
	})
}`
	want := `
package hosting_test

import (
	"code.gitea.io/sdk/gitea"
)

func TestNewGiteaConnector(t *testing.T) {
	t.Parallel()

	t.Run("top-level test 1", func(t *testing.T) {
		t.Parallel()
		t.Run("nested test 1a", func(t *testing.T) {
			t.Parallel()
			give := 123
		})
		t.Run("nested test 1b", func(t *testing.T) {
			t.Parallel()
			give := 123
		})
	})

	t.Run("top-level test 2", func(t *testing.T) {
		t.Parallel()
		t.Run("nested test 2a", func(t *testing.T) {
			t.Parallel()
			give := 123
		})
		t.Run("nested test 2b", func(t *testing.T) {
			t.Parallel()
			give := 123
		})
	})
}`
	have := formatFileContent(give)
	assertEqual(want, have, "formatContent with nested subtests")
}

func testFormatContentWithoutSubTests() {
	give := `
package hosting_test

import (
	"code.gitea.io/sdk/gitea"
)

func TestNewGiteaConnector(t *testing.T) {
	t.Parallel()
	give := "123"
}`
	want := `
package hosting_test

import (
	"code.gitea.io/sdk/gitea"
)

func TestNewGiteaConnector(t *testing.T) {
	t.Parallel()
	give := "123"
}`
	have := formatFileContent(give)
	assertEqual(want, have, "formatContent without subtests")
}

func testIsGoTestFile() {
	tests := map[string]bool{
		"/one/two/three_test.go": true,
		"/one/two/three.go":      false,
		"/one/two_test/three.go": false,
	}
	for give, want := range tests {
		have := isGoTestFile(give)
		assertEqual(want, have, "isGoTestFile")
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
