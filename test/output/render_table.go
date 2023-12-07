package output

import (
	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/v11/test/datatable"
)

// RenderTable provides the textual Gherkin representation of the given Gherkin table.
func RenderTable(table *messages.PickleStepArgument_PickleTable) string {
	result := datatable.DataTable{}
	for _, row := range table.Rows {
		values := []string{}
		for _, cell := range row.Cells {
			values = append(values, cell.Value)
		}
		result.AddRow(values...)
	}
	return result.String()
}
