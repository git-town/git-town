package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/git-town/git-town/src/script"
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
		repo := git.NewProdRepo()
		config, err := getHackConfig(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		stepList := getAppendStepList(config, repo)
		runState := steps.NewRunState("hack", stepList)
		err = steps.Run(runState, repo)
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
		return validateIsConfigured()
	},
}

func getParentBranch(targetBranch string) string {
	if promptForParent {
		parentBranch := prompt.AskForBranchParent(targetBranch, git.Config().GetMainBranch())
		prompt.EnsureKnowsParentBranches([]string{parentBranch})
		return parentBranch
	}
	return git.Config().GetMainBranch()
}

func getHackConfig(args []string) (result appendConfig, err error) {
	result.targetBranch = args[0]
	result.parentBranch = getParentBranch(result.targetBranch)
	result.hasOrigin = git.HasRemote("origin")
	result.shouldNewBranchPush = git.Config().ShouldNewBranchPush()
	result.isOffline = git.Config().IsOffline()
	if git.HasRemote("origin") && !git.Config().IsOffline() {
		err := script.Fetch()
		if err != nil {
			return result, err
		}
	}
	git.EnsureDoesNotHaveBranch(result.targetBranch)
	return
}

func init() {
	hackCmd.Flags().BoolVarP(&promptForParent, "prompt", "p", false, "Prompt for the parent branch")
	RootCmd.AddCommand(hackCmd)
}
