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
	pullRequestNumber            int64
	branchToShip                 string
	branchToMergeInto            string
	initialBranch                string
	defaultCommitMessage         string
	canShipWithDriver            bool
	hasOrigin                    bool
	hasTrackingBranch            bool
	isOffline                    bool
	isShippingInitialBranch      bool
	shouldShipDeleteRemoteBranch bool
	childBranches                []string
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
		stepList := getShipStepList(config)
		runState := steps.NewRunState("ship", stepList)
		err = steps.Run(runState)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := git.ValidateIsRepository(); err != nil {
			return err
		}
		return validateIsConfigured()
	},
}

func gitShipConfig(args []string) (result shipConfig, err error) {
	result.initialBranch = git.GetCurrentBranchName()
	if len(args) == 0 {
		result.branchToShip = result.initialBranch
	} else {
		result.branchToShip = args[0]
	}
	if result.branchToShip == result.initialBranch {
		git.EnsureDoesNotHaveOpenChanges("Did you mean to commit them before shipping?")
	}
	if git.HasRemote("origin") && !git.Config().IsOffline() {
		err := script.Fetch()
		if err != nil {
			return result, err
		}
	}
	if result.branchToShip != result.initialBranch {
		git.EnsureHasBranch(result.branchToShip)
	}
	git.Config().EnsureIsFeatureBranch(result.branchToShip, "Only feature branches can be shipped.")
	prompt.EnsureKnowsParentBranches([]string{result.branchToShip})
	ensureParentBranchIsMainOrPerennialBranch(result.branchToShip)
	result.hasTrackingBranch = git.HasTrackingBranch(result.branchToShip)
	result.hasOrigin = git.HasRemote("origin")
	result.isOffline = git.Config().IsOffline()
	result.isShippingInitialBranch = result.branchToShip == result.initialBranch
	result.branchToMergeInto = git.Config().GetParentBranch(result.branchToShip)
	result.canShipWithDriver, result.defaultCommitMessage, result.pullRequestNumber, err = getCanShipWithDriver(result.branchToShip, result.branchToMergeInto)
	result.childBranches = git.Config().GetChildBranches(result.branchToShip)
	result.shouldShipDeleteRemoteBranch = git.Config().ShouldShipDeleteRemoteBranch()
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

func getShipStepList(config shipConfig) (result steps.StepList) {
	result.AppendList(steps.GetSyncBranchSteps(config.branchToMergeInto, true))
	result.AppendList(steps.GetSyncBranchSteps(config.branchToShip, false))
	result.Append(&steps.EnsureHasShippableChangesStep{BranchName: config.branchToShip})
	result.Append(&steps.CheckoutBranchStep{BranchName: config.branchToMergeInto})
	if config.canShipWithDriver {
		result.Append(&steps.PushBranchStep{BranchName: config.branchToShip})
		result.Append(&steps.DriverMergePullRequestStep{BranchName: config.branchToShip, PullRequestNumber: config.pullRequestNumber, CommitMessage: commitMessage, DefaultCommitMessage: config.defaultCommitMessage})
		result.Append(&steps.PullBranchStep{})
	} else {
		result.Append(&steps.SquashMergeBranchStep{BranchName: config.branchToShip, CommitMessage: commitMessage})
	}
	if config.hasOrigin && !config.isOffline {
		result.Append(&steps.PushBranchStep{BranchName: config.branchToMergeInto, Undoable: true})
	}
	// NOTE: when shipping with a driver, we can always delete the remote branch because:
	// - we know we have a tracking branch (otherwise there would be no PR to ship via driver)
	// - we have updated the PRs of all child branches (because we have API access)
	// - we know we are online
	if config.canShipWithDriver || (config.hasTrackingBranch && len(config.childBranches) == 0 && !config.isOffline) {
		if config.shouldShipDeleteRemoteBranch {
			result.Append(&steps.DeleteRemoteBranchStep{BranchName: config.branchToShip, IsTracking: true})
		}
	}
	result.Append(&steps.DeleteLocalBranchStep{BranchName: config.branchToShip})
	result.Append(&steps.DeleteParentBranchStep{BranchName: config.branchToShip})
	for _, child := range config.childBranches {
		result.Append(&steps.SetParentBranchStep{BranchName: child, ParentBranchName: config.branchToMergeInto})
	}
	if !config.isShippingInitialBranch {
		result.Append(&steps.CheckoutBranchStep{BranchName: config.initialBranch})
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: !config.isShippingInitialBranch})
	return result
}

func getCanShipWithDriver(branch, parentBranch string) (canShip bool, defaultCommitMessage string, pullRequestNumber int64, err error) {
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
