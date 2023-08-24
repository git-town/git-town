package hosting_test

import "github.com/git-town/git-town/v9/src/domain"

// emptySHAForBranch is a dummy implementation for hosting.SHAForBranchfunc to be used in tests.
func emptySHAForBranch(domain.BranchName) (domain.SHA, error) {
	return domain.SHA{}, nil
}
