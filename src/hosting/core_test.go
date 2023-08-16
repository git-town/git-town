package hosting_test

import "github.com/git-town/git-town/v9/src/domain"

// emptyShaForBranch is a dummy implementation for hosting.ShaForBranchfunc to be used in tests.
func emptyShaForBranch(string) (domain.SHA, error) {
	return domain.SHA{}, nil
}
