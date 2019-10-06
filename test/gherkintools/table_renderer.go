package gherkintools

import (
	"fmt"

	"github.com/Originate/git-town/test/helpers"
)

// TableRenderer creates Gherkin-table-like representations of 2d string lists.
type TableRenderer struct {
	// cells contains a list of columns
	// organized as rows and columns
	cells [][]string
}

// AddLine adds the given line to this table.
func (renderer *TableRenderer) AddLine(elements ...string) {
	renderer.cells = append(renderer.cells, elements)
}

// columns provides the data organized in columns.
func (renderer *TableRenderer) columns() (result [][]string) {
	for column := range renderer.cells[0] {
		colData := []string{}
		for row := range renderer.cells {
			colData = append(colData, renderer.cells[row][column])
		}
		result = append(result, colData)
	}
	return result
}

// String provides the textual representation of this table.
func (renderer *TableRenderer) String() (result string) {
	formatStrings := []string{}
	for _, width := range renderer.widths() {
		formatStrings = append(formatStrings, fmt.Sprintf("| %%-%dv ", width))
	}
	for row := range renderer.cells {
		for col := range renderer.cells[row] {
			result += fmt.Sprintf(formatStrings[col], renderer.cells[row][col])
		}
		result += "|\n"
	}
	return result
}

// widths provides the widths of all columns.
func (renderer *TableRenderer) widths() (result []int) {
	for _, column := range renderer.columns() {
		result = append(result, helpers.LongestStringLength(column))
	}
	return result
}
