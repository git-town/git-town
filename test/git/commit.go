package git

import (
	"log"

	"github.com/git-town/git-town/v9/src/domain"
)

// Commit describes a Git commit.
type Commit struct {
	Author      string `exhaustruct:"optional"`
	Branch      domain.LocalBranchName
	FileContent string   `exhaustruct:"optional"`
	FileName    string   `exhaustruct:"optional"`
	Locations   []string `exhaustruct:"optional"`
	Message     string
	SHA         domain.SHA `exhaustruct:"optional"`
}

// Set assigns the given value to the property with the given Gherkin table name.
func (commit *Commit) Set(name, value string) {
	switch name {
	case "BRANCH":
		commit.Branch = domain.NewLocalBranchName(value)
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
		log.Fatalf("unknown Commit property: %s", name)
	}
}
