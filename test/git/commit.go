package git

import (
	"log"

	"github.com/cucumber/godog"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/test/helpers"
)

// Commit describes a Git commit.
type Commit struct {
	Author      gitdomain.Author `exhaustruct:"optional"`
	Branch      gitdomain.LocalBranchName
	FileContent string    `exhaustruct:"optional"`
	FileName    string    `exhaustruct:"optional"`
	Locations   Locations `exhaustruct:"optional"`
	Message     gitdomain.CommitMessage
	SHA         gitdomain.SHA `exhaustruct:"optional"`
}

var counter helpers.AtomicCounter //nolint:gochecknoglobals

// Set assigns the given value to the property with the given Gherkin table name.
func (self *Commit) Set(name, value string) {
	switch name {
	case "BRANCH":
		self.Branch = gitdomain.NewLocalBranchName(value)
	case "LOCATION":
		self.Locations = NewLocations(value)
	case "MESSAGE":
		self.Message = gitdomain.CommitMessage(value)
	case "FILE NAME":
		self.FileName = value
	case "FILE CONTENT":
		self.FileContent = value
	case "AUTHOR":
		self.Author = gitdomain.Author(value)
	default:
		log.Fatalf("unknown Commit property: %s", name)
	}
}

// DefaultCommit provides a new Commit instance populated with the default values used in the absence of value specified by the test.
func DefaultCommit() Commit {
	return Commit{
		Branch:      gitdomain.NewLocalBranchName("main"),
		FileContent: "default file content",
		FileName:    "default_file_name_" + counter.NextAsString(),
		Locations:   Locations{LocationLocal, LocationOrigin},
		Message:     "default commit message",
	}
}

// FromGherkinTable provides a Commit collection representing the data in the given Gherkin table.
func FromGherkinTable(table *godog.Table) []Commit {
	columnNames := helpers.TableFields(table)
	lastBranch := ""
	lastLocationName := ""
	result := []Commit{}
	for _, row := range table.Rows[1:] {
		commit := DefaultCommit()
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
