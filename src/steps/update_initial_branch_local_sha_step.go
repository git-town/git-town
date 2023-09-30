package steps

import "github.com/git-town/git-town/v9/src/domain"

type UpdateInitialBranchLocalSHAStep struct {
	Branch    domain.LocalBranchName
	EmptyStep `exhaustruct:"optional"`
}

func (step *UpdateInitialBranchLocalSHAStep) Run(args RunArgs) error {
	newSHA, err := args.Runner.Backend.SHAForBranch(step.Branch.BranchName())
	if err != nil {
		return err
	}
	return args.UpdateInitialBranchLocalSHA(step.Branch, newSHA)
}
