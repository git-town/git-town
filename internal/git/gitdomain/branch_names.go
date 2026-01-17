package gitdomain

import "strings"

type BranchNames []BranchName

func (self BranchNames) Join(sep string) string {
	return strings.Join(self.Strings(), sep)
}

func (self BranchNames) LocalBranchNames() LocalBranchNames {
	result := LocalBranchNames{}
	for _, branchName := range self {
		if branchName.IsLocal() {
			result = append(result, branchName.LocalName())
		}
	}
	return result
}

func (self BranchNames) Strings() []string {
	result := make([]string, len(self))
	for i, branchName := range self {
		result[i] = branchName.String()
	}
	return result
}
