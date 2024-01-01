package gitdomain

import "sort"

type RemoteBranchNames []RemoteBranchName

// Sort orders the branches in this collection alphabetically.
func (self RemoteBranchNames) Sort() {
	sort.Slice(self, func(a, b int) bool {
		return self[a] < self[b]
	})
}

// Strings provides these remote branch names as strings.
func (self RemoteBranchNames) Strings() []string {
	result := make([]string, len(self))
	for b, branch := range self {
		result[b] = branch.String()
	}
	return result
}
