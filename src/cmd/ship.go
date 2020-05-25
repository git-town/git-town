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
	BranchToShip            string
	InitialBranch           string // the name of the branch that was checked out when running this command
	branchToMergeInto       string
	defaultCommitMessage    string
	childBranches           []string
	isOffline               bool
	isShippingInitialBranch bool
	canShipWithDriver       bool
	hasOrigin               bool
	pullRequestNumber       int
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

Ships direct children of the main branch.
To ship a nested child branch, ship or kill all ancestor branches first.

If you use GitHub, this command can squash merge pull requests via the GitHub API. Setup:
1. Get a GitHub personal access token with the "repo" scope
2. Run 'git config git-town.github-token XXX' (optionally add the '--global' flag)
Now anytime you ship a branch with a pull request on GitHub, it will squash merge via the GitHub API.
It will also update the base branch for any pull requests against that branch.

If your origin server deletes shipped branches, for example
GitHub's feature to automatically delete head branches,
run "git config git-town.ship-delete-remote-branch false"
and Git Town will leave it up to your origin server to delete the remote branch.`,
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
	result.hasOrigin = git.HasRemote("origin")
	result.isOffline = git.Config().IsOffline()
	result.isShippingInitialBranch = result.BranchToShip == result.InitialBranch
	result.branchToMergeInto = git.Config().GetParentBranch(result.BranchToShip)
	result.canShipWithDriver, result.defaultCommitMessage, result.pullRequestNumber, err = getCanShipWithDriver(result.BranchToShip, result.branchToMergeInto)
	result.childBranches = git.Config().GetChildBranches(result.BranchToShip)
	return result, err
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
	result.AppendList(steps.GetSyncBranchSteps(config.branchToMergeInto, true))
	result.AppendList(steps.GetSyncBranchSteps(config.BranchToShip, false))
	result.Append(&steps.EnsureHasShippableChangesStep{BranchName: config.BranchToShip})
	result.Append(&steps.CheckoutBranchStep{BranchName: config.branchToMergeInto})
	if config.canShipWithDriver {
		result.Append(&steps.PushBranchStep{BranchName: config.BranchToShip})
		result.Append(&steps.DriverMergePullRequestStep{BranchName: config.BranchToShip, PullRequestNumber: config.pullRequestNumber, CommitMessage: commitMessage, DefaultCommitMessage: config.defaultCommitMessage})
		result.Append(&steps.PullBranchStep{})
	} else {
		result.Append(&steps.SquashMergeBranchStep{BranchName: config.BranchToShip, CommitMessage: commitMessage})
	}
	if config.hasOrigin && !config.isOffline {
		result.Append(&steps.PushBranchStep{BranchName: config.branchToMergeInto, Undoable: true})
	}
	// NOTE: when shipping with a driver, we can always delete the remote branch because:
	// - we know we have a tracking branch (otherwise there would be no PR to ship via driver)
	// - we have updated the PRs of all child branches (because we have API access)
	// - we know we are online
	if config.canShipWithDriver || (git.HasTrackingBranch(config.BranchToShip) && len(config.childBranches) == 0 && !config.isOffline) {
		if git.Config().ShouldShipDeleteRemoteBranch() {
			result.Append(&steps.DeleteRemoteBranchStep{BranchName: config.BranchToShip, IsTracking: true})
		}
	}
	result.Append(&steps.DeleteLocalBranchStep{BranchName: config.BranchToShip})
	result.Append(&steps.DeleteParentBranchStep{BranchName: config.BranchToShip})
	for _, child := range config.childBranches {
		result.Append(&steps.SetParentBranchStep{BranchName: child, ParentBranchName: config.branchToMergeInto})
	}
	if !config.isShippingInitialBranch {
		result.Append(&steps.CheckoutBranchStep{BranchName: config.InitialBranch})
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: !config.isShippingInitialBranch})
	return result, nil
}

func getCanShipWithDriver(branch, parentBranch string) (canShip bool, defaultCommitMessage string, pullRequestNumber int, err error) {
	if !git.HasRemote("origin") {
		return false, "", 0, nil
	}
	if git.Config().IsOffline() {
		return false, "", 0, nil
	}
	driver := drivers.GetActiveDriver()
	if driver == nil {
		return false, "", 0, nil
	}
	return driver.CanMergePullRequest(branch, parentBranch)
}

func init() {
	shipCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Specify the commit message for the squash commit")
	RootCmd.AddCommand(shipCmd)
}
