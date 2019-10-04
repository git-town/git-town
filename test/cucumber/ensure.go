// This file defines functions that verify state in Cucumber steps.
// Ensure functions return an error if they fail.

package cucumber

import (
	"fmt"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// EnsureStringSliceMatchesTable compares the given string slice to the given Gherkin table.
// If they don't match, it returns an error
// and might print additional information to the console.
// The comparison ignores whitespace around strings.
func EnsureStringSliceMatchesTable(actual []string, expected *gherkin.DataTable) error {
	if len(expected.Rows) == 0 {
		return fmt.Errorf("Empty table given")
	}
	if len(expected.Rows[0].Cells) != 1 {
		return fmt.Errorf("Table with more than one column given")
	}
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(RenderSlice(actual), RenderTable(expected), false)
	if len(diffs) > 1 {
		fmt.Println(dmp.DiffPrettyText(diffs))
		return fmt.Errorf("Found %d differences", len(diffs))
	}
	return nil
}
