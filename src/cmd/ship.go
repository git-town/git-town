package cmd

import (
	"strings"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/drivers"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/prompt"
	"github.com/Originate/git-town/src/script"
	"github.com/Originate/git-town/src/steps"
	"github.com/Originate/git-town/src/util"

	"github.com/spf13/cobra"
)

type shipConfig struct {
	BranchToShip  string
	InitialBranch string
}

var commitMessage string

var shipCmd = &cobra.Command{
	Use:   "ship",
	Short: "Deliver a completed feature branch",
	Long: `Deliver a completed feature branch

Squash-merges the current branch, or <branch_name> if given,
into the main branch, resulting in linear history on the main branch.

- syncs the main branch
- pulls remote updates for <branch_name>
- merges the main branch into <branch_name>
- squash-merges <branch_name> into the main branch
  with commit message specified by the user
- pushes the main branch to the remote repository
- deletes <branch_name> from the local and remote repositories

Only shipping of direct children of the main branch is allowed.
To ship a nested child branch, all ancestor branches have to be shipped or killed.

If you are using GitHub, this command can squash merge pull requests via the GitHub API. Setup:
1. Get a GitHub personal access token with the "repo" scope
2. Run 'git config git-town.github-token XXX' (optionally add the '--global' flag)
Now anytime you ship a branch with a pull request on GitHub, it will squash merge via the GitHub API.
It will also update the base branch for any pull requests against that branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := gitShipConfig(args)
		stepList := getShipStepList(config)
		runState := steps.NewRunState("ship", stepList)
		steps.Run(runState)
	},
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
		)
	},
}

func gitShipConfig(args []string) (result shipConfig) {
	result.InitialBranch = git.GetCurrentBranchName()
	if len(args) == 0 {
		result.BranchToShip = result.InitialBranch
	} else {
		result.BranchToShip = args[0]
	}
	if result.BranchToShip == result.InitialBranch {
		git.EnsureDoesNotHaveOpenChanges("Did you mean to commit them before shipping?")
	}
	if git.HasRemote("origin") && !git.IsOffline() {
		script.Fetch()
	}
	if result.BranchToShip != result.InitialBranch {
		git.EnsureHasBranch(result.BranchToShip)
	}
	git.EnsureIsFeatureBranch(result.BranchToShip, "Only feature branches can be shipped.")
	prompt.EnsureKnowsParentBranches([]string{result.BranchToShip})
	ensureParentBranchIsMainOrPerennialBranch(result.BranchToShip)
	return
}

func ensureParentBranchIsMainOrPerennialBranch(branchName string) {
	parentBranch := git.GetParentBranch(branchName)
	if !git.IsMainBranch(parentBranch) && !git.IsPerennialBranch(parentBranch) {
		ancestors := git.GetAncestorBranches(branchName)
		ancestorsWithoutMainOrPerennial := ancestors[1:]
		oldestAncestor := ancestorsWithoutMainOrPerennial[0]
		util.ExitWithErrorMessage(
			"Shipping this branch would ship "+strings.Join(ancestorsWithoutMainOrPerennial, ", ")+" as well.",
			"Please ship \""+oldestAncestor+"\" first.",
		)
	}
}

func getShipStepList(config shipConfig) (result steps.StepList) {
	var isOffline = git.IsOffline()
	branchToMergeInto := git.GetParentBranch(config.BranchToShip)
	isShippingInitialBranch := config.BranchToShip == config.InitialBranch
	result.AppendList(steps.GetSyncBranchSteps(branchToMergeInto, true))
	result.AppendList(steps.GetSyncBranchSteps(config.BranchToShip, false))
	result.Append(&steps.EnsureHasShippableChangesStep{BranchName: config.BranchToShip})
	result.Append(&steps.CheckoutBranchStep{BranchName: branchToMergeInto})
	canShipWithDriver, defaultCommitMessage := getCanShipWithDriver(config.BranchToShip, branchToMergeInto)
	if canShipWithDriver {
		result.Append(&steps.PushBranchStep{BranchName: config.BranchToShip})
		result.Append(&steps.DriverMergePullRequestStep{BranchName: config.BranchToShip, CommitMessage: commitMessage, DefaultCommitMessage: defaultCommitMessage})
		result.Append(&steps.PullBranchStep{})
	} else {
		result.Append(&steps.SquashMergeBranchStep{BranchName: config.BranchToShip, CommitMessage: commitMessage})
	}
	if git.HasRemote("origin") && !isOffline {
		result.Append(&steps.PushBranchStep{BranchName: branchToMergeInto, Undoable: true})
	}
	childBranches := git.GetChildBranches(config.BranchToShip)
	if canShipWithDriver || (git.HasTrackingBranch(config.BranchToShip) && len(childBranches) == 0 && !isOffline) {
		result.Append(&steps.DeleteRemoteBranchStep{BranchName: config.BranchToShip, IsTracking: true})
	}
	result.Append(&steps.DeleteLocalBranchStep{BranchName: config.BranchToShip})
	result.Append(&steps.DeleteParentBranchStep{BranchName: config.BranchToShip})
	for _, child := range childBranches {
		result.Append(&steps.SetParentBranchStep{BranchName: child, ParentBranchName: branchToMergeInto})
	}
	if !isShippingInitialBranch {
		result.Append(&steps.CheckoutBranchStep{BranchName: config.InitialBranch})
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: !isShippingInitialBranch})
	return
}

func getCanShipWithDriver(branch, parentBranch string) (bool, string) {
	if !git.HasRemote("origin") {
		return false, ""
	}
	if git.IsOffline() {
		return false, ""
	}
	driver := drivers.GetActiveDriver()
	if driver == nil {
		return false, ""
	}
	canMerge, defaultCommitMessage, err := driver.CanMergePullRequest(branch, parentBranch)
	exit.If(err)
	return canMerge, defaultCommitMessage
}

func init() {
	shipCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Specify the commit message for the squash commit")
	RootCmd.AddCommand(shipCmd)
}
