package sync

import (
	"github.com/git-town/git-town/v12/src/git/gitdomain"
)

// ParkedBranchProgram adds the opcodes to sync the observed branch with the given name.
func ParkedBranchProgram(initialBranch gitdomain.LocalBranchName, args featureBranchArgs) {
	if args.branch.LocalName == initialBranch {
		FeatureBranchProgram(args)
	}
}
