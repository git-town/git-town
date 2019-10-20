package test

import (
	"fmt"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/Originate/git-town/test/helpers"
)

// DataTable compares Gherkin tables with user-generated data.
type DataTable struct {
	// cells contains table data organized as rows and columns
	cells [][]string
}

// FromGherkin provides a DataTable instance populated with data from the given Gherkin table.
func FromGherkin(table *gherkin.DataTable) (result DataTable) {
	for _, tableRow := range table.Rows {
		mortaRow := make([]string, len(tableRow.Cells))
		for i, tableCell := range tableRow.Cells {
			mortaRow[i] = tableCell.Value
		}
		result.AddRow(mortaRow...)
	}
	return result
}

// AddRow adds the given row of table data to this table.
func (morta *DataTable) AddRow(elements ...string) {
	morta.cells = append(morta.cells, elements)
}

// columns provides the table data organized into columns.
func (morta *DataTable) columns() (result [][]string) {
	for column := range morta.cells[0] {
		colData := []string{}
		for row := range morta.cells {
			colData = append(colData, morta.cells[row][column])
		}
		result = append(result, colData)
	}
	return result
}

// Equal indicates whether this DataTable instance is equal to the given Gherkin table.
// If both are equal it returns an empty string,
// otherwise a diff printable on the console.
func (morta *DataTable) Equal(table *gherkin.DataTable) (diff string, errorCount int) {
	if len(morta.cells) == 0 {
		return "your data is empty", 1
	}
	gherkinTable := FromGherkin(table)
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(gherkinTable.String(), morta.String(), false)
	if len(diffs) == 1 && diffs[0].Type == 0 {
		return "", 0
	}
	return dmp.DiffPrettyText(diffs), len(diffs)
}

// String provides the data in this DataTable instance formatted in Gherkin table format.
func (morta *DataTable) String() (result string) {
	// determine how to format each column
	formatStrings := []string{}
	for _, width := range morta.widths() {
		formatStrings = append(formatStrings, fmt.Sprintf("| %%-%dv ", width))
	}

	// render the table using this format
	for row := range morta.cells {
		for col := range morta.cells[row] {
			result += fmt.Sprintf(formatStrings[col], morta.cells[row][col])
		}
		result += "|\n"
	}
	return result
}

// widths provides the widths of all columns.
func (morta *DataTable) widths() (result []int) {
	for _, column := range morta.columns() {
		result = append(result, helpers.LongestStringLength(column))
	}
	return result
}
