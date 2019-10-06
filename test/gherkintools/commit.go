package gherkintools

import (
	"fmt"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test/helpers"
)

// Commit describes a Git commit.
type Commit struct {
	Branch      string
	Location    string
	Message     string
	FileName    string
	FileContent string
}

// DefaultCommit provides a new Commit instance populated with the default values used in the absence of value specified by the test.
func DefaultCommit() Commit {
	return Commit{
		FileName:    "default_file_name_" + helpers.UniqueString(),
		Message:     "default commit message",
		Location:    "local and remote",
		Branch:      "main",
		FileContent: "default file content",
	}
}

// FromGherkinTable provides a Commit collection representing the data in the given Gherkin table.
func FromGherkinTable(table *gherkin.DataTable) []Commit {
	result := []Commit{}
	columnNames := []string{}
	for _, cell := range table.Rows[0].Cells {
		columnNames = append(columnNames, cell.Value)
	}
	for _, row := range table.Rows[1:] {
		commit := DefaultCommit()
		for i, cell := range row.Cells {
			commit.Set(columnNames[i], cell.Value)
		}
		result = append(result, commit)
	}
	return result
}

// Set assigns the given value to the property with the given name.
func (commit *Commit) Set(name, value string) (err error) {
	switch value {
	case "BRANCH":
		commit.Branch = value
	case "LOCATION":
		commit.Location = value
	case "MESSAGE":
		commit.Message = value
	default:
		err = fmt.Errorf("unknown CommitData property: %s", name)
	}
	return err
}
