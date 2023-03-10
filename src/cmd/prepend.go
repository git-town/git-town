package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	. "github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

func prependCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:     "prepend <branch>",
		GroupID: "lineage",
		Args:    cobra.ExactArgs(1),
		PreRunE: Ensure(repo, HasGitVersion, IsRepository, IsConfigured),
		Short:   "Creates a new feature branch as the parent of the current branch",
		Long: `Creates a new feature branch as the parent of the current branch

Syncs the parent branch,
cuts a new feature branch with the given name off the parent branch,
makes the new branch the parent of the current branch,
pushes the new feature branch to the origin repository
(if "push-new-branches" is true),
and brings over all uncommitted changes to the new feature branch.

See "sync" for upstream remote options.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := determinePrependConfig(args, repo)
			if err != nil {
				return err
			}
			stepList, err := prependStepList(config, repo)
			if err != nil {
				return err
			}
			runState := runstate.New("prepend", stepList)
			return runstate.Execute(runState, repo, nil)
		},
	}
}

type prependConfig struct {
	ancestorBranches    []string
	hasOrigin           bool
	initialBranch       string
	isOffline           bool
	noPushHook          bool
	parentBranch        string
	shouldNewBranchPush bool
	targetBranch        string
}

func determinePrependConfig(args []string, repo *git.ProdRepo) (*prependConfig, error) {
	ec := runstate.ErrorChecker{}
	initialBranch := ec.String(repo.Silent.CurrentBranch())
	hasOrigin := ec.Bool(repo.Silent.HasOrigin())
	shouldNewBranchPush := ec.Bool(repo.Config.ShouldNewBranchPush())
	pushHook := ec.Bool(repo.Config.PushHook())
	isOffline := ec.Bool(repo.Config.IsOffline())
	if ec.Err != nil {
		return nil, ec.Err
	}
	if hasOrigin && !isOffline {
		err := repo.Logging.Fetch()
		if err != nil {
			return nil, err
		}
	}
	targetBranch := args[0]
	hasBranch, err := repo.Silent.HasLocalOrOriginBranch(targetBranch)
	if err != nil {
		return nil, err
	}
	if hasBranch {
		return nil, fmt.Errorf("a branch named %q already exists", targetBranch)
	}
	if !repo.Config.IsFeatureBranch(initialBranch) {
		return nil, fmt.Errorf("the branch %q is not a feature branch. Only feature branches can have parent branches", initialBranch)
	}
	parentDialog := dialog.ParentBranches{}
	err = parentDialog.EnsureKnowsParentBranches([]string{initialBranch}, repo)
	if err != nil {
		return nil, err
	}
	return &prependConfig{
		hasOrigin:           hasOrigin,
		initialBranch:       initialBranch,
		isOffline:           isOffline,
		noPushHook:          !pushHook,
		parentBranch:        repo.Config.ParentBranch(initialBranch),
		ancestorBranches:    repo.Config.AncestorBranches(initialBranch),
		shouldNewBranchPush: shouldNewBranchPush,
		targetBranch:        targetBranch,
	}, nil
}

func prependStepList(config *prependConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branch := range config.ancestorBranches {
		updateBranchSteps(&list, branch, true, repo)
	}
	list.Add(&steps.CreateBranchStep{Branch: config.targetBranch, StartingPoint: config.parentBranch})
	list.Add(&steps.SetParentStep{Branch: config.targetBranch, ParentBranch: config.parentBranch})
	list.Add(&steps.SetParentStep{Branch: config.initialBranch, ParentBranch: config.targetBranch})
	list.Add(&steps.CheckoutStep{Branch: config.targetBranch})
	if config.hasOrigin && config.shouldNewBranchPush && !config.isOffline {
		list.Add(&steps.CreateTrackingBranchStep{Branch: config.targetBranch, NoPushHook: config.noPushHook})
	}
	list.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, repo)
	return list.Result()
}
