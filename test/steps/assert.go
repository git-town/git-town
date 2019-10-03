package steps

import (
	"fmt"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test/cucumber"
	"github.com/sergi/go-diff/diffmatchpatch"
)

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
	diffs := dmp.DiffMain(cucumber.RenderSlice(actual), cucumber.RenderTable(expected), false)
	if len(diffs) > 1 {
		fmt.Println(dmp.DiffPrettyText(diffs))
		return fmt.Errorf("Found %d differences", len(diffs))
	}
	return nil
}
