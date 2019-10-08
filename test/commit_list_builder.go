package test

import (
	"fmt"
	"strings"

	"github.com/Originate/git-town/test/gherkintools"
)

// CommitListBuilder provides lists of existing commits in Git repositories.
type CommitListBuilder struct {

	// commits stores which branches contain which commits.
	// The first key is the branch name, the second key is the commit SHA.
	commits map[string]map[string]*gherkintools.Commit
}

// NewCommitListBuilder provides a fully initialized instance of commitListBuilder.
func NewCommitListBuilder() CommitListBuilder {
	return CommitListBuilder{commits: make(map[string]map[string]*gherkintools.Commit)}
}

// Add inserts the given commit into this list.
func (builder *CommitListBuilder) Add(commit gherkintools.Commit) {
	fmt.Println("7777777777777777777", commit)
	_, exists := builder.commits[commit.Branch]
	if !exists {
		builder.commits[commit.Branch] = make(map[string]*gherkintools.Commit)
	}
	_, exists = builder.commits[commit.Branch][commit.SHA]
	if exists {
		c := builder.commits[commit.Branch][commit.SHA]
		c.Location = append(c.Location, commit.Location...)
	} else {
		builder.commits[commit.Branch][commit.SHA] = &commit
	}
}

// AddAll inserts all given commits into this list.
func (builder *CommitListBuilder) AddAll(commits []gherkintools.Commit) {
	for _, commit := range commits {
		builder.Add(commit)
	}
}

// Table provides the data accumulated by this CommitListBuilder as a Mortadella table.
func (builder *CommitListBuilder) Table(fields []string) (result gherkintools.Mortadella) {
	result.AddRow(fields...)
	for _, commitMap := range builder.commits {
		for _, commit := range commitMap {
			row := []string{}
			for _, field := range fields {
				switch field {
				case "BRANCH":
					row = append(row, commit.Branch)
				case "LOCATION":
					row = append(row, strings.Join(commit.Location, ", "))
				case "MESSAGE":
					row = append(row, commit.Message)
				default:
					panic("unknown table field: " + field)
				}
			}
			result.AddRow(row...)
		}
	}
	return result
}
