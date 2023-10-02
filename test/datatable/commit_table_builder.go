package datatable

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/test/git"
	"github.com/git-town/git-town/v9/test/helpers"
)

// CommitTableBuilder collects data about commits in Git repositories
// in the same way that our Gherkin tables describing commits in repos are organized.
type CommitTableBuilder struct {
	// commits stores data about commits.
	//
	// Structure:
	//   commit 1 SHA:  commit 1
	//   commit 2 SHA:  commit 2
	commits map[domain.SHA]git.Commit

	// commitsInBranch stores which branches contain which commits.
	//
	// Structure:
	//   branch 1 name: [commit 1 SHA, commit 2 SHA]
	//   branch 2 name: [commit 1 SHA, commit 3 SHA]
	commitsInBranch map[domain.LocalBranchName]helpers.OrderedSet[domain.SHA]

	// locations stores which commits occur in which repositories.
	//
	// Structure:
	//   commit 1 SHA + branch 1 name:  ["local"]
	//   commit 1 SHA + branch 2 name:  ["local", "origin"]
	locations map[string]helpers.OrderedSet[string]
}

// NewCommitTableBuilder provides a fully initialized instance of CommitTableBuilder.
func NewCommitTableBuilder() CommitTableBuilder {
	result := CommitTableBuilder{
		commits:         make(map[domain.SHA]git.Commit),
		commitsInBranch: make(map[domain.LocalBranchName]helpers.OrderedSet[domain.SHA]),
		locations:       make(map[string]helpers.OrderedSet[string]),
	}
	return result
}

// Add registers the given commit from the given location into this table.
func (ctb *CommitTableBuilder) Add(commit git.Commit, location string) {
	ctb.commits[commit.SHA] = commit
	commitsInBranch, exists := ctb.commitsInBranch[commit.Branch]
	if exists {
		ctb.commitsInBranch[commit.Branch] = commitsInBranch.Add(commit.SHA)
	} else {
		ctb.commitsInBranch[commit.Branch] = helpers.NewOrderedSet(commit.SHA)
	}
	locationKey := commit.SHA.String() + commit.Branch.String()
	locations, exists := ctb.locations[locationKey]
	if exists {
		ctb.locations[locationKey] = locations.Add(location)
	} else {
		ctb.locations[locationKey] = helpers.NewOrderedSet(location)
	}
}

// AddMany registers the given commits from the given location into this table.
func (ctb *CommitTableBuilder) AddMany(commits []git.Commit, location string) {
	for _, commit := range commits {
		ctb.Add(commit, location)
	}
}

// branches provides the names of the all branches known to this CommitTableBuilder,
// sorted alphabetically, with the main branch first.
func (ctb *CommitTableBuilder) branches() domain.LocalBranchNames {
	result := make(domain.LocalBranchNames, 0, len(ctb.commitsInBranch))
	hasMain := false
	for branch := range ctb.commitsInBranch {
		if branch == domain.NewLocalBranchName("main") {
			hasMain = true
		} else {
			result = append(result, branch)
		}
	}
	result.Sort()
	if hasMain {
		return append(domain.NewLocalBranchNames("main"), result...)
	}
	return result
}

// Table provides the data accumulated by this CommitTableBuilder as a DataTable.
func (ctb *CommitTableBuilder) Table(fields []string) DataTable {
	result := DataTable{}
	result.AddRow(fields...)
	lastBranch := domain.LocalBranchName{}
	lastLocation := ""
	for _, branch := range ctb.branches() {
		SHAs := ctb.commitsInBranch[branch]
		for _, SHA := range SHAs.Elements() {
			commit := ctb.commits[SHA]
			row := []string{}
			for _, field := range fields {
				switch field {
				case "BRANCH":
					if branch == lastBranch {
						row = append(row, "")
					} else {
						row = append(row, branch.String())
					}
				case "LOCATION":
					locations := ctb.locations[SHA.String()+branch.String()].Join(", ")
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
