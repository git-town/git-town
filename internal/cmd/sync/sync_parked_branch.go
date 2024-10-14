package sync

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
)

// ParkedBranchProgram adds the opcodes to sync the parked branch with the given name.
func ParkedBranchProgram(initialBranch gitdomain.LocalBranchName, args featureBranchArgs) {
	if args.localName == initialBranch {
		FeatureBranchProgram(args)
	}
}
