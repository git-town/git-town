package test

import (
	"fmt"

	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/test/helpers"
)

// DefaultCommit provides a new Commit instance populated with the default values used in the absence of value specified by the test.
func DefaultCommit(filenameSuffix string) git.Commit {
	return git.Commit{
		FileName:    "default_file_name_" + filenameSuffix,
		Message:     "default commit message",
		Locations:   []string{"local", config.OriginRemote},
		Branch:      "main",
		FileContent: "default file content",
	}
}

// FromGherkinTable provides a Commit collection representing the data in the given Gherkin table.
func FromGherkinTable(table *messages.PickleStepArgument_PickleTable) ([]git.Commit, error) {
	columnNames := helpers.TableFields(table)
	lastBranch := ""
	lastLocationName := ""
	result := []git.Commit{}
	counter := helpers.Counter{}
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
			err := commit.Set(columnName, cellValue)
			if err != nil {
				return result, fmt.Errorf("cannot set property %q to %q: %w", columnNames[cellNo], cell.Value, err)
			}
		}
		result = append(result, commit)
	}
	return result, nil
}
