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
func (c RemoteBranchName) String() string { return c.value }
