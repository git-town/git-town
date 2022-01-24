package cmd

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/git-town/git-town/v7/src/userinput"
	"github.com/spf13/cobra"
)

type shipConfig struct {
	branchToShip                 string
	branchToMergeInto            string
	canShipWithDriver            bool
	childBranches                []string
	defaultCommitMessage         string
	hasOrigin                    bool
	hasTrackingBranch            bool
	initialBranch                string
	isShippingInitialBranch      bool
	isOffline                    bool
	pullRequestNumber            int64
	shouldShipDeleteRemoteBranch bool
}

// optional commit message provided via the command line.
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
		driver := hosting.NewDriver(&prodRepo.Config, &prodRepo.Silent, cli.PrintDriverAction)
		config, err := gitShipConfig(args, driver, prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		stepList, err := createShipStepList(config, prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		runState := runstate.New("ship", stepList)
		err = runstate.Execute(runState, prodRepo, driver)
		if err != nil {
			cli.Exit(err)
		}
	},
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := ValidateIsRepository(prodRepo); err != nil {
			return err
		}
		return validateIsConfigured(prodRepo)
	},
}

//nolint:funlen
func gitShipConfig(args []string, driver hosting.Driver, repo *git.ProdRepo) (result shipConfig, err error) {
	result.initialBranch, err = repo.Silent.CurrentBranch()
	if err != nil {
		return result, err
	}
	if len(args) == 0 {
		result.branchToShip = result.initialBranch
	} else {
		result.branchToShip = args[0]
	}
	if result.branchToShip == result.initialBranch {
		hasOpenChanges, err := repo.Silent.HasOpenChanges()
		if err != nil {
			return result, err
		}
		if hasOpenChanges {
			return result, fmt.Errorf("you have uncommitted changes. Did you mean to commit them before shipping?")
		}
	}
	result.hasOrigin, err = repo.Silent.HasRemote("origin")
	if err != nil {
		return result, err
	}
	if result.hasOrigin && !repo.Config.IsOffline() {
		err := repo.Logging.Fetch()
		if err != nil {
			return result, err
		}
	}
	if result.branchToShip != result.initialBranch {
		hasBranch, err := repo.Silent.HasLocalOrRemoteBranch(result.branchToShip)
		if err != nil {
			return result, err
		}
		if !hasBranch {
			return result, fmt.Errorf("there is no branch named %q", result.branchToShip)
		}
	}
	if !repo.Config.IsFeatureBranch(result.branchToShip) {
		return result, fmt.Errorf("the branch %q is not a feature branch. Only feature branches can be shipped", result.branchToShip)
	}
	err = userinput.EnsureKnowsParentBranches([]string{result.branchToShip}, repo)
	if err != nil {
		return result, err
	}
	ensureParentBranchIsMainOrPerennialBranch(result.branchToShip)
	result.hasTrackingBranch, err = repo.Silent.HasTrackingBranch(result.branchToShip)
	if err != nil {
		return result, err
	}
	result.isOffline = repo.Config.IsOffline()
	result.isShippingInitialBranch = result.branchToShip == result.initialBranch
	result.branchToMergeInto = repo.Config.ParentBranch(result.branchToShip)
	prInfo, err := createPullRequestInfo(result.branchToShip, result.branchToMergeInto, driver)
	result.canShipWithDriver = prInfo.CanMergeWithAPI
	result.defaultCommitMessage = prInfo.DefaultCommitMessage
	result.pullRequestNumber = prInfo.PullRequestNumber
	result.childBranches = repo.Config.ChildBranches(result.branchToShip)
	result.shouldShipDeleteRemoteBranch = prodRepo.Config.ShouldShipDeleteRemoteBranch()
	return result, err
}

func ensureParentBranchIsMainOrPerennialBranch(branchName string) {
	parentBranch := prodRepo.Config.ParentBranch(branchName)
	if !prodRepo.Config.IsMainBranch(parentBranch) && !prodRepo.Config.IsPerennialBranch(parentBranch) {
		ancestors := prodRepo.Config.AncestorBranches(branchName)
		ancestorsWithoutMainOrPerennial := ancestors[1:]
		oldestAncestor := ancestorsWithoutMainOrPerennial[0]
		cli.Exit(fmt.Errorf(`shipping this branch would ship %q as well,
please ship %q first`, strings.Join(ancestorsWithoutMainOrPerennial, ", "), oldestAncestor))
	}
}

func createShipStepList(config shipConfig, repo *git.ProdRepo) (result runstate.StepList, err error) {
	syncSteps, err := runstate.SyncBranchSteps(config.branchToMergeInto, true, repo)
	if err != nil {
		return result, err
	}
	result.AppendList(syncSteps)
	syncSteps, err = runstate.SyncBranchSteps(config.branchToShip, false, repo)
	if err != nil {
		return result, err
	}
	result.AppendList(syncSteps)
	result.Append(&steps.EnsureHasShippableChangesStep{BranchName: config.branchToShip})
	result.Append(&steps.CheckoutBranchStep{BranchName: config.branchToMergeInto})
	if config.canShipWithDriver {
		result.Append(&steps.PushBranchStep{BranchName: config.branchToShip})
		result.Append(&steps.DriverMergePullRequestStep{
			BranchName:           config.branchToShip,
			PullRequestNumber:    config.pullRequestNumber,
			CommitMessage:        commitMessage,
			DefaultCommitMessage: config.defaultCommitMessage,
		})
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
	err = result.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: !config.isShippingInitialBranch}, repo)
	return result, err
}

func createPullRequestInfo(branch, parentBranch string, driver hosting.Driver) (result hosting.PullRequestInfo, err error) {
	hasOrigin, err := prodRepo.Silent.HasRemote("origin")
	if err != nil {
		return result, err
	}
	if !hasOrigin {
		return result, nil
	}
	if prodRepo.Config.IsOffline() {
		return result, nil
	}
	if driver == nil {
		return result, nil
	}
	return driver.LoadPullRequestInfo(branch, parentBranch)
}

func init() {
	shipCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Specify the commit message for the squash commit")
	RootCmd.AddCommand(shipCmd)
}
