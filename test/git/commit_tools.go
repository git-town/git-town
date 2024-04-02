package git

import (
	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/test/helpers"
)

// DefaultCommit provides a new Commit instance populated with the default values used in the absence of value specified by the test.
func DefaultCommit(filenameSuffix string) Commit {
	return Commit{
		Branch:      gitdomain.NewLocalBranchName("main"),
		FileContent: "default file content",
		FileName:    "default_file_name_" + filenameSuffix,
		Locations:   Locations{LocationLocal, LocationOrigin},
		Message:     "default commit message",
	}
}

// FromGherkinTable provides a Commit collection representing the data in the given Gherkin table.
func FromGherkinTable(table *messages.PickleStepArgument_PickleTable) []Commit {
	columnNames := helpers.TableFields(table)
	lastBranch := ""
	lastLocationName := ""
	result := []Commit{}
	counter := helpers.AtomicCounter{}
	for _, row := range table.Rows[1:] {
		commit := DefaultCommit(counter.ToString())
		for cellNo, cell := range row.Cells {
			columnName := columnNames[cellNo]
			cellValue := cell.Value
			if columnName == "BRANCH" {
				if cell.Value == "" {
					cellValue = lastBranch
				} else {
					lastBranch = cellValue
				}
			}
			if columnName == "LOCATION" {
				if cell.Value == "" {
					cellValue = lastLocationName
				} else {
					lastLocationName = cellValue
				}
			}
			commit.Set(columnName, cellValue)
		}
		result = append(result, commit)
	}
	return result
}
