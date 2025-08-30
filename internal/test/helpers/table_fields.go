package helpers

import "github.com/cucumber/godog"

// TableFields provides the header fields of the given table.
func TableFields(table *godog.Table) []string {
	cells := table.Rows[0].Cells
	result := make([]string, len(cells))
	for c, cell := range cells {
		result[c] = cell.Value
	}
	return result
}
