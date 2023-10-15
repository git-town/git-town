package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

type UpdateInitialBranchLocalSHA struct {
	Branch                  domain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (op *UpdateInitialBranchLocalSHA) Run(args shared.RunArgs) error {
	newSHA, err := args.Runner.Backend.SHAForBranch(op.Branch.BranchName())
	if err != nil {
		return err
	}
	return args.UpdateInitialBranchLocalSHA(op.Branch, newSHA)
}
