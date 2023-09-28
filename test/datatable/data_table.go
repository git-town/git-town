package datatable

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/test/helpers"
	"github.com/sergi/go-diff/diffmatchpatch"
	"golang.org/x/exp/maps"
)

// DataTable allows comparing user-generated data with Gherkin tables.
// The zero value is an empty DataTable.
type DataTable struct {
	// table data organized as rows and columns
	Cells [][]string `exhaustruct:"optional"`
}

// FromGherkin provides a DataTable instance populated with data from the given Gherkin table.
func FromGherkin(table *messages.PickleStepArgument_PickleTable) DataTable {
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
func (table *DataTable) AddRow(elements ...string) {
	table.Cells = append(table.Cells, elements)
}

// columns provides the table data organized into columns.
func (table *DataTable) columns() [][]string {
	result := [][]string{}
	for column := range table.Cells[0] {
		colData := []string{}
		for row := range table.Cells {
			colData = append(colData, table.Cells[row][column])
		}
		result = append(result, colData)
	}
	return result
}

// EqualDataTable compares this DataTable instance to the given DataTable.
// If both are equal it returns an empty string, otherwise a diff printable on the console.
func (table *DataTable) EqualDataTable(other DataTable) (diff string, errorCount int) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(other.String(), table.String(), false)
	if len(diffs) == 1 && diffs[0].Type == 0 {
		return "", 0
	}
	return dmp.DiffPrettyText(diffs), len(diffs)
}

// EqualGherkin compares this DataTable instance to the given Gherkin table.
// If both are equal it returns an empty string, otherwise a diff printable on the console.
func (table *DataTable) EqualGherkin(other *messages.PickleStepArgument_PickleTable) (diff string, errorCount int) {
	if len(table.Cells) == 0 {
		return "your data is empty", 1
	}
	dataTable := FromGherkin(other)
	return table.EqualDataTable(dataTable)
}

// Expand returns a new DataTable instance with the placeholders in this datatable replaced with the given values.
func (table *DataTable) Expand(localRepo runner, remoteRepo runner, initialSHAs map[string]domain.SHA) DataTable {
	var templateRE *regexp.Regexp
	var templateOnce sync.Once
	result := DataTable{}
	for row := range table.Cells {
		cells := []string{}
		for col := range table.Cells[row] {
			cell := table.Cells[row][col]
			if strings.Contains(cell, "{{") {
				templateOnce.Do(func() { templateRE = regexp.MustCompile(`\{\{.*?\}\}`) })
				match := templateRE.FindString(cell)
				switch {
				case strings.HasPrefix(match, "{{ sha "):
					commitName := match[8 : len(match)-4]
					sha := localRepo.SHAForCommit(commitName)
					cell = strings.Replace(cell, match, sha, 1)
				case strings.HasPrefix(match, "{{ sha-in-origin "):
					commitName := match[18 : len(match)-4]
					sha := remoteRepo.SHAForCommit(commitName)
					cell = strings.Replace(cell, match, sha, 1)
				case strings.HasPrefix(match, "{{ sha-before-run "):
					commitName := match[19 : len(match)-4]
					sha, found := initialSHAs[commitName]
					if !found {
						fmt.Printf("I cannot find the initial commit %q.\n", commitName)
						fmt.Println("I have these commits:")
						for _, key := range maps.Keys(initialSHAs) {
							fmt.Println("  -", key)
						}
					}
					cell = strings.Replace(cell, match, sha.String(), 1)
				default:
					log.Fatalf("DataTable.Expand: unknown template expression %q", cell)
				}
			}
			cells = append(cells, cell)
		}
		result.AddRow(cells...)
	}
	return result
}

// RemoveText deletes the given text from each cell.
func (table *DataTable) RemoveText(text string) {
	for row := range table.Cells {
		for col := range table.Cells[row] {
			table.Cells[row][col] = strings.Replace(table.Cells[row][col], text, "", 1)
		}
	}
}

// Sorted provides a new DataTable that contains the content of this DataTable sorted by the first column.
func (table *DataTable) Sort() {
	sort.Slice(table.Cells, func(a, b int) bool {
		return table.Cells[a][0] < table.Cells[b][0]
	})
}

// String provides the data in this DataTable instance formatted in Gherkin table format.
func (table *DataTable) String() string {
	// determine how to format each column
	formatStrings := []string{}
	for _, width := range table.widths() {
		formatStrings = append(formatStrings, fmt.Sprintf("| %%-%dv ", width))
	}
	// render the table using this format
	result := ""
	for row := range table.Cells {
		for col := range table.Cells[row] {
			result += fmt.Sprintf(formatStrings[col], table.Cells[row][col])
		}
		result += "|\n"
	}
	return result
}

// widths provides the widths of all columns.
func (table *DataTable) widths() []int {
	result := []int{}
	for _, column := range table.columns() {
		result = append(result, helpers.LongestStringLength(column))
	}
	return result
}
