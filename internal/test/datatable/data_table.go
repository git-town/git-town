package datatable

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cucumber/godog"
	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/test/handlebars"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// DataTable allows comparing user-generated data with Gherkin tables.
// The zero value is an empty DataTable.
type DataTable struct {
	// table data organized as rows and columns
	Cells [][]string `exhaustruct:"optional"`
}

// FromGherkin provides a DataTable instance populated with data from the given Gherkin table.
func FromGherkin(table *godog.Table) DataTable {
	result := DataTable{}
	for _, tableRow := range table.Rows {
		resultRow := make([]string, len(tableRow.Cells))
		for t, tableCell := range tableRow.Cells {
			resultRow[t] = tableCell.Value
		}
		result.AddRow(resultRow...)
	}
	return result
}

// AddRow adds the given row of table data to this table.
func (self *DataTable) AddRow(elements ...string) {
	self.Cells = append(self.Cells, elements)
}

// EqualDataTable compares this DataTable instance to the given DataTable.
// If both are equal it returns an empty string, otherwise a diff printable on the console.
func (self *DataTable) EqualDataTable(other DataTable) (diff string, errorCount int) { //nolint:nonamedreturns
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(other.String(), self.String(), false)
	if len(diffs) == 1 && diffs[0].Type == 0 {
		return "", 0
	}
	result := dmp.DiffPrettyText(diffs)
	result += "\n\nReceived this table:\n\n"
	result += self.String()
	return result, len(diffs)
}

// EqualGherkin compares this DataTable instance to the given Gherkin self.
// If both are equal it returns an empty string, otherwise a diff printable on the console.
func (self *DataTable) EqualGherkin(other *godog.Table) (diff string, errorCount int) { //nolint:nonamedreturns
	if len(self.Cells) == 0 {
		return "your data is empty", 1
	}
	dataTable := FromGherkin(other)
	return self.EqualDataTable(dataTable)
}

// Expand returns a new DataTable instance with the placeholders in this datatable replaced with the given values.
func (self *DataTable) Expand(args handlebars.ExpandArgs) DataTable {
	result := DataTable{}
	for row := range self.Cells {
		var cells []string
		for col := range self.Cells[row] {
			cell := self.Cells[row][col]
			cell = handlebars.Expand(cell, args)
			cells = append(cells, cell)
		}
		result.AddRow(cells...)
	}
	return result
}

// RemoveText deletes the given text from each cell.
func (self *DataTable) RemoveText(text string) {
	for row := range self.Cells {
		for col := range self.Cells[row] {
			self.Cells[row][col] = strings.Replace(self.Cells[row][col], text, "", 1)
		}
	}
}

// Sorted provides a new DataTable that contains the content of this DataTable sorted by the first column.
func (self *DataTable) Sort() {
	sort.Slice(self.Cells, func(a, b int) bool {
		return self.Cells[a][0] < self.Cells[b][0]
	})
}

// String provides the data in this DataTable instance formatted in Gherkin self format.
func (self *DataTable) String() string {
	if len(self.Cells) == 0 {
		return ""
	}
	// determine how to format each column
	widths := self.widths()
	formatStrings := make([]string, len(widths))
	for w, width := range widths {
		formatStrings[w] = fmt.Sprintf("| %%-%dv ", width)
	}
	// render the self using this format
	result := strings.Builder{}
	for row := range self.Cells {
		for col := range self.Cells[row] {
			result.WriteString(fmt.Sprintf(formatStrings[col], gohacks.EscapeNewLines(self.Cells[row][col])))
		}
		result.WriteString("|\n")
	}
	return result.String()
}

// columns provides the self data organized into columns.
func (self *DataTable) columns() [][]string {
	columns := self.Cells[0]
	result := make([][]string, len(columns))
	for c := range columns {
		var colData []string
		for row := range self.Cells {
			colData = append(colData, self.Cells[row][c])
		}
		result[c] = colData
	}
	return result
}

// widths provides the widths of all columns.
func (self *DataTable) widths() []int {
	columns := self.columns()
	result := make([]int, len(columns))
	for c, column := range columns {
		result[c] = stringslice.Longest(column)
	}
	return result
}
