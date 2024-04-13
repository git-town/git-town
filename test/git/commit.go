package git

import (
	"log"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
)

// Commit describes a Git commit.
type Commit struct {
	Author      string `exhaustruct:"optional"`
	Branch      gitdomain.LocalBranchName
	FileContent string    `exhaustruct:"optional"`
	FileName    string    `exhaustruct:"optional"`
	Locations   Locations `exhaustruct:"optional"`
	Message     string
	SHA         gitdomain.SHA `exhaustruct:"optional"`
}

// Set assigns the given value to the property with the given Gherkin table name.
func (self *Commit) Set(name, value string) {
	switch name {
	case "BRANCH":
		self.Branch = gitdomain.NewLocalBranchName(value)
	case "LOCATION":
		self.Locations = NewLocations(value)
	case "MESSAGE":
		self.Message = value
	case "FILE NAME":
		self.FileName = value
	case "FILE CONTENT":
		self.FileContent = value
	case "AUTHOR":
		self.Author = value
	default:
		log.Fatalf("unknown Commit property: %s", name)
	}
}
