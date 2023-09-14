package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// DeleteRemoteBranchStep deletes the tracking branch of the given local branch.
type DeleteRemoteBranchStep struct {
	Branch     domain.LocalBranchName
	NoPushHook bool
	branchSHA  domain.SHA `exhaustruct:"optional"`
	EmptyStep
}

func (step *DeleteRemoteBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&CreateRemoteBranchStep{Branch: step.Branch, SHA: step.branchSHA, NoPushHook: step.NoPushHook}}, nil
}

func (step *DeleteRemoteBranchStep) Run(args RunArgs) error {
	remoteBranch := step.Branch.AtRemote(domain.OriginRemote)
	var err error
	step.branchSHA, err = args.Runner.Backend.SHAForBranch(remoteBranch.BranchName())
	if err != nil {
		return err
	}
	return args.Runner.Frontend.DeleteRemoteBranch(step.Branch)
}
