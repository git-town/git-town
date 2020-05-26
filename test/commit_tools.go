package test

import (
	"fmt"

	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/test/helpers"
)

// DefaultCommit provides a new Commit instance populated with the default values used in the absence of value specified by the test.
func DefaultCommit() git.Commit {
	return git.Commit{
		FileName:    "default_file_name_" + helpers.UniqueString(),
		Message:     "default commit message",
		Locations:   []string{"local", "remote"},
		Branch:      "main",
		FileContent: "default file content",
	}
}

// FromGherkinTable provides a Commit collection representing the data in the given Gherkin table.
func FromGherkinTable(table *messages.PickleStepArgument_PickleTable) (result []git.Commit, err error) {
	columnNames := helpers.TableFields(table)
	lastBranchName := ""
	lastLocationName := ""
	for _, row := range table.Rows[1:] {
		commit := DefaultCommit()
		for i, cell := range row.Cells {
			columnName := columnNames[i]
			cellValue := cell.Value
			if columnName == "BRANCH" {
				if cell.Value == "" {
					cellValue = lastBranchName
				} else {
					lastBranchName = cellValue
				}
			}
			if columnName == "LOCATION" {
				if cell.Value == "" {
					cellValue = lastLocationName
				} else {
					lastLocationName = cellValue
				}
			}
			err := commit.Set(columnName, cellValue)
			if err != nil {
				return result, fmt.Errorf("cannot set property %q to %q: %w", columnNames[i], cell.Value, err)
			}
		}
		result = append(result, commit)
	}
	return result, nil
}
