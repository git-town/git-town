package gherkintools

import (
	"fmt"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test/helpers"
	"github.com/pkg/errors"
)

// Commit describes a Git commit.
type Commit struct {
	Author      string
	Branch      string
	FileContent string
	FileName    string
	Location    []string
	Message     string
	SHA         string
}

// DefaultCommit provides a new Commit instance populated with the default values used in the absence of value specified by the test.
func DefaultCommit() Commit {
	return Commit{
		FileName:    "default_file_name_" + helpers.UniqueString(),
		Message:     "default commit message",
		Location:    []string{"local", "remote"},
		Branch:      "main",
		FileContent: "default file content",
	}
}

// FromGherkinTable provides a Commit collection representing the data in the given Gherkin table.
func FromGherkinTable(table *gherkin.DataTable) (result []Commit, err error) {
	columnNames := []string{}
	for _, cell := range table.Rows[0].Cells {
		columnNames = append(columnNames, cell.Value)
	}
	for _, row := range table.Rows[1:] {
		commit := DefaultCommit()
		for i, cell := range row.Cells {
			err := commit.Set(columnNames[i], cell.Value)
			if err != nil {
				return result, errors.Wrapf(err, "cannot set property %q to %q", columnNames[i], cell.Value)
			}
		}
		result = append(result, commit)
	}
	return result, nil
}

// Set assigns the given value to the property with the given name.
func (commit *Commit) Set(name, value string) (err error) {
	fmt.Printf("setting %q to %q\n", name, value)
	switch name {
	case "BRANCH":
		commit.Branch = value
	case "LOCATION":
		commit.Location = []string{value}
	case "MESSAGE":
		commit.Message = value
	default:
		return fmt.Errorf("unknown Commit property: %s", name)
	}
	return nil
}
