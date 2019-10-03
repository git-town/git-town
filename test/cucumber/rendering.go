package cucumber

import (
	"fmt"

	"github.com/DATA-DOG/godog/gherkin"
)

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
