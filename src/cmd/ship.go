package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/git-town/git-town/src/script"
	"github.com/git-town/git-town/src/steps"
	"github.com/git-town/git-town/src/util"

	"github.com/spf13/cobra"
)

type shipConfig struct {
	BranchToShip  string
	InitialBranch string // the name of the branch that was checked out when running this command
}

// optional commit message provided via the command line
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

Ships only direct children of the main branch.
To ship a nested child branch, all ancestor branches must be shipped or killed first.

If you are using GitHub, this command can squash merge pull requests via the GitHub API. Setup:
1. Get a GitHub personal access token with the "repo" scope
2. Run 'git config git-town.github-token XXX' (optionally add the '--global' flag)
Now anytime you ship a branch with a pull request on GitHub, it will squash merge via the GitHub API.
It will also update the base branch for any pull requests against that branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := gitShipConfig(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		stepList, err := getShipStepList(config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		runState := steps.NewRunState("ship", stepList)
		err = steps.Run(runState)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
		)
	},
}

func gitShipConfig(args []string) (result shipConfig, err error) {
	result.InitialBranch = git.GetCurrentBranchName()
	if len(args) == 0 {
		result.BranchToShip = result.InitialBranch
	} else {
		result.BranchToShip = args[0]
	}
	if result.BranchToShip == result.InitialBranch {
		git.EnsureDoesNotHaveOpenChanges("Did you mean to commit them before shipping?")
	}
	if git.HasRemote("origin") && !git.Config().IsOffline() {
		err := script.Fetch()
		if err != nil {
			return result, err
		}
	}
	if result.BranchToShip != result.InitialBranch {
		git.EnsureHasBranch(result.BranchToShip)
	}
	git.Config().EnsureIsFeatureBranch(result.BranchToShip, "Only feature branches can be shipped.")
	prompt.EnsureKnowsParentBranches([]string{result.BranchToShip})
	ensureParentBranchIsMainOrPerennialBranch(result.BranchToShip)
	return
}

func ensureParentBranchIsMainOrPerennialBranch(branchName string) {
	parentBranch := git.Config().GetParentBranch(branchName)
	if !git.Config().IsMainBranch(parentBranch) && !git.Config().IsPerennialBranch(parentBranch) {
		ancestors := git.Config().GetAncestorBranches(branchName)
		ancestorsWithoutMainOrPerennial := ancestors[1:]
		oldestAncestor := ancestorsWithoutMainOrPerennial[0]
		util.ExitWithErrorMessage(
			"Shipping this branch would ship "+strings.Join(ancestorsWithoutMainOrPerennial, ", ")+" as well.",
			"Please ship \""+oldestAncestor+"\" first.",
		)
	}
}

func getShipStepList(config shipConfig) (steps.StepList, error) {
	result := steps.StepList{}
	isOffline := git.Config().IsOffline()
	branchToMergeInto := git.Config().GetParentBranch(config.BranchToShip)
	isShippingInitialBranch := config.BranchToShip == config.InitialBranch
	result.AppendList(steps.GetSyncBranchSteps(branchToMergeInto, true))
	result.AppendList(steps.GetSyncBranchSteps(config.BranchToShip, false))
	result.Append(&steps.EnsureHasShippableChangesStep{BranchName: config.BranchToShip})
	result.Append(&steps.CheckoutBranchStep{BranchName: branchToMergeInto})
	canShipWithDriver, defaultCommitMessage, err := getCanShipWithDriver(config.BranchToShip, branchToMergeInto)
	if err != nil {
		return result, err
	}
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
	childBranches := git.Config().GetChildBranches(config.BranchToShip)
	// NOTE: when shipping with a driver, we can always delete the remote branch because:
	// - we know we have a tracking branch (otherwise there would be no PR to ship via driver)
	// - we have updated the PRs of all child branches (because we have API access)
	// - we know we are online
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
	return result, nil
}

func getCanShipWithDriver(branch, parentBranch string) (bool, string, error) {
	if !git.HasRemote("origin") {
		return false, "", nil
	}
	if git.Config().IsOffline() {
		return false, "", nil
	}
	driver := drivers.GetActiveDriver()
	if driver == nil {
		return false, "", nil
	}
	return driver.CanMergePullRequest(branch, parentBranch)
}

func init() {
	shipCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Specify the commit message for the squash commit")
	RootCmd.AddCommand(shipCmd)
}
