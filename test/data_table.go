package test

import (
	"fmt"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test/helpers"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// DataTable allows comparing user-generated data with Gherkin tables.
// The zero value is an empty DataTable.
type DataTable struct {
	// cells contains table data organized as rows and columns
	cells [][]string
}

// FromGherkin provides a DataTable instance populated with data from the given Gherkin table.
func FromGherkin(table *gherkin.DataTable) (result DataTable) {
	for _, tableRow := range table.Rows {
		resultRow := make([]string, len(tableRow.Cells))
		for i, tableCell := range tableRow.Cells {
			resultRow[i] = tableCell.Value
		}
		result.AddRow(resultRow...)
	}
	return result
}

// AddRow adds the given row of table data to this table.
func (table *DataTable) AddRow(elements ...string) {
	table.cells = append(table.cells, elements)
}

// columns provides the table data organized into columns.
func (table *DataTable) columns() (result [][]string) {
	for column := range table.cells[0] {
		colData := []string{}
		for row := range table.cells {
			colData = append(colData, table.cells[row][column])
		}
		result = append(result, colData)
	}
	return result
}

// Equal indicates whether this DataTable instance is equal to the given Gherkin table.
// If both are equal it returns an empty string,
// otherwise a diff printable on the console.
func (table *DataTable) Equal(other *gherkin.DataTable) (diff string, errorCount int) {
	if len(table.cells) == 0 {
		return "your data is empty", 1
	}
	gherkinTable := FromGherkin(other)
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(gherkinTable.String(), table.String(), false)
	if len(diffs) == 1 && diffs[0].Type == 0 {
		return "", 0
	}
	return dmp.DiffPrettyText(diffs), len(diffs)
}

// String provides the data in this DataTable instance formatted in Gherkin table format.
func (table *DataTable) String() (result string) {
	// determine how to format each column
	formatStrings := []string{}
	for _, width := range table.widths() {
		formatStrings = append(formatStrings, fmt.Sprintf("| %%-%dv ", width))
	}

	// render the table using this format
	for row := range table.cells {
		for col := range table.cells[row] {
			result += fmt.Sprintf(formatStrings[col], table.cells[row][col])
		}
		result += "|\n"
	}
	return result
}

// widths provides the widths of all columns.
func (table *DataTable) widths() (result []int) {
	for _, column := range table.columns() {
		result = append(result, helpers.LongestStringLength(column))
	}
	return result
}
