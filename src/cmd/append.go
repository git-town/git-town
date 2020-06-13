package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/git-town/git-town/src/steps"

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
			fmt.Println(err)
			os.Exit(1)
		}
		stepList, err := getAppendStepList(config, prodRepo)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		runState := steps.NewRunState("append", stepList)
		err = steps.Run(runState, prodRepo, nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := git.ValidateIsRepository(); err != nil {
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
	hasRemote, err := repo.Silent.HasRemote("origin")
	if err != nil {
		return result, err
	}
	if hasRemote && !git.Config().IsOffline() {
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
	result.ancestorBranches = git.Config().GetAncestorBranches(result.parentBranch)
	result.hasOrigin, err = repo.Silent.HasRemote("origin")
	if err != nil {
		return result, err
	}
	result.shouldNewBranchPush = git.Config().ShouldNewBranchPush()
	result.isOffline = git.Config().IsOffline()
	return result, err
}

func init() {
	RootCmd.AddCommand(appendCommand)
}
