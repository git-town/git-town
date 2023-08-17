package domain

import (
	"fmt"
	"strings"
)

// RemoteBranchName is the name of a remote branch.
// Examples:
// - the local branch `foo` has the RemoteBranchName `origin/foo`
type RemoteBranchName struct {
	BranchName
}

func NewRemoteBranchName(value string) RemoteBranchName {
	if !isValidRemoteBranchName(value) {
		panic(fmt.Sprintf("%q is not a valid remote branch name", value))
	}
	return RemoteBranchName{BranchName{Location{value}}}
}

func isValidRemoteBranchName(value string) bool {
	if len(value) < 3 {
		return false
	}
	if !strings.Contains(value, "/") {
		return false
	}
	return true
}

func (r RemoteBranchName) LocalBranchName() LocalBranchName {
	return NewLocalBranchName(strings.TrimPrefix(r.value, "origin/"))
}

// Implements the fmt.Stringer interface.
func (r RemoteBranchName) String() string { return r.value }

// type RemoteBranchNames []RemoteBranchName

// func (r RemoteBranchNames) Join(sep string) string {
// 	return strings.Join(r.Strings(), sep)
// }

// func (r RemoteBranchNames) Sort() {
// 	sort.Slice(r, func(i, j int) bool {
// 		return r[i].value < r[j].value
// 	})
// }

// func (r RemoteBranchNames) Strings() []string {
// 	result := make([]string, len(r))
// 	for b, branch := range r {
// 		result[b] = branch.String()
// 	}
// 	return result
// }
