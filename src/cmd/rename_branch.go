package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/script"
	"github.com/git-town/git-town/src/steps"
	"github.com/git-town/git-town/src/util"
	"github.com/spf13/cobra"
)

type renameBranchConfig struct {
	OldBranchName string
	NewBranchName string
}

var forceFlag bool

var renameBranchCommand = &cobra.Command{
	Use:   "rename-branch [<old_branch_name>] <new_branch_name>",
	Short: "Renames a branch both locally and remotely",
	Long: `Renames a branch both locally and remotely

Renames the given branch in the local and origin repository.
Aborts if the new branch name already exists or the tracking branch is out of sync.

- creates a branch with the new name
- deletes the old branch

When there is a remote repository
- syncs the repository

When there is a tracking branch
- pushes the new branch to the remote repository
- deletes the old branch from the remote repository

When run on a perennial branch
- confirm with the "-f" option
- registers the new perennial branch name in the local Git Town configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := getRenameBranchConfig(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		stepList := getRenameBranchStepList(config)
		runState := steps.NewRunState("rename-branch", stepList)
		err = steps.Run(runState)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: cobra.RangeArgs(1, 2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
		)
	},
}

func getRenameBranchConfig(args []string) (result renameBranchConfig, err error) {
	if len(args) == 1 {
		result.OldBranchName = git.GetCurrentBranchName()
		result.NewBranchName = args[0]
	} else {
		result.OldBranchName = args[0]
		result.NewBranchName = args[1]
	}
	git.EnsureIsNotMainBranch(result.OldBranchName, "The main branch cannot be renamed.")
	if !forceFlag {
		git.EnsureIsNotPerennialBranch(result.OldBranchName, fmt.Sprintf("%q is a perennial branch. Renaming a perennial branch typically requires other updates. If you are sure you want to do this, use '--force'.", result.OldBranchName))
	}
	if result.OldBranchName == result.NewBranchName {
		util.ExitWithErrorMessage("Cannot rename branch to current name.")
	}
	if !git.Config().IsOffline() {
		err := script.Fetch()
		if err != nil {
			return result, err
		}
	}
	git.EnsureHasBranch(result.OldBranchName)
	git.EnsureBranchInSync(result.OldBranchName, "Please sync the branches before renaming.")
	git.EnsureDoesNotHaveBranch(result.NewBranchName)
	return
}

func getRenameBranchStepList(config renameBranchConfig) (result steps.StepList) {
	result.Append(&steps.CreateBranchStep{BranchName: config.NewBranchName, StartingPoint: config.OldBranchName})
	if git.GetCurrentBranchName() == config.OldBranchName {
		result.Append(&steps.CheckoutBranchStep{BranchName: config.NewBranchName})
	}
	if git.Config().IsPerennialBranch(config.OldBranchName) {
		result.Append(&steps.RemoveFromPerennialBranches{BranchName: config.OldBranchName})
		result.Append(&steps.AddToPerennialBranches{BranchName: config.NewBranchName})
	} else {
		result.Append(&steps.DeleteParentBranchStep{BranchName: config.OldBranchName})
		result.Append(&steps.SetParentBranchStep{BranchName: config.NewBranchName, ParentBranchName: git.Config().GetParentBranch(config.OldBranchName)})
	}
	for _, child := range git.Config().GetChildBranches(config.OldBranchName) {
		result.Append(&steps.SetParentBranchStep{BranchName: child, ParentBranchName: config.NewBranchName})
	}
	if git.HasTrackingBranch(config.OldBranchName) && !git.Config().IsOffline() {
		result.Append(&steps.CreateTrackingBranchStep{BranchName: config.NewBranchName})
		result.Append(&steps.DeleteRemoteBranchStep{BranchName: config.OldBranchName, IsTracking: true})
	}
	result.Append(&steps.DeleteLocalBranchStep{BranchName: config.OldBranchName})
	result.Wrap(steps.WrapOptions{RunInGitRoot: false, StashOpenChanges: false})
	return
}

func init() {
	renameBranchCommand.Flags().BoolVar(&forceFlag, "force", false, "Force rename of perennial branch")
	RootCmd.AddCommand(renameBranchCommand)
}
