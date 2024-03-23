package sync

import (
	"github.com/git-town/git-town/v13/src/git/gitdomain"
)

// ParkedBranchProgram adds the opcodes to sync the parked branch with the given name.
func ParkedBranchProgram(initialBranch gitdomain.LocalBranchName, args featureBranchArgs) {
	if args.branch.LocalName == initialBranch {
		FeatureBranchProgram(args)
	}
}
