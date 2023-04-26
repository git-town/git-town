package gherkin

import (
	"sort"
	"strings"

	"github.com/git-town/git-town/v8/test/git"
	"github.com/git-town/git-town/v8/test/helpers"
)

// CommitTableBuilder collects data about commits in Git repositories
// in the same way that our Gherkin tables describing commits in repos are organized.
type CommitTableBuilder struct {
	// commits stores data about commits.
	//
	// Structure:
	//   commit 1 SHA:  commit 1
	//   commit 2 SHA:  commit 2
	commits map[string]git.Commit

	// commitsInBranch stores which branches contain which commits.
	//
	// Structure:
	//   branch 1 name: [commit 1 SHA, commit 2 SHA]
	//   branch 2 name: [commit 1 SHA, commit 3 SHA]
	commitsInBranch map[string]helpers.OrderedStringSet

	// locations stores which commits occur in which repositories.
	//
	// Structure:
	//   commit 1 SHA + branch 1 name:  ["local"]
	//   commit 1 SHA + branch 2 name:  ["local", "origin"]
	locations map[string]helpers.OrderedStringSet
}

// NewCommitTableBuilder provides a fully initialized instance of CommitTableBuilder.
func NewCommitTableBuilder() CommitTableBuilder {
	result := CommitTableBuilder{
		commits:         make(map[string]git.Commit),
		commitsInBranch: make(map[string]helpers.OrderedStringSet),
		locations:       make(map[string]helpers.OrderedStringSet),
	}
	return result
}

// Add registers the given commit from the given location into this table.
func (builder *CommitTableBuilder) Add(commit git.Commit, location string) {
	builder.commits[commit.SHA] = commit
	commitsInBranch, exists := builder.commitsInBranch[commit.Branch]
	if exists {
		builder.commitsInBranch[commit.Branch] = commitsInBranch.Add(commit.SHA)
	} else {
		builder.commitsInBranch[commit.Branch] = helpers.NewOrderedStringSet(commit.SHA)
	}
	locationKey := commit.SHA + commit.Branch
	locations, exists := builder.locations[locationKey]
	if exists {
		builder.locations[locationKey] = locations.Add(location)
	} else {
		builder.locations[locationKey] = helpers.NewOrderedStringSet(location)
	}
}

// AddMany registers the given commits from the given location into this table.
func (builder *CommitTableBuilder) AddMany(commits []git.Commit, location string) {
	for _, commit := range commits {
		builder.Add(commit, location)
	}
}

// branches provides the names of the all branches known to this CommitTableBuilder,
// sorted alphabetically, with the main branch first.
func (builder *CommitTableBuilder) branches() []string {
	result := make([]string, 0, len(builder.commitsInBranch))
	hasMain := false
	for branch := range builder.commitsInBranch {
		if branch == "main" {
			hasMain = true
		} else {
			result = append(result, branch)
		}
	}
	sort.Strings(result)
	if hasMain {
		return append([]string{"main"}, result...)
	}
	return result
}

// Table provides the data accumulated by this CommitTableBuilder as a DataTable.
func (builder *CommitTableBuilder) Table(fields []string) DataTable {
	result := DataTable{}
	result.AddRow(fields...)
	lastBranch := ""
	lastLocation := ""
	for _, branch := range builder.branches() {
		SHAs := builder.commitsInBranch[branch]
		for _, SHA := range SHAs.Slice() {
			commit := builder.commits[SHA]
			row := []string{}
			for _, field := range fields {
				switch field {
				case "BRANCH":
					if branch == lastBranch {
						row = append(row, "")
					} else {
						row = append(row, branch)
					}
				case "LOCATION":
					locations := strings.Join(builder.locations[SHA+branch].Slice(), ", ")
					if locations == lastLocation && branch == lastBranch {
						row = append(row, "")
					} else {
						lastLocation = locations
						row = append(row, locations)
					}
				case "MESSAGE":
					row = append(row, commit.Message)
				case "FILE NAME":
					row = append(row, commit.FileName)
				case "FILE CONTENT":
					row = append(row, commit.FileContent)
				case "AUTHOR":
					row = append(row, commit.Author)
				default:
					panic("unknown table field: " + field)
				}
			}
			result.AddRow(row...)
			lastBranch = branch
		}
	}
	return result
}
