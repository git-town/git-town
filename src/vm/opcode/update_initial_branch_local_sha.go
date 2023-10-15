package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

type UpdateInitialBranchLocalSHA struct {
	Branch                  domain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (step *UpdateInitialBranchLocalSHA) Run(args shared.RunArgs) error {
	newSHA, err := args.Runner.Backend.SHAForBranch(step.Branch.BranchName())
	if err != nil {
		return err
	}
	return args.UpdateInitialBranchLocalSHA(step.Branch, newSHA)
}
