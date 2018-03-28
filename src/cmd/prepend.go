package cmd

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/prompt"
	"github.com/Originate/git-town/src/script"
	"github.com/Originate/git-town/src/steps"
	"github.com/Originate/git-town/src/util"

	"github.com/spf13/cobra"
)

type prependConfig struct {
	InitialBranch string
	ParentBranch  string
	TargetBranch  string
}

var prependCommand = &cobra.Command{
	Use:   "prepend <branch>",
	Short: "Creates a new feature branch as the parent of the current branch",
	Long: `Creates a new feature branch as the parent of the current branch

Syncs the parent branch
forks a new feature branch with the given name off the parent branch,
makes the new branch the parent of the current branch,
pushes the new feature branch to the remote repository,
and brings over all uncommitted changes to the new feature branch.

Additionally, when there is a remote upstream,
the main branch is synced with its upstream counterpart.
This can be disabled by toggling the "new-branch-push-flag" configuration:

	git town new-branch-push-flag false`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getPrependConfig(args)
		stepList := getPrependStepList(config)
		runState := steps.NewRunState("prepend", stepList)
		steps.Run(runState)
	},
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
		)
	},
}

func getPrependConfig(args []string) (result prependConfig) {
	result.InitialBranch = git.GetCurrentBranchName()
	result.TargetBranch = args[0]
	if git.HasRemote("origin") && !git.IsOffline() {
		script.Fetch()
	}
	git.EnsureDoesNotHaveBranch(result.TargetBranch)
	git.EnsureIsFeatureBranch(result.InitialBranch, "Only feature branches can have parent branches.")
	prompt.EnsureKnowsParentBranches([]string{result.InitialBranch})
	result.ParentBranch = git.GetParentBranch(result.InitialBranch)
	return
}

func getPrependStepList(config prependConfig) (result steps.StepList) {
	for _, branchName := range git.GetAncestorBranches(config.InitialBranch) {
		result.AppendList(steps.GetSyncBranchSteps(branchName, true))
	}
	result.Append(&steps.CreateBranchStep{BranchName: config.TargetBranch, StartingPoint: config.ParentBranch})
	result.Append(&steps.SetParentBranchStep{BranchName: config.TargetBranch, ParentBranchName: config.ParentBranch})
	result.Append(&steps.SetParentBranchStep{BranchName: config.InitialBranch, ParentBranchName: config.TargetBranch})
	result.Append(&steps.CheckoutBranchStep{BranchName: config.TargetBranch})
	if git.HasRemote("origin") && git.ShouldNewBranchPush() && !git.IsOffline() {
		result.Append(&steps.CreateTrackingBranchStep{BranchName: config.TargetBranch})
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
	return
}

func init() {
	RootCmd.AddCommand(prependCommand)
}
