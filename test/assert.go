package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
)

func assertFileExists(t *testing.T, dir, filename string) {
	filePath := path.Join(dir, filename)
	info, err := os.Stat(filePath)
	assert.Nilf(t, err, "file %q does not exist", filePath)
	assert.Falsef(t, info.IsDir(), "%q is a directory", filePath)
}

func assertFileExistsWithContent(t *testing.T, dir, filename, expectedContent string) {
	fileContent, err := ioutil.ReadFile(path.Join(dir, filename))
	assert.Nil(t, err)
	assert.Equal(t, expectedContent, string(fileContent))
}

func assertFolderExists(t *testing.T, dir string) {
	_, err := os.Stat(dir)
	assert.Falsef(t, os.IsNotExist(err), "directory %q not found", dir)
}

func assertHasGitBranch(t *testing.T, dir, expectedBranch string) {
	runner := NewShellRunner(dir)
	output, err := runner.Run("git", "branch")
	assert.Nilf(t, err, "cannot run 'git status' in %q", dir)
	assert.Contains(t, output, expectedBranch, "doesn't have Git branch")
}

func assertIsBareGitRepo(t *testing.T, dir string) {
	assertFolderExists(t, dir)
	entries, err := ioutil.ReadDir(dir)
	assert.Nilf(t, err, "cannot list directory %q", dir)
	testData := []string{"HEAD", "branches", "config", "description", "hooks", "info", "objects", "refs"}
	for i, expected := range testData {
		assert.Equal(t, expected, entries[i].Name())
	}
}

func assertIsNormalGitRepo(t *testing.T, dir string) {
	assertFolderExists(t, dir)
	entries, err := ioutil.ReadDir(dir)
	assert.Nilf(t, err, "cannot list directory %q", dir)
	testData := []string{".git"}
	for i, expected := range testData {
		assert.Equal(t, expected, entries[i].Name())
	}
}

// AssertStringSliceMatchesTable compares the given string slice to the given Gherkin table.
// If they don't match, it returns an error
// and might print additional information to the console.
// The comparison ignores whitespace around strings.
func AssertStringSliceMatchesTable(actual []string, expected *gherkin.DataTable) error {

	// ensure we have a valid table
	if len(expected.Rows) == 0 {
		return fmt.Errorf("Empty table given")
	}
	if len(expected.Rows[0].Cells) != 1 {
		return fmt.Errorf("Table with more than one column given")
	}

	// render the slice
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(RenderSlice(actual), RenderTable(expected), false)
	if len(diffs) > 1 {
		fmt.Println(dmp.DiffPrettyText(diffs))
		return fmt.Errorf("Found %d differences", len(diffs))
	}
	return nil
}
