package opcodes

import (
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// DeleteTrackingBranch deletes the tracking branch of the given local branch.
type DeleteTrackingBranch struct {
	Branch                  gitdomain.RemoteBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *DeleteTrackingBranch) Run(args shared.RunArgs) error {
	// TODO: add AllowError setting and if true, ignore the error here
	// this is useful when shipping via API
	// since there is nothing the user can do to retry this in a meaningful way
	return args.Git.DeleteTrackingBranch(args.Frontend, self.Branch)
}
