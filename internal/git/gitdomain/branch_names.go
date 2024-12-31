package gitdomain

import "strings"

type BranchNames []BranchName

func (self BranchNames) Join(sep string) string {
	return strings.Join(self.Strings(), sep)
}

func (self BranchNames) Strings() []string {
	result := make([]string, len(self))
	for i, branchName := range self {
		result[i] = branchName.String()
	}
	return result
}
