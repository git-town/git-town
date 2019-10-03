package test

import "github.com/dchest/uniuri"

// CommitTableEntry contains the elements of a Gherkin table defining commit data.
type CommitTableEntry struct {
	branch      string
	location    string
	message     string
	fileName    string
	fileContent string
}

// NewCommitTableEntry provides a new CommitTableEntry with default values
func NewCommitTableEntry() CommitTableEntry {
	return CommitTableEntry{
		fileName:    "default_file_name_" + uniuri.NewLen(10),
		message:     "default commit message",
		location:    "local and remote",
		branch:      "main",
		fileContent: "default file content",
	}
}
