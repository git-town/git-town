package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/prompt"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/spf13/cobra"
)

type appendConfig struct {
	ancestorBranches    []string
	parentBranch        string
	targetBranch        string
	hasOrigin           bool
	isOffline           bool
	shouldNewBranchPush bool
}

var appendCommand = &cobra.Command{
	Use:   "append <branch>",
	Short: "Creates a new feature branch as a child of the current branch",
	Long: `Creates a new feature branch as a direct child of the current branch.

Syncs the current branch,
forks a new feature branch with the given name off the current branch,
makes the new branch a child of the current branch,
pushes the new feature branch to the remote repository
(if and only if "new-branch-push-flag" is true),
and brings over all uncommitted changes to the new feature branch.

See "sync" for information regarding remote upstream.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := getAppendConfig(args, prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		stepList, err := getAppendStepList(config, prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		runState := steps.NewRunState("append", stepList)
		err = steps.Run(runState, prodRepo, nil)
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

func getAppendConfig(args []string, repo *git.ProdRepo) (result appendConfig, err error) {
	result.parentBranch, err = repo.Silent.CurrentBranch()
	if err != nil {
		return result, err
	}
	result.targetBranch = args[0]
	result.hasOrigin, err = repo.Silent.HasRemote("origin")
	if err != nil {
		return result, err
	}
	if result.hasOrigin && !repo.Config.IsOffline() {
		err := repo.Logging.Fetch()
		if err != nil {
			return result, err
		}
	}
	hasBranch, err := repo.Silent.HasLocalOrRemoteBranch(result.targetBranch)
	if err != nil {
		return result, err
	}
	if hasBranch {
		return result, fmt.Errorf("a branch named %q already exists", result.targetBranch)
	}
	err = prompt.EnsureKnowsParentBranches([]string{result.parentBranch}, repo)
	if err != nil {
		return result, err
	}
	result.ancestorBranches = repo.Config.GetAncestorBranches(result.parentBranch)
	result.shouldNewBranchPush = repo.Config.ShouldNewBranchPush()
	result.isOffline = repo.Config.IsOffline()
	return result, err
}

func init() {
	RootCmd.AddCommand(appendCommand)
}
