package cucumber

import (
	"fmt"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test/helpers"
)

// RenderSlice provides the textual Gherkin table representation of the given slice.
func RenderSlice(data []string) (result string) {
	width := helpers.LongestStringLength(data)
	formatStr := fmt.Sprintf("| %%-%dv |\n", width)
	for i := range data {
		result += fmt.Sprintf(formatStr, data[i])
	}
	return result
}

// RenderTable provides the textual Gherkin representation of the given Gherkin table.
func RenderTable(table *gherkin.DataTable) string {
	slice := []string{}
	for i := 1; i < len(table.Rows); i++ {
		slice = append(slice, table.Rows[i].Cells[0].Value)
	}
	return RenderSlice(slice)
}
