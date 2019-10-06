// This file defines functions that verify state in Cucumber steps.
// Ensure functions return an error if they fail.

package test

import (
	"fmt"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// EnsureExecutedGitCommandsMatchTable compares the given list of executed Git commands contains the same data as the given Gherkin table.
// If they don't match, it returns an error
// and might print additional information to the console.
// The comparison ignores whitespace around strings.
func EnsureExecutedGitCommandsMatchTable(actual []ExecutedGitCommand, expected *gherkin.DataTable) error {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(RenderExecutedGitCommands(actual), RenderTable(expected), false)
	if len(diffs) == 1 && diffs[0].Type == 0 {
		return nil
	}
	fmt.Println(dmp.DiffPrettyText(diffs))
	return fmt.Errorf("found %d differences", len(diffs))
}
