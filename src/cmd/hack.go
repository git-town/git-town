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
			config, err := determineHackConfig(args, promptForParentFlag, repo)
			if err != nil {
				cli.Exit(err)
			}
			stepList, err := appendStepList(config, repo)
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
		GroupID: "basic",
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

func determineHackConfig(args []string, promptForParent bool, repo *git.ProdRepo) (*appendConfig, error) {
	ec := runstate.ErrorChecker{}
	targetBranch := args[0]
	parentBranch := ec.String(determineParentBranch(targetBranch, promptForParent, repo))
	hasOrigin := ec.Bool(repo.Silent.HasOrigin())
	shouldNewBranchPush := ec.Bool(repo.Config.ShouldNewBranchPush())
	isOffline := ec.Bool(repo.Config.IsOffline())
	if ec.Err == nil && hasOrigin && !isOffline {
		ec.Check(repo.Logging.Fetch())
	}
	hasBranch := ec.Bool(repo.Silent.HasLocalOrOriginBranch(targetBranch))
	pushHook := ec.Bool(repo.Config.PushHook())
	if hasBranch {
		return nil, fmt.Errorf("a branch named %q already exists", targetBranch)
	}
	return &appendConfig{
		ancestorBranches:    []string{},
		targetBranch:        targetBranch,
		parentBranch:        parentBranch,
		hasOrigin:           hasOrigin,
		shouldNewBranchPush: shouldNewBranchPush,
		noPushHook:          !pushHook,
		isOffline:           isOffline,
	}, ec.Err
}
