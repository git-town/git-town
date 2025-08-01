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
	err := filepath.WalkDir(".", func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil || dirEntry.IsDir() || !IsGoFile(path) || shouldIgnorePath(path) {
			return err
		}
		fmt.Print(".")
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		text := string(content)
		newContent := FormatFileContent(text)
		if newContent == text {
			return nil
		}
		return os.WriteFile(path, []byte(newContent), dirEntry.Type().Perm())
	})
	fmt.Println()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

func FormatFileContent(content string) string {
	lines := strings.Split(content, "\n")
	result := make([]string, len(lines))
	for l, line := range lines {
		result[l] = FormatLine(line)
	}
	return strings.Join(result, "\n")
}

func FormatLine(line string) string {
	if !strings.HasPrefix(line, "func (") {
		return line
	}
	matches := instanceRE.FindStringSubmatch(line)
	if len(matches) < 2 {
		return line
	}
	return strings.Replace(line, "("+matches[1], "(self", 1)
}

var instanceRE = regexp.MustCompile(`func \((\w+) (\*?[\w\[\]]+\).*)$`)

// IsGoFile indicates whether the given file path is a Go source code file.
// Tests are not considered source code in the context of this source code formatter.
func IsGoFile(path string) bool {
	if strings.HasSuffix(path, "_test.go") {
		return false
	}
	return strings.HasSuffix(path, ".go")
}

// shouldIgnorePath indicates whether the file with the given path should be ignored (not formatted).
func shouldIgnorePath(path string) bool {
	return strings.HasPrefix(path, "vendor/") ||
		path == "internal/config/configdomain/offline.go" ||
		path == "internal/cli/dialog/switch_branch.go" ||
		path == "internal/gohacks/slice/natural_sort.go" ||
		path == "tools/tests_sorted/matcher/matcher.go" ||
		strings.HasPrefix(path, "tools/stats_release")
}
