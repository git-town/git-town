package cmd

import (
	"errors"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/prompt"
	"github.com/Originate/git-town/src/script"
	"github.com/Originate/git-town/src/steps"

	"github.com/spf13/cobra"
)

type appendConfig struct {
	InitialBranch string
	TargetBranch  string
}

var appendCommand = &cobra.Command{
	Use:   "append <branch>",
	Short: "Creates a new feature branch as a child of the current branch",
	Long: `Creates a new feature branch as a child of the current branch.

Syncs the current branch,
forks a new feature branch with the given name off the current branch,
makes the new branch a child of the current branch,
pushes the new feature branch to the remote repository,
and brings over all uncommitted changes to the new feature branch.

Additionally, when there is a remote upstream,
the main branch is synced with its upstream counterpart.
This can be disabled by toggling the "hack-push-flag" configuration:
$ git town hack-push-flag false`,
	Run: func(cmd *cobra.Command, args []string) {
		git.EnsureIsRepository()
		prompt.EnsureIsConfigured()
		steps.Run(steps.RunOptions{
			CanSkip:              func() bool { return false },
			Command:              "append",
			IsAbort:              abortFlag,
			IsContinue:           continueFlag,
			IsSkip:               false,
			IsUndo:               undoFlag,
			SkipMessageGenerator: func() string { return "" },
			StepListGenerator: func() steps.StepList {
				config := checkAppendPreconditions(args)
				return getAppendStepList(config)
			},
		})
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && !abortFlag && !continueFlag && !undoFlag {
			return errors.New("no branch name provided")
		}
		return validateMaxArgs(args, 1)
	},
}

func checkAppendPreconditions(args []string) (result appendConfig) {
	result.InitialBranch = git.GetCurrentBranchName()
	result.TargetBranch = args[0]
	if git.HasRemote("origin") {
		script.Fetch()
	}
	git.EnsureDoesNotHaveBranch(result.TargetBranch)
	prompt.EnsureKnowsParentBranches([]string{result.InitialBranch})
	return
}

func getAppendStepList(config appendConfig) (result steps.StepList) {
	for _, branchName := range append(git.GetAncestorBranches(config.InitialBranch), config.InitialBranch) {
		result.AppendList(steps.GetSyncBranchSteps(branchName))
	}
	result.Append(steps.CreateBranchStep{BranchName: config.TargetBranch, StartingPoint: config.InitialBranch})
	result.Append(steps.SetParentBranchStep{BranchName: config.TargetBranch, ParentBranchName: config.InitialBranch})
	result.Append(steps.CheckoutBranchStep{BranchName: config.TargetBranch})
	if git.HasRemote("origin") && git.ShouldHackPush() {
		result.Append(steps.CreateTrackingBranchStep{BranchName: config.TargetBranch})
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
	return
}

func init() {
	appendCommand.Flags().BoolVar(&abortFlag, "abort", false, abortFlagDescription)
	appendCommand.Flags().BoolVar(&continueFlag, "continue", false, continueFlagDescription)
	appendCommand.Flags().BoolVar(&undoFlag, "undo", false, undoFlagDescription)
	RootCmd.AddCommand(appendCommand)
}
