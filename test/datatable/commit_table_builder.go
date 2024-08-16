package datatable

import (
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/test/git"
	"github.com/git-town/git-town/v15/test/helpers"
)

// CommitTableBuilder collects data about commits in Git repositories
// in the same way that our Gherkin tables describing commits in repos are organized.
type CommitTableBuilder struct {
	// commits stores data about commits.
	//
	// Structure:
	//   commit 1 SHA:  commit 1
	//   commit 2 SHA:  commit 2
	commits map[gitdomain.SHA]git.Commit

	// commitsInBranch stores which branches contain which commits.
	//
	// Structure:
	//   branch 1 name: [commit 1 SHA, commit 2 SHA]
	//   branch 2 name: [commit 1 SHA, commit 3 SHA]
	commitsInBranch map[gitdomain.LocalBranchName]helpers.OrderedSet[gitdomain.SHA]

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
		commits:         make(map[gitdomain.SHA]git.Commit),
		commitsInBranch: make(map[gitdomain.LocalBranchName]helpers.OrderedSet[gitdomain.SHA]),
		locations:       make(map[string]helpers.OrderedSet[string]),
	}
	return result
}

// Add registers the given commit from the given location into this table.
func (self *CommitTableBuilder) Add(commit git.Commit, location string) {
	self.commits[commit.SHA] = commit
	commitsInBranch, exists := self.commitsInBranch[commit.Branch]
	if exists {
		self.commitsInBranch[commit.Branch] = commitsInBranch.Add(commit.SHA)
	} else {
		self.commitsInBranch[commit.Branch] = helpers.NewOrderedSet(commit.SHA)
	}
	locationKey := commit.SHA.String() + commit.Branch.String()
	locations, exists := self.locations[locationKey]
	if exists {
		self.locations[locationKey] = locations.Add(location)
	} else {
		self.locations[locationKey] = helpers.NewOrderedSet(location)
	}
}

// AddMany registers the given commits from the given location into this table.
func (self *CommitTableBuilder) AddMany(commits []git.Commit, location string) {
	for _, commit := range commits {
		self.Add(commit, location)
	}
}

// Table provides the data accumulated by this CommitTableBuilder as a DataTable.
func (self *CommitTableBuilder) Table(fields []string) DataTable {
	result := DataTable{}
	result.AddRow(fields...)
	var lastBranch gitdomain.LocalBranchName
	lastLocation := ""
	for _, branch := range self.branches() {
		SHAs := self.commitsInBranch[branch]
		for _, SHA := range SHAs.Elements() {
			commit := self.commits[SHA]
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
					locations := self.locations[SHA.String()+branch.String()].Join(", ")
					if locations == lastLocation && branch == lastBranch {
						row = append(row, "")
					} else {
						lastLocation = locations
						row = append(row, locations)
					}
				case "MESSAGE":
					row = append(row, commit.Message.String())
				case "FILE NAME":
					row = append(row, commit.FileName)
				case "FILE CONTENT":
					row = append(row, commit.FileContent)
				case "AUTHOR":
					row = append(row, commit.Author.String())
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

// branches provides the names of the all branches known to this CommitTableBuilder,
// sorted alphabetically, with the main branch first.
func (self *CommitTableBuilder) branches() gitdomain.LocalBranchNames {
	result := make(gitdomain.LocalBranchNames, 0, len(self.commitsInBranch))
	hasMain := false
	for branch := range self.commitsInBranch {
		if branch == gitdomain.NewLocalBranchName("main") {
			hasMain = true
		} else {
			result = append(result, branch)
		}
	}
	result.Sort()
	if hasMain {
		return append(gitdomain.NewLocalBranchNames("main"), result...)
	}
	return result
}
