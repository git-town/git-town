package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/git-town/git-town/src/steps"

	"github.com/spf13/cobra"
)

var promptForParent bool

var hackCmd = &cobra.Command{
	Use:   "hack <branch>",
	Short: "Creates a new feature branch off the main development branch",
	Long: `Creates a new feature branch off the main development branch

Syncs the main branch,
forks a new feature branch with the given name off the main branch,
pushes the new feature branch to the remote repository
(if and only if "new-branch-push-flag" is true),
and brings over all uncommitted changes to the new feature branch.

See "sync" for information regarding remote upstream.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := getHackConfig(args, prodRepo)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		stepList, err := getAppendStepList(config, prodRepo)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		runState := steps.NewRunState("hack", stepList)
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

func getParentBranch(targetBranch string, repo *git.ProdRepo) (string, error) {
	if promptForParent {
		parentBranch, err := prompt.AskForBranchParent(targetBranch, git.Config().GetMainBranch(), repo)
		if err != nil {
			return "", err
		}
		err = prompt.EnsureKnowsParentBranches([]string{parentBranch}, repo)
		if err != nil {
			return "", err
		}
		return parentBranch, nil
	}
	return git.Config().GetMainBranch(), nil
}

func getHackConfig(args []string, repo *git.ProdRepo) (result appendConfig, err error) {
	result.targetBranch = args[0]
	result.parentBranch, err = getParentBranch(result.targetBranch, repo)
	if err != nil {
		return result, err
	}
	result.hasOrigin, err = repo.Silent.HasRemote("origin")
	if err != nil {
		return result, err
	}
	result.shouldNewBranchPush = git.Config().ShouldNewBranchPush()
	result.isOffline = git.Config().IsOffline()
	if result.hasOrigin && !git.Config().IsOffline() {
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
	return
}

func init() {
	hackCmd.Flags().BoolVarP(&promptForParent, "prompt", "p", false, "Prompt for the parent branch")
	RootCmd.AddCommand(hackCmd)
}
