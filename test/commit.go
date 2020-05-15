package test

import (
	"fmt"

	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/test/helpers"
)

// Commit describes a Git commit.
type Commit struct {
	Author      string
	Branch      string
	FileContent string
	FileName    string
	Locations   []string
	Message     string
	SHA         string
}

// DefaultCommit provides a new Commit instance populated with the default values used in the absence of value specified by the test.
func DefaultCommit() Commit {
	return Commit{
		FileName:    "default_file_name_" + helpers.UniqueString(),
		Message:     "default commit message",
		Locations:   []string{"local", "remote"},
		Branch:      "main",
		FileContent: "default file content",
	}
}

// FromGherkinTable provides a Commit collection representing the data in the given Gherkin table.
func FromGherkinTable(table *messages.PickleStepArgument_PickleTable) (result []Commit, err error) {
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
			err := commit.set(columnName, cellValue)
			if err != nil {
				return result, fmt.Errorf("cannot set property %q to %q: %w", columnNames[i], cell.Value, err)
			}
		}
		result = append(result, commit)
	}
	return result, nil
}

// Set assigns the given value to the property with the given Gherkin table name.
func (commit *Commit) set(name, value string) (err error) {
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
