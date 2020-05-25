package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/git-town/git-town/src/script"
	"github.com/git-town/git-town/src/steps"
	"github.com/git-town/git-town/src/util"

	"github.com/spf13/cobra"
)

type prependConfig struct {
	initialBranch       string
	parentBranch        string
	targetBranch        string
	ancestorBranches    []string
	hasOrigin           bool
	shouldNewBranchPush bool
	isOffline           bool
}

var prependCommand = &cobra.Command{
	Use:   "prepend <branch>",
	Short: "Creates a new feature branch as the parent of the current branch",
	Long: `Creates a new feature branch as the parent of the current branch

Syncs the parent branch,
cuts a new feature branch with the given name off the parent branch,
makes the new branch the parent of the current branch,
pushes the new feature branch to the remote repository
(if "new-branch-push-flag" is true),
and brings over all uncommitted changes to the new feature branch.

See "sync" for remote upstream options.
`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := getPrependConfig(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		stepList := getPrependStepList(config)
		runState := steps.NewRunState("prepend", stepList)
		err = steps.Run(runState)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
		)
	},
}

func getPrependConfig(args []string) (result prependConfig, err error) {
	result.initialBranch = git.GetCurrentBranchName()
	result.targetBranch = args[0]
	result.hasOrigin = git.HasRemote("origin")
	result.shouldNewBranchPush = git.Config().ShouldNewBranchPush()
	result.isOffline = git.Config().IsOffline()
	if result.hasOrigin && !result.isOffline {
		err := script.Fetch()
		if err != nil {
			return result, err
		}
	}
	git.EnsureDoesNotHaveBranch(result.targetBranch)
	git.Config().EnsureIsFeatureBranch(result.initialBranch, "Only feature branches can have parent branches.")
	prompt.EnsureKnowsParentBranches([]string{result.initialBranch})
	result.parentBranch = git.Config().GetParentBranch(result.initialBranch)
	result.ancestorBranches = git.Config().GetAncestorBranches(result.initialBranch)
	return
}

func getPrependStepList(config prependConfig) (result steps.StepList) {
	for _, branchName := range config.ancestorBranches {
		result.AppendList(steps.GetSyncBranchSteps(branchName, true))
	}
	result.Append(&steps.CreateBranchStep{BranchName: config.targetBranch, StartingPoint: config.parentBranch})
	result.Append(&steps.SetParentBranchStep{BranchName: config.targetBranch, ParentBranchName: config.parentBranch})
	result.Append(&steps.SetParentBranchStep{BranchName: config.initialBranch, ParentBranchName: config.targetBranch})
	result.Append(&steps.CheckoutBranchStep{BranchName: config.targetBranch})
	if config.hasOrigin && config.shouldNewBranchPush && !config.isOffline {
		result.Append(&steps.CreateTrackingBranchStep{BranchName: config.targetBranch})
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
	return
}

func init() {
	RootCmd.AddCommand(prependCommand)
}
