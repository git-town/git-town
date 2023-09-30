package domain

import "sort"

type RemoteBranchNames []RemoteBranchName

// Strings provides these remote branch names as strings.
func (r RemoteBranchNames) Strings() []string {
	result := make([]string, len(r))
	for b, branch := range r {
		result[b] = branch.String()
	}
	return result
}

// Sort orders the branches in this collection alphabetically.
func (l RemoteBranchNames) Sort() {
	sort.Slice(l, func(a, b int) bool {
		return l[a].id < l[b].id
	})
}
