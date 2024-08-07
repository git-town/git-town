package datatable

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/cucumber/godog"

	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/gohacks/stringslice"
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
func (self DataTable) EqualDataTable(other DataTable) (diff string, errorCount int) {
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
func (self *DataTable) EqualGherkin(other *godog.Table) (diff string, errorCount int) {
	if len(self.Cells) == 0 {
		return "your data is empty", 1
	}
	dataTable := FromGherkin(other)
	return self.EqualDataTable(dataTable)
}

// Expand returns a new DataTable instance with the placeholders in this datatable replaced with the given values.
func (self *DataTable) Expand(localRepo runner, remoteRepo runner, worktreeRepo runner, initialDevSHAs map[string]gitdomain.SHA, initialOriginSHAsOpt, initialWorktreeSHAsOpt Option[map[string]gitdomain.SHA]) DataTable {
	var templateRE *regexp.Regexp
	var templateOnce sync.Once
	result := DataTable{}
	for row := range self.Cells {
		cells := []string{}
		for col := range self.Cells[row] {
			cell := self.Cells[row][col]
			if strings.Contains(cell, "{{") {
				templateOnce.Do(func() { templateRE = regexp.MustCompile(`\{\{.*?\}\}`) })
				match := templateRE.FindString(cell)
				switch {
				case strings.HasPrefix(match, "{{ sha "):
					commitName := match[8 : len(match)-4]
					shas := localRepo.SHAsForCommit(commitName)
					if len(shas) == 0 {
						panic(fmt.Sprintf("test workspace has no commit %q", commitName))
					}
					sha := shas.First()
					cell = strings.Replace(cell, match, sha.String(), 1)
				case strings.HasPrefix(match, "{{ sha-in-origin "):
					commitName := match[18 : len(match)-4]
					shas := remoteRepo.SHAsForCommit(commitName)
					sha := shas.First()
					cell = strings.Replace(cell, match, sha.String(), 1)
				case strings.HasPrefix(match, "{{ sha-before-run "):
					commitName := match[19 : len(match)-4]
					sha, found := initialDevSHAs[commitName]
					if !found {
						fmt.Printf("I cannot find the initial dev commit %q.\n", commitName)
						fmt.Printf("I have records about %d commits:\n", len(initialDevSHAs))
						for _, key := range maps.Keys(initialDevSHAs) {
							fmt.Println("  -", key)
						}
						panic("see error above")
					}
					cell = strings.Replace(cell, match, sha.String(), 1)
				case strings.HasPrefix(match, "{{ sha-in-origin-before-run "):
					initialOriginSHAs, has := initialOriginSHAsOpt.Get()
					if !has {
						panic("no origin SHAs recorded")
					}
					commitName := match[29 : len(match)-4]
					sha, found := initialOriginSHAs[commitName]
					if !found {
						fmt.Printf("I cannot find the initial origin commit %q.\n", commitName)
						fmt.Printf("I have records about %d commits:\n", len(initialOriginSHAs))
						for _, key := range maps.Keys(initialOriginSHAs) {
							fmt.Println("  -", key)
						}
					}
					cell = strings.Replace(cell, match, sha.String(), 1)
				case strings.HasPrefix(match, "{{ sha-in-worktree "):
					commitName := match[20 : len(match)-4]
					shas := worktreeRepo.SHAsForCommit(commitName)
					sha := shas.First()
					cell = strings.Replace(cell, match, sha.String(), 1)
				case strings.HasPrefix(match, "{{ sha-in-worktree-before-run "):
					commitName := match[31 : len(match)-4]
					initialWorktreeSHAs, has := initialWorktreeSHAsOpt.Get()
					if !has {
						panic("no initial worktree SHAs recorded")
					}
					sha, found := initialWorktreeSHAs[commitName]
					if !found {
						fmt.Printf("I cannot find the initial worktree commit %q.\n", commitName)
						fmt.Printf("I have records about %d commits:\n", len(initialWorktreeSHAs))
						for _, key := range maps.Keys(initialWorktreeSHAs) {
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
	formatStrings := []string{}
	for _, width := range self.widths() {
		formatStrings = append(formatStrings, fmt.Sprintf("| %%-%dv ", width))
	}
	// render the self using this format
	result := ""
	for row := range self.Cells {
		for col := range self.Cells[row] {
			result += fmt.Sprintf(formatStrings[col], self.Cells[row][col])
		}
		result += "|\n"
	}
	return result
}

// columns provides the self data organized into columns.
func (self *DataTable) columns() [][]string {
	result := [][]string{}
	for column := range self.Cells[0] {
		colData := []string{}
		for row := range self.Cells {
			colData = append(colData, self.Cells[row][column])
		}
		result = append(result, colData)
	}
	return result
}

// widths provides the widths of all columns.
func (self *DataTable) widths() []int {
	result := []int{}
	for _, column := range self.columns() {
		result = append(result, stringslice.Longest(column))
	}
	return result
}
