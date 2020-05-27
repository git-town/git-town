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
		config, err := getAppendConfig(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		stepList := getAppendStepList(config)
		runState := steps.NewRunState("append", stepList)
		err = steps.Run(runState, git.NewProdRepo())
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

func getAppendConfig(args []string) (result appendConfig, err error) {
	result.parentBranch = git.GetCurrentBranchName()
	result.targetBranch = args[0]
	if git.HasRemote("origin") && !git.Config().IsOffline() {
		err := script.Fetch()
		if err != nil {
			return result, err
		}
	}
	git.EnsureDoesNotHaveBranch(result.targetBranch)
	prompt.EnsureKnowsParentBranches([]string{result.parentBranch})
	result.ancestorBranches = git.Config().GetAncestorBranches(result.parentBranch)
	result.hasOrigin = git.HasRemote("origin")
	result.shouldNewBranchPush = git.Config().ShouldNewBranchPush()
	result.isOffline = git.Config().IsOffline()
	return
}

func init() {
	RootCmd.AddCommand(appendCommand)
}
