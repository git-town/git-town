// This file defines functions that verify state in Cucumber steps.
// Ensure functions return an error if they fail.

package test

import (
	"fmt"

	"github.com/DATA-DOG/godog/gherkin"
)

// EnsureExecutedGitCommandsMatchTable compares the given list of executed Git commands contains the same data as the given Gherkin table.
// If they don't match, it returns an error
// and might print additional information to the console.
// The comparison ignores whitespace around strings.
func EnsureExecutedGitCommandsMatchTable(actual []ExecutedGitCommand, expected *gherkin.DataTable) error {
	morta := RenderExecutedGitCommands(actual, expected)
	diff, errorCount := morta.Equal(expected)
	if errorCount != 0 {
		fmt.Println(diff)
		return fmt.Errorf("found %d differences", errorCount)
	}
	return nil
}
