package opcode

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// ResetRemoteBranchToSHA sets the given remote branch to the given SHA,
// but only if it currently has a particular SHA.
type ResetRemoteBranchToSHA struct {
	Branch      domain.RemoteBranchName
	MustHaveSHA domain.SHA
	SetToSHA    domain.SHA
	undeclaredOpcodeMethods
}

func (op *ResetRemoteBranchToSHA) Run(args shared.RunArgs) error {
	currentSHA, err := args.Runner.Backend.SHAForBranch(op.Branch.BranchName())
	if err != nil {
		return err
	}
	if currentSHA != op.MustHaveSHA {
		return fmt.Errorf(messages.BranchHasWrongSHA, op.Branch, op.SetToSHA, op.MustHaveSHA, currentSHA)
	}
	return args.Runner.Frontend.ResetRemoteBranchToSHA(op.Branch, op.SetToSHA)
}
