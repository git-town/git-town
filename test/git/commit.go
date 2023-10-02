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
func (c *Commit) Set(name, value string) {
	switch name {
	case "BRANCH":
		c.Branch = domain.NewLocalBranchName(value)
	case "LOCATION":
		c.Locations = []string{value}
	case "MESSAGE":
		c.Message = value
	case "FILE NAME":
		c.FileName = value
	case "FILE CONTENT":
		c.FileContent = value
	case "AUTHOR":
		c.Author = value
	default:
		log.Fatalf("unknown Commit property: %s", name)
	}
}
