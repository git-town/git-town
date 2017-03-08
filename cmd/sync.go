package cmd

import (
  "log"

  "github.com/Originate/gt/cmd/config"
  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/script"
  "github.com/Originate/gt/cmd/steps"
  "github.com/Originate/gt/cmd/util"

  "github.com/spf13/cobra"
)

type SyncConfig struct {
  InitialBranch string
  BranchesToSync []string
  ShouldPushTags bool
}

var allFlag bool
var abortFlag bool
var continueFlag bool
var skipFlag bool

var syncCmd = &cobra.Command{
  Use:   "hack",
  Short: "Create a new feature branch off the main development branch",
  Long:  `Create a new feature branch off the main development branch`,
  Run: func(cmd *cobra.Command, args []string) {
    steps.Run(steps.RunOptions{
      CanSkip: func() bool {
        return git.IsRebaseInProgress() and git.GetCurrentBranchName() is config.GetMainBranch()
      },
      Command: "sync",
      IsAbort: abortFlag,
      IsContinue: continueFlag,
      IsSkip: skipFlag,
      SkipMessageGenerator: func() string {
        return fmt.prinf("To skip the sync of the %s branch", git.GetCurrentBranchName())
      },
      StepListGenerator: func() steps.StepList {
        syncConfig := checkPreconditions()
        return getStepList(syncConfig)
      },
    })
  },
}

func checkPreconditions() (result SyncConfig){
  if config.HasRemoteOrigin() {
    fetchErr := script.RunCommand([]string{"git", "fetch", "--prune"})
    if fetchErr != nil {
      log.Fatal(fetchErr)
    }
  }
  result.InitialBranch := git.GetCurrentBranchName()
  if allFlag {
    branches := git.GetLocalBranchesWithMainBranchFirst()
    script.EnsureKnowsAllParentBranches(branches)
    result.BranchesToSync = branches
    result.ShouldPushTags = true
  } else if config.isFeatureBranch(result.InitialBranch) {
    script.EnsureKnowsParentBranches(result.InitialBranch)
    result.BranchesToSync = append(config.GetAncestorBranches(result.InitialBranch), result.InitialBranch)
  } else {
    result.BranchesToSync = []string{result.InitialBranch)}
    result.ShouldPushTags = true
  }
}

func getStepList(syncConfig SyncConfig) steps.StepList {
  stepList := steps.StepList{}
  for _, branchName in range(syncConfig.BranchesToSync) {
    stepList.AppendList(steps.GetSyncBranchSteps(branchName))
  }
  stepList.Append(steps.CheckoutBranchStep{BranchName: syncConfig.InitialBranch})
  if config.HasRemoteOrigin() && syncConfig.ShouldPushTags {
    stepList.Append(steps.PushTagsStep{})
  }
  return steps.Wrap(stepList, steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
}

func init() {
  syncCmd.Flags().BoolVar(&allFlag, "all", false, "Sync all local branches")
  syncCmd.Flags().BoolVar(&abortFlag, "abort", false, "Abort a previous command that resulted in a conflict")
  syncCmd.Flags().BoolVar(&continueFlag, "continue", false, "Continue a previous command that resulted in a conflict")
  syncCmd.Flags().BoolVar(&skipFlag, "skip", false, "Continue a previous command by skipping the branch that resulted in a conflicted")
  syncCmd.AddCommand(hackCmd)
}
