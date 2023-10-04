package domain

import (
	"sort"
	"strings"
)

type LocalBranchNames []LocalBranchName

func NewLocalBranchNames(names ...string) LocalBranchNames {
	result := make(LocalBranchNames, len(names))
	for n, name := range names {
		result[n] = NewLocalBranchName(name)
	}
	return result
}

func (lbns LocalBranchNames) Categorize(branchTypes BranchTypes) (perennials, features LocalBranchNames) {
	for _, branch := range lbns {
		if branchTypes.IsFeatureBranch(branch) {
			features = append(features, branch)
		} else {
			perennials = append(perennials, branch)
		}
	}
	return
}

// Join provides the names of all branches in this collection connected by the given separator.
func (lbns LocalBranchNames) Join(sep string) string {
	return strings.Join(lbns.Strings(), sep)
}

// Sort orders the branches in this collection alphabetically.
func (lbns LocalBranchNames) Sort() {
	sort.Slice(lbns, func(a, b int) bool {
		return lbns[a].ID < lbns[b].ID
	})
}

// Strings provides the names of all branches in this collection as strings.
func (lbns LocalBranchNames) Strings() []string {
	result := make([]string, len(lbns))
	for b, branch := range lbns {
		result[b] = branch.String()
	}
	return result
}
