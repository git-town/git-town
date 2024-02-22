package sync

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/git/gitdomain"
)

// PerennialBranchProgram adds the opcodes to sync the observed branch with the given name.
func ParkedBranchProgram(initialBranch gitdomain.LocalBranchName, args featureBranchArgs) {
	fmt.Println("111111111111111")
	if args.branch.LocalName == initialBranch {
		FeatureBranchProgram(args)
	}
}
