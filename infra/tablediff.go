package infra

import (
	"fmt"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// EqualsStringSlice compares the given string slice to the given Gherkin table.
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

// RenderSlice returns the Gherkin table representation of the given slice
func RenderSlice(data []string) string {

	// determine the width of the longest string
	width := 0
	for _, text := range data {
		if len(text) > width {
			width = len(text)
		}
	}

	// render
	result := ""
	formatStr := fmt.Sprintf("| %%-%dv |\n", width)
	for _, text := range data {
		result += fmt.Sprintf(formatStr, text)
	}
	return result
}

// RenderTable returns the Gherkin representation of the given Gherkin table
func RenderTable(table *gherkin.DataTable) string {

	// determine the width of the table
	width := 0
	for _, row := range table.Rows {
		cellWidth := len(row.Cells[0].Value)
		if (cellWidth) > width {
			width = cellWidth
		}
	}

	// convert table to slice
	slice := []string{}
	for i := 1; i < len(table.Rows); i++ {
		slice = append(slice, table.Rows[i].Cells[0].Value)
	}

	return RenderSlice(slice)
}
