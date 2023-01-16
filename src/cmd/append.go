package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

type appendConfig struct {
	ancestorBranches    []string
	hasOrigin           bool
	isOffline           bool
	noPushHook          bool
	parentBranch        string
	shouldNewBranchPush bool
	targetBranch        string
}

func addAppendCmd(rootCmd *cobra.Command) {
	appendCmd := &cobra.Command{
		Use:   "append <branch>",
		Short: "Creates a new feature branch as a child of the current branch",
		Long: `Creates a new feature branch as a direct child of the current branch.

Syncs the current branch,
forks a new feature branch with the given name off the current branch,
makes the new branch a child of the current branch,
pushes the new feature branch to the origin repository
(if and only if "push-new-branches" is true),
and brings over all uncommitted changes to the new feature branch.

See "sync" for information regarding upstream remotes.`,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := createAppendConfig(args, prodRepo)
			if err != nil {
				cli.Exit(err)
			}
			stepList, err := createAppendStepList(config, prodRepo)
			if err != nil {
				cli.Exit(err)
			}
			runState := runstate.New("append", stepList)
			err = runstate.Execute(runState, prodRepo, nil)
			if err != nil {
				cli.Exit(err)
			}
		},
		Args: cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(prodRepo); err != nil {
				return err
			}
			return validateIsConfigured(prodRepo)
		},
	}
	rootCmd.AddCommand(appendCmd)
}

func createAppendConfig(args []string, repo *git.ProdRepo) (appendConfig, error) {
	parentBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return appendConfig{}, err
	}
	result := appendConfig{
		parentBranch: parentBranch,
		targetBranch: args[0],
	}
	result.hasOrigin, err = repo.Silent.HasOrigin()
	if err != nil {
		return appendConfig{}, err
	}
	isOffline, err := repo.Config.IsOffline()
	if err != nil {
		return appendConfig{}, err
	}
	if result.hasOrigin && !isOffline {
		err := repo.Logging.Fetch()
		if err != nil {
			return appendConfig{}, err
		}
	}
	hasBranch, err := repo.Silent.HasLocalOrOriginBranch(result.targetBranch)
	if err != nil {
		return appendConfig{}, err
	}
	if hasBranch {
		return appendConfig{}, fmt.Errorf("a branch named %q already exists", result.targetBranch)
	}
	err = dialog.EnsureKnowsParentBranches([]string{result.parentBranch}, repo)
	if err != nil {
		return appendConfig{}, err
	}
	result.ancestorBranches = repo.Config.AncestorBranches(result.parentBranch)
	pushHook, err := repo.Config.PushHook()
	if err != nil {
		return appendConfig{}, err
	}
	result.noPushHook = !pushHook
	result.shouldNewBranchPush, err = repo.Config.ShouldNewBranchPush()
	if err != nil {
		return appendConfig{}, err
	}
	result.isOffline = isOffline
	return result, err
}
