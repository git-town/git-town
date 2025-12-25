package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// BranchTrackingCreateIfLocalExists pushes the given local branch up to origin
// and marks it as tracking the current branch,
// but only if the given branch still exists locally.
type BranchTrackingCreateIfLocalExists struct {
	Branch gitdomain.LocalBranchName
}

func (self *BranchTrackingCreateIfLocalExists) Run(args shared.RunArgs) error {
	if args.Git.BranchExists(args.Backend, self.Branch) {
		args.PrependOpcodes(&BranchTrackingCreate{Branch: self.Branch})
	}
	return nil
}
