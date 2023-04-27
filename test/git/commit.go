package git

import (
	"fmt"

	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/v8/src/config"
	"github.com/git-town/git-town/v8/test/helpers"
)

// Commit describes a Git commit.
type Commit struct {
	Author      string `exhaustruct:"optional"`
	Branch      string
	FileContent string   `exhaustruct:"optional"`
	FileName    string   `exhaustruct:"optional"`
	Locations   []string `exhaustruct:"optional"`
	Message     string
	SHA         string `exhaustruct:"optional"`
}

// Set assigns the given value to the property with the given Gherkin table name.
func (commit *Commit) Set(name, value string) error {
	switch name {
	case "BRANCH":
		commit.Branch = value
	case "LOCATION":
		commit.Locations = []string{value}
	case "MESSAGE":
		commit.Message = value
	case "FILE NAME":
		commit.FileName = value
	case "FILE CONTENT":
		commit.FileContent = value
	case "AUTHOR":
		commit.Author = value
	default:
		return fmt.Errorf("unknown Commit property: %s", name)
	}
	return nil
}

// DefaultCommit provides a new Commit instance populated with the default values used in the absence of value specified by the test.
func DefaultCommit(filenameSuffix string) Commit {
	return Commit{
		FileName:    "default_file_name_" + filenameSuffix,
		Message:     "default commit message",
		Locations:   []string{"local", config.OriginRemote},
		Branch:      "main",
		FileContent: "default file content",
	}
}

// FromGherkinTable provides a Commit collection representing the data in the given Gherkin table.
func FromGherkinTable(table *messages.PickleStepArgument_PickleTable) ([]Commit, error) {
	columnNames := helpers.TableFields(table)
	lastBranch := ""
	lastLocationName := ""
	result := []Commit{}
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
