package cmd

import (
	"fmt"
	"log"

	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/gitconfig"
	"github.com/Originate/git-town/lib/prompt"
	"github.com/Originate/git-town/lib/steps"

	"github.com/spf13/cobra"
)

type SyncConfig struct {
	InitialBranch  string
	BranchesToSync []string
	ShouldPushTags bool
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Update the current branch with all relevant changes",
	Long:  `Update the current branch with all relevant changes`,
	Run: func(cmd *cobra.Command, args []string) {
		steps.Run(steps.RunOptions{
			CanSkip: func() bool {
				return !(git.IsRebaseInProgress() && gitconfig.IsMainBranch(git.GetCurrentBranchName()))
			},
			Command:    "sync",
			IsAbort:    abortFlag,
			IsContinue: continueFlag,
			IsSkip:     skipFlag,
			IsUndo:     false,
			SkipMessageGenerator: func() string {
				return fmt.Sprintf("the sync of the '%s' branch", git.GetCurrentBranchName())
			},
			StepListGenerator: func() steps.StepList {
				syncConfig := checkSyncPreconditions()
				return getSyncStepList(syncConfig)
			},
		})
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateMaxArgs(args, 0)
	},
}

func checkSyncPreconditions() (result SyncConfig) {
	if gitconfig.HasRemote("origin") {
		err := steps.FetchStep{}.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
	result.InitialBranch = git.GetCurrentBranchName()
	if allFlag {
		branches := git.GetLocalBranchesWithMainBranchFirst()
		prompt.EnsureKnowsParentBranches(branches)
		result.BranchesToSync = branches
		result.ShouldPushTags = true
	} else if gitconfig.IsFeatureBranch(result.InitialBranch) {
		prompt.EnsureKnowsParentBranches([]string{result.InitialBranch})
		result.BranchesToSync = append(gitconfig.GetAncestorBranches(result.InitialBranch), result.InitialBranch)
	} else {
		result.BranchesToSync = []string{result.InitialBranch}
		result.ShouldPushTags = true
	}
	return
}

func getSyncStepList(syncConfig SyncConfig) steps.StepList {
	stepList := steps.StepList{}
	for _, branchName := range syncConfig.BranchesToSync {
		stepList.AppendList(steps.GetSyncBranchSteps(branchName))
	}
	stepList.Append(steps.CheckoutBranchStep{BranchName: syncConfig.InitialBranch})
	if gitconfig.HasRemote("origin") && syncConfig.ShouldPushTags {
		stepList.Append(steps.PushTagsStep{})
	}
	return steps.Wrap(stepList, steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
}

func init() {
	syncCmd.Flags().BoolVar(&allFlag, "all", false, "Sync all local branches")
	syncCmd.Flags().BoolVar(&abortFlag, "abort", false, "Abort a previous command that resulted in a conflict")
	syncCmd.Flags().BoolVar(&continueFlag, "continue", false, "Continue a previous command that resulted in a conflict")
	syncCmd.Flags().BoolVar(&skipFlag, "skip", false, "Continue a previous command by skipping the branch that resulted in a conflicted")
	RootCmd.AddCommand(syncCmd)
}
