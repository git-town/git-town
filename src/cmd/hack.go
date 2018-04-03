package cmd

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/script"
	"github.com/Originate/git-town/src/steps"
	"github.com/Originate/git-town/src/util"

	"github.com/spf13/cobra"
)

type hackConfig struct {
	TargetBranch string
}

var hackCmd = &cobra.Command{
	Use:   "hack <branch>",
	Short: "Creates a new feature branch off the main development branch",
	Long: `Creates a new feature branch off the main development branch

Syncs the main branch and forks a new feature branch with the given name off it.

If (and only if) [new-branch-push-flag](./new-branch-push-flag.md) is true,
pushes the new feature branch to the remote repository.

Finally, brings over all uncommitted changes to the new feature branch.

Additionally, when there is a remote upstream,
the main branch is synced with its upstream counterpart.
This can be disabled by toggling the "new-branch-push-flag" configuration:
$ git town new-branch-push-flag false`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getHackConfig(args)
		stepList := getHackStepList(config)
		runState := steps.NewRunState("hack", stepList)
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

func getHackConfig(args []string) (result hackConfig) {
	result.TargetBranch = args[0]
	if git.HasRemote("origin") && !git.IsOffline() {
		script.Fetch()
	}
	git.EnsureDoesNotHaveBranch(result.TargetBranch)
	return
}

func getHackStepList(config hackConfig) (result steps.StepList) {
	mainBranchName := git.GetMainBranch()
	result.AppendList(steps.GetSyncBranchSteps(mainBranchName, true))
	result.Append(&steps.CreateAndCheckoutBranchStep{BranchName: config.TargetBranch, ParentBranchName: mainBranchName})
	if git.HasRemote("origin") && git.ShouldNewBranchPush() && !git.IsOffline() {
		result.Append(&steps.CreateTrackingBranchStep{BranchName: config.TargetBranch})
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
	return
}

func init() {
	RootCmd.AddCommand(hackCmd)
}
