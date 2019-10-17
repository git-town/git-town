package test

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Originate/git-town/test/gherkintools"
	"github.com/Originate/git-town/test/helpers"
)

// CommitListBuilder provides lists of existing commits in Git repositories.
type CommitListBuilder struct {

	// commits stores data about commits.
	//
	// Structure (key  =>  value):
	//   commit 1 SHA  =>  commit 1
	//   commit 2 SHA  =>  commit 2
	commits map[string]gherkintools.Commit

	// commitsInBranch stores which branches contain which commits.
	//
	// Structure (key  =>  value):
	//   branch 1 name  =>  [commit 1 SHA, commit 2 SHA]
	//   branch 2 name  =>  [commit 1 SHA, commit 3 SHA]
	commitsInBranch map[string]*helpers.OrderedStringSet

	// locations stores which commits occur in which repositories.
	//
	// Structure (key  =>  value):
	//   commit 1 SHA  =>  ["local"]
	//   commit 2 SHA  =>  ["local", "remote"]
	locations map[string]*helpers.OrderedStringSet
}

// NewCommitListBuilder provides a fully initialized instance of commitListBuilder.
func NewCommitListBuilder() CommitListBuilder {
	result := CommitListBuilder{
		commits:         make(map[string]gherkintools.Commit),
		commitsInBranch: make(map[string]*helpers.OrderedStringSet),
		locations:       make(map[string]*helpers.OrderedStringSet),
	}
	return result
}

// Add inserts the given commit from the given location into this list.
func (builder *CommitListBuilder) Add(commit gherkintools.Commit, location string) {
	builder.commits[commit.SHA] = commit

	_, exists := builder.commitsInBranch[commit.Branch]
	if !exists {
		builder.commitsInBranch[commit.Branch] = &helpers.OrderedStringSet{}
	}
	builder.commitsInBranch[commit.Branch].Add(commit.SHA)

	_, exists = builder.locations[commit.SHA]
	if !exists {
		builder.locations[commit.SHA] = &helpers.OrderedStringSet{}
	}
	builder.locations[commit.SHA].Add(location)
	fmt.Printf("COMMIT %s LOC %s\n", commit.Message, builder.locations[commit.SHA])
}

// branches provides the names of the branches known to this CommitListBuilder.
func (builder *CommitListBuilder) branches() []string {
	result := make([]string, 0, len(builder.commitsInBranch))
	for branch := range builder.commitsInBranch {
		result = append(result, branch)
	}
	return sort.StringSlice(result)
}

// Table provides the data accumulated by this CommitListBuilder as a Mortadella table.
func (builder *CommitListBuilder) Table(fields []string) (result gherkintools.Mortadella) {
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
					locations := builder.locations[SHA]
					row = append(row, strings.Join(locations.Slice(), ", "))
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
