package cmd

import (
	"errors"

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

Syncs the main branch,
forks a new feature branch with the given name off it,
pushes the new feature branch to the remote repository,
and brings over all uncommitted changes to the new feature branch.

Additionally, when there is a remote upstream,
the main branch is synced with its upstream counterpart.
This can be disabled by toggling the "new-branch-push-flag" configuration:
$ git town new-branch-push-flag false`,
	Run: func(cmd *cobra.Command, args []string) {
		steps.Run(steps.RunOptions{
			CanSkip:              func() bool { return false },
			Command:              "hack",
			IsAbort:              abortFlag,
			IsContinue:           continueFlag,
			IsSkip:               false,
			IsUndo:               false,
			SkipMessageGenerator: func() string { return "" },
			StepListGenerator: func() steps.StepList {
				config := getHackConfig(args)
				return getHackStepList(config)
			},
		})
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && !abortFlag && !continueFlag {
			return errors.New("no branch name provided")
		}
		return util.FirstError(
			validateMaxArgsFunc(args, 1),
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
	result.AppendList(steps.GetSyncBranchSteps(mainBranchName))
	result.Append(&steps.CreateAndCheckoutBranchStep{BranchName: config.TargetBranch, ParentBranchName: mainBranchName})
	if git.HasRemote("origin") && git.ShouldNewBranchPush() && !git.IsOffline() {
		result.Append(&steps.CreateTrackingBranchStep{BranchName: config.TargetBranch})
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
	return
}

func init() {
	hackCmd.Flags().BoolVar(&abortFlag, "abort", false, abortFlagDescription)
	hackCmd.Flags().BoolVar(&continueFlag, "continue", false, continueFlagDescription)
	RootCmd.AddCommand(hackCmd)
}
