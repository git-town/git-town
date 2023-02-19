package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/spf13/cobra"
)

type prependConfig struct {
	ancestorBranches        []string
	branchesDeletedOnRemote []string // local branches whose tracking branches have been deleted
	hasOrigin               bool
	initialBranch           string
	isOffline               bool
	noPushHook              bool
	parentBranch            string
	shouldNewBranchPush     bool
	targetBranch            string
}

func prependCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "prepend <branch>",
		Short: "Creates a new feature branch as the parent of the current branch",
		Long: `Creates a new feature branch as the parent of the current branch

Syncs the parent branch,
cuts a new feature branch with the given name off the parent branch,
makes the new branch the parent of the current branch,
pushes the new feature branch to the origin repository
(if "push-new-branches" is true),
and brings over all uncommitted changes to the new feature branch.

See "sync" for upstream remote options.
`,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := determinePrependConfig(args, repo)
			if err != nil {
				cli.Exit(err)
			}
			stepList, err := prependStepList(config, repo)
			if err != nil {
				cli.Exit(err)
			}
			runState := runstate.New("prepend", stepList)
			err = runstate.Execute(runState, repo, nil)
			if err != nil {
				fmt.Println(err)
				cli.Exit(err)
			}
		},
		Args: cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(repo); err != nil {
				return err
			}
			return validateIsConfigured(repo)
		},
	}
}

func determinePrependConfig(args []string, repo *git.ProdRepo) (*prependConfig, error) {
	initialBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return nil, err
	}
	hasOrigin, err := repo.Silent.HasOrigin()
	if err != nil {
		return nil, err
	}
	shouldNewBranchPush, err := repo.Config.ShouldNewBranchPush()
	if err != nil {
		return nil, err
	}
	isOffline, err := repo.Config.IsOffline()
	if err != nil {
		return nil, err
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
	pushHook, err := repo.Config.PushHook()
	if err != nil {
		return nil, err
	}
	branchesDeletedOnRemote, err := repo.Silent.LocalBranchesWithDeletedTrackingBranches()
	if err != nil {
		return nil, err
	}
	return &prependConfig{
		branchesDeletedOnRemote: branchesDeletedOnRemote,
		hasOrigin:               hasOrigin,
		initialBranch:           initialBranch,
		isOffline:               isOffline,
		noPushHook:              !pushHook,
		parentBranch:            repo.Config.ParentBranch(initialBranch),
		ancestorBranches:        repo.Config.AncestorBranches(initialBranch),
		shouldNewBranchPush:     shouldNewBranchPush,
		targetBranch:            targetBranch,
	}, nil
}

func prependStepList(config *prependConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	for _, branch := range config.ancestorBranches {
		steps, err := updateBranchSteps(branch, true, config.branchesDeletedOnRemote, repo)
		if err != nil {
			return runstate.StepList{}, err
		}
		result.AppendList(steps)
	}
	result.Append(&steps.CreateBranchStep{Branch: config.targetBranch, StartingPoint: config.parentBranch})
	result.Append(&steps.SetParentBranchStep{Branch: config.targetBranch, ParentBranch: config.parentBranch})
	result.Append(&steps.SetParentBranchStep{Branch: config.initialBranch, ParentBranch: config.targetBranch})
	result.Append(&steps.CheckoutBranchStep{Branch: config.targetBranch})
	if config.hasOrigin && config.shouldNewBranchPush && !config.isOffline {
		result.Append(&steps.CreateTrackingBranchStep{Branch: config.targetBranch, NoPushHook: config.noPushHook})
	}
	err := result.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, repo)
	return result, err
}
