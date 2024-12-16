package sync

import (
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
)

// ParkedBranchProgram adds the opcodes to sync the parked branch with the given name.
func ParkedBranchProgram(syncStrategy configdomain.SyncStrategy, initialBranch gitdomain.LocalBranchName, args featureBranchArgs) {
	if args.localName == initialBranch {
		FeatureBranchProgram(syncStrategy, args)
	}
}
