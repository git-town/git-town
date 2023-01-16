package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func hackCmd(repo *git.ProdRepo) *cobra.Command {
	promptForParentFlag := false
	hackCmd := cobra.Command{
		Use:   "hack <branch>",
		Short: "Creates a new feature branch off the main development branch",
		Long: `Creates a new feature branch off the main development branch

Syncs the main branch,
forks a new feature branch with the given name off the main branch,
pushes the new feature branch to origin
(if and only if "push-new-branches" is true),
and brings over all uncommitted changes to the new feature branch.

See "sync" for information regarding upstream remotes.`,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := createHackConfig(args, promptForParentFlag, repo)
			if err != nil {
				cli.Exit(err)
			}
			stepList, err := createAppendStepList(config, repo)
			if err != nil {
				cli.Exit(err)
			}
			runState := runstate.New("hack", stepList)
			err = runstate.Execute(runState, repo, nil)
			if err != nil {
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
	hackCmd.Flags().BoolVarP(&promptForParentFlag, "prompt", "p", false, "Prompt for the parent branch")
	return &hackCmd
}

func determineParentBranch(targetBranch string, promptForParent bool, repo *git.ProdRepo) (string, error) {
	if promptForParent {
		parentDialog := dialog.ParentBranches{}
		parentBranch, err := parentDialog.AskForBranchParent(targetBranch, repo.Config.MainBranch(), repo)
		if err != nil {
			return "", err
		}
		err = parentDialog.EnsureKnowsParentBranches([]string{parentBranch}, repo)
		if err != nil {
			return "", err
		}
		return parentBranch, nil
	}
	return repo.Config.MainBranch(), nil
}

func createHackConfig(args []string, promptForParent bool, repo *git.ProdRepo) (appendConfig, error) {
	targetBranch := args[0]
	parentBranch, err := determineParentBranch(targetBranch, promptForParent, repo)
	if err != nil {
		return appendConfig{}, err
	}
	hasOrigin, err := repo.Silent.HasOrigin()
	if err != nil {
		return appendConfig{}, err
	}
	shouldNewBranchPush, err := repo.Config.ShouldNewBranchPush()
	if err != nil {
		return appendConfig{}, err
	}
	isOffline, err := repo.Config.IsOffline()
	if err != nil {
		return appendConfig{}, err
	}
	if hasOrigin && !isOffline {
		err := repo.Logging.Fetch()
		if err != nil {
			return appendConfig{}, err
		}
	}
	hasBranch, err := repo.Silent.HasLocalOrOriginBranch(targetBranch)
	if err != nil {
		return appendConfig{}, err
	}
	if hasBranch {
		return appendConfig{}, fmt.Errorf("a branch named %q already exists", targetBranch)
	}
	pushHook, err := repo.Config.PushHook()
	if err != nil {
		return appendConfig{}, err
	}
	result := appendConfig{
		targetBranch:        targetBranch,
		parentBranch:        parentBranch,
		hasOrigin:           hasOrigin,
		shouldNewBranchPush: shouldNewBranchPush,
		noPushHook:          !pushHook,
		isOffline:           isOffline,
	}
	return result, nil
}
