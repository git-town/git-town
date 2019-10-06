package gherkintools

import (
	"log"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test/helpers"
)

// CommitTableEntry contains the elements of a Gherkin table defining commit data.
type CommitTableEntry struct {
	Branch      string
	Location    string
	Message     string
	FileName    string
	FileContent string
}

// NewCommitTableEntry provides a new CommitTableEntry with default values
func NewCommitTableEntry() CommitTableEntry {
	return CommitTableEntry{
		FileName:    "default_file_name_" + helpers.UniqueString(),
		Message:     "default commit message",
		Location:    "local and remote",
		Branch:      "main",
		FileContent: "default file content",
	}
}

// ParseGherkinTable provides a CommitTableEntry slice describing the given Gherkin table.
func ParseGherkinTable(table *gherkin.DataTable) []CommitTableEntry {
	result := []CommitTableEntry{}
	columnNames := []string{}
	for _, cell := range table.Rows[0].Cells {
		columnNames = append(columnNames, cell.Value)
	}
	for _, row := range table.Rows[1:] {
		commit := NewCommitTableEntry()
		for i, cell := range row.Cells {
			switch columnNames[i] {
			case "BRANCH":
				commit.Branch = cell.Value
			case "LOCATION":
				commit.Location = cell.Value
			case "MESSAGE":
				commit.Message = cell.Value
			default:
				log.Fatalf("GitRepository.parseCommitsTable: unknown column name: %s", columnNames[i])
			}
		}
		result = append(result, commit)
	}
	return result
}
