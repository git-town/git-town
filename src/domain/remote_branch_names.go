package domain

import "sort"

type RemoteBranchNames []RemoteBranchName

// Sort orders the branches in this collection alphabetically.
func (rbns RemoteBranchNames) Sort() {
	sort.Slice(rbns, func(a, b int) bool {
		return rbns[a].ID < rbns[b].ID
	})
}

// Strings provides these remote branch names as strings.
func (rbns RemoteBranchNames) Strings() []string {
	result := make([]string, len(rbns))
	for b, branch := range rbns {
		result[b] = branch.String()
	}
	return result
}
