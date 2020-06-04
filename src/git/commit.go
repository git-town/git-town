package git

import (
	"fmt"
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

// Set assigns the given value to the property with the given Gherkin table name.
func (commit *Commit) Set(name, value string) (err error) {
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
