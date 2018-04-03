package cmd

import (
	"fmt"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/script"
	"github.com/Originate/git-town/src/steps"
	"github.com/Originate/git-town/src/util"
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

Renames the given branch on both the local machine and the remote if one is configured.
Aborts if the new branch name already exists or the tracking branch is out of sync.
This command is intended for feature branches.
Renaming perennial branches has to be confirmed with the "-f" option.

- Creates a branch with the new name
- Deletes the old branch

When there is a remote repository
- Syncs the repository

When there is a tracking branch
- Pushes the new branch to the remote repository
- Deletes the old branch from the remote repository

When run on a perennial branch
- Requires the use of the "-f" option
- Reconfigures git town locally for the perennial branch`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getRenameBranchConfig(args)
		stepList := getRenameBranchStepList(config)
		runState := steps.NewRunState("rename-branch", stepList)
		steps.Run(runState)
	},
	Args: cobra.RangeArgs(1, 2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
		)
	},
}

func getRenameBranchConfig(args []string) (result renameBranchConfig) {
	if len(args) == 1 {
		result.OldBranchName = git.GetCurrentBranchName()
		result.NewBranchName = args[0]
	} else {
		result.OldBranchName = args[0]
		result.NewBranchName = args[1]
	}
	git.EnsureIsNotMainBranch(result.OldBranchName, "The main branch cannot be renamed.")
	if !forceFlag {
		git.EnsureIsNotPerennialBranch(result.OldBranchName, fmt.Sprintf("'%s' is a perennial branch. Renaming a perennial branch typically requires other updates. If you are sure you want to do this, use '--force'.", result.OldBranchName))
	}
	if result.OldBranchName == result.NewBranchName {
		util.ExitWithErrorMessage("Cannot rename branch to current name.")
	}
	if !git.IsOffline() {
		script.Fetch()
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
	if git.IsPerennialBranch(config.OldBranchName) {
		result.Append(&steps.RemoveFromPerennialBranches{BranchName: config.OldBranchName})
		result.Append(&steps.AddToPerennialBranches{BranchName: config.NewBranchName})
	} else {
		result.Append(&steps.DeleteParentBranchStep{BranchName: config.OldBranchName})
		result.Append(&steps.SetParentBranchStep{BranchName: config.NewBranchName, ParentBranchName: git.GetParentBranch(config.OldBranchName)})
	}
	for _, child := range git.GetChildBranches(config.OldBranchName) {
		result.Append(&steps.SetParentBranchStep{BranchName: child, ParentBranchName: config.NewBranchName})
	}
	if git.HasTrackingBranch(config.OldBranchName) && !git.IsOffline() {
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
