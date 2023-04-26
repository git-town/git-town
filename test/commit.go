package test

import (
	"fmt"
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
