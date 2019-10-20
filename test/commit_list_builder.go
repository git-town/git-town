package test

import (
	"sort"
	"strings"

	"github.com/Originate/git-town/test/gherkintools"
	"github.com/Originate/git-town/test/helpers"
)

// CommitTableBuilder collects data about commits in Git repositories
// in the same way that our Gherkin tables describing commits in repos are organized.
type CommitTableBuilder struct {

	// commits stores data about commits.
	//
	// Structure:
	//   commit 1 SHA:  commit 1
	//   commit 2 SHA:  commit 2
	commits map[string]gherkintools.Commit

	// commitsInBranch stores which branches contain which commits.
	//
	// Structure:
	//   branch 1 name: [commit 1 SHA, commit 2 SHA]
	//   branch 2 name: [commit 1 SHA, commit 3 SHA]
	commitsInBranch map[string]*helpers.OrderedStringSet

	// locations stores which commits occur in which repositories.
	//
	// Structure:
	//   commit 1 SHA + branch 1 name:  ["local"]
	//   commit 1 SHA + branch 2 name:  ["local", "remote"]
	locations map[string]*helpers.OrderedStringSet
}

// NewCommitTableBuilder provides a fully initialized instance of commitListBuilder.
func NewCommitTableBuilder() CommitTableBuilder {
	result := CommitTableBuilder{
		commits:         make(map[string]gherkintools.Commit),
		commitsInBranch: make(map[string]*helpers.OrderedStringSet),
		locations:       make(map[string]*helpers.OrderedStringSet),
	}
	return result
}

// Add inserts the given commit from the given location into this list.
func (builder *CommitTableBuilder) Add(commit gherkintools.Commit, location string) {
	builder.commits[commit.SHA] = commit

	_, exists := builder.commitsInBranch[commit.Branch]
	if !exists {
		builder.commitsInBranch[commit.Branch] = &helpers.OrderedStringSet{}
	}
	builder.commitsInBranch[commit.Branch].Add(commit.SHA)

	locationKey := commit.SHA + commit.Branch
	_, exists = builder.locations[locationKey]
	if !exists {
		builder.locations[locationKey] = &helpers.OrderedStringSet{}
	}
	builder.locations[locationKey].Add(location)
}

// branches provides the names of the branches known to this CommitListBuilder.
func (builder *CommitTableBuilder) branches() []string {
	result := make([]string, 0, len(builder.commitsInBranch))
	for branch := range builder.commitsInBranch {
		result = append(result, branch)
	}
	sort.Strings(result)
	return result
}

// Table provides the data accumulated by this CommitListBuilder as a Mortadella table.
func (builder *CommitTableBuilder) Table(fields []string) (result gherkintools.Mortadella) {
	result.AddRow(fields...)
	// Note: need to create a sorted list of branch names here,
	// because iterating builder.commitsInBranch directly provides the branch names in random order.
	for _, branch := range builder.branches() {
		SHAs := builder.commitsInBranch[branch]
		for _, SHA := range SHAs.Slice() {
			commit := builder.commits[SHA]
			row := []string{}
			for _, field := range fields {
				switch field {
				case "BRANCH":
					row = append(row, branch)
				case "LOCATION":
					locations := builder.locations[SHA+branch]
					row = append(row, strings.Join(locations.Slice(), ", "))
				case "MESSAGE":
					row = append(row, commit.Message)
				case "FILE NAME":
					row = append(row, commit.FileName)
				case "FILE CONTENT":
					row = append(row, commit.FileContent)
				default:
					panic("unknown table field: " + field)
				}
			}
			result.AddRow(row...)
		}
	}
	return result
}
