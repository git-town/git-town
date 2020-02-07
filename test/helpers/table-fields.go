package helpers

import (
	"github.com/cucumber/godog/gherkin"
)

// TableFields provides the header fields of the given table.
func TableFields(table *gherkin.DataTable) (result []string) {
	for _, cell := range table.Rows[0].Cells {
		result = append(result, cell.Value)
	}
	return result
}
