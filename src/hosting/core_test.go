package hosting_test

import (
	"github.com/git-town/git-town/v9/src/git"
	testgit "github.com/git-town/git-town/v9/test/git"
)

// emptyShaForBranch is a dummy implementation for hosting.ShaForBranchfunc to be used in tests.
func emptyShaForBranch(string) (git.SHA, error) {
	return testgit.ZeroValueSHA(), nil
}
