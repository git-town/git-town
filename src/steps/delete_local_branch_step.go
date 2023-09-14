package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// DeleteLocalBranchStep deletes the branch with the given name.
type DeleteLocalBranchStep struct {
	Branch    domain.LocalBranchName
	Parent    domain.Location
	Force     bool
	branchSHA domain.SHA `exhaustruct:"optional"`
	EmptyStep
}

func (step *DeleteLocalBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&CreateBranchStep{Branch: step.Branch, StartingPoint: step.branchSHA.Location()}}, nil
}

func (step *DeleteLocalBranchStep) Run(args RunArgs) error {
	var err error
	step.branchSHA, err = args.Run.Backend.SHAForBranch(step.Branch.BranchName())
	if err != nil {
		return err
	}
	hasUnmergedCommits, err := args.Run.Backend.BranchHasUnmergedCommits(step.Branch, step.Parent)
	if err != nil {
		return err
	}
	return args.Run.Frontend.DeleteLocalBranch(step.Branch, step.Force || hasUnmergedCommits)
}
