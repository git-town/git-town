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

func (l LocalBranchNames) Categorize(branchTypes BranchTypes) (perennials, features LocalBranchNames) {
	for _, branch := range l {
		if branchTypes.IsFeatureBranch(branch) {
			features = append(features, branch)
		} else {
			perennials = append(perennials, branch)
		}
	}
	return perennials, features
}

// Join provides the names of all branches in this collection connected by the given separator.
func (l LocalBranchNames) Join(sep string) string {
	return strings.Join(l.Strings(), sep)
}

// Sort orders the branches in this collection alphabetically.
func (l LocalBranchNames) Sort() {
	sort.Slice(l, func(a, b int) bool {
		return l[a].id < l[b].id
	})
}

// Strings provides the names of all branches in this collection as strings.
func (l LocalBranchNames) Strings() []string {
	result := make([]string, len(l))
	for b, branch := range l {
		result[b] = branch.String()
	}
	return result
}
