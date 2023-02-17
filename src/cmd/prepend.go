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
	ancestorBranches    []string
	hasOrigin           bool
	initialBranch       string
	isOffline           bool
	noPushHook          bool
	parentBranch        string
	shouldNewBranchPush bool
	targetBranch        string
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

func determinePrependConfig(args []string, repo *git.ProdRepo) (prependConfig, error) {
	initialBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return prependConfig{}, err
	}
	result := prependConfig{
		initialBranch: initialBranch,
		targetBranch:  args[0],
	}
	result.hasOrigin, err = repo.Silent.HasOrigin()
	if err != nil {
		return prependConfig{}, err
	}
	result.shouldNewBranchPush, err = repo.Config.ShouldNewBranchPush()
	if err != nil {
		return prependConfig{}, err
	}
	isOffline, err := repo.Config.IsOffline()
	if err != nil {
		return prependConfig{}, err
	}
	result.isOffline = isOffline
	if result.hasOrigin && !result.isOffline {
		err := repo.Logging.Fetch()
		if err != nil {
			return prependConfig{}, err
		}
	}
	hasBranch, err := repo.Silent.HasLocalOrOriginBranch(result.targetBranch)
	if err != nil {
		return prependConfig{}, err
	}
	if hasBranch {
		return prependConfig{}, fmt.Errorf("a branch named %q already exists", result.targetBranch)
	}
	if !repo.Config.IsFeatureBranch(result.initialBranch) {
		return prependConfig{}, fmt.Errorf("the branch %q is not a feature branch. Only feature branches can have parent branches", result.initialBranch)
	}
	parentDialog := dialog.ParentBranches{}
	err = parentDialog.EnsureKnowsParentBranches([]string{result.initialBranch}, repo)
	if err != nil {
		return prependConfig{}, err
	}
	pushHook, err := repo.Config.PushHook()
	if err != nil {
		return prependConfig{}, err
	}
	result.noPushHook = !pushHook
	result.parentBranch = repo.Config.ParentBranch(result.initialBranch)
	result.ancestorBranches = repo.Config.AncestorBranches(result.initialBranch)
	return result, nil
}

func prependStepList(config prependConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	for _, branchName := range config.ancestorBranches {
		steps, err := syncBranchSteps(branchName, true, repo)
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
