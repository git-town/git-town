package cmd

import (
  "fmt"
  "log"

  "github.com/Originate/gt/cmd/config"
  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/script"
  "github.com/Originate/gt/cmd/steps"

  "github.com/spf13/cobra"
)

type SyncConfig struct {
  InitialBranch string
  BranchesToSync []string
  ShouldPushTags bool
}

type SyncFlags struct {
  All bool
  Abort bool
  Continue bool
  Skip bool
}

var syncFlags SyncFlags

var syncCmd = &cobra.Command{
  Use:   "hack",
  Short: "Create a new feature branch off the main development branch",
  Long:  `Create a new feature branch off the main development branch`,
  Run: func(cmd *cobra.Command, args []string) {
    steps.Run(steps.RunOptions{
      CanSkip: func() bool {
        return git.IsRebaseInProgress() && git.GetCurrentBranchName() == config.GetMainBranch()
      },
      Command: "sync",
      IsAbort: syncFlags.Abort,
      IsContinue: syncFlags.Continue,
      IsSkip: syncFlags.Skip,
      SkipMessageGenerator: func() string {
        return fmt.Sprintf("To skip the sync of the %s branch", git.GetCurrentBranchName())
      },
      StepListGenerator: func() steps.StepList {
        syncConfig := checkSyncPreconditions()
        return getSyncStepList(syncConfig)
      },
    })
  },
}

func checkSyncPreconditions() (result SyncConfig){
  if config.HasRemoteOrigin() {
    fetchErr := script.RunCommand([]string{"git", "fetch", "--prune"})
    if fetchErr != nil {
      log.Fatal(fetchErr)
    }
  }
  result.InitialBranch = git.GetCurrentBranchName()
  if syncFlags.All {
    branches := git.GetLocalBranchesWithMainBranchFirst()
    prompt.EnsureKnowsAllParentBranches(branches)
    result.BranchesToSync = branches
    result.ShouldPushTags = true
  } else if config.IsFeatureBranch(result.InitialBranch) {
    prompt.EnsureKnowsParentBranches(result.InitialBranch)
    result.BranchesToSync = append(config.GetAncestorBranches(result.InitialBranch), result.InitialBranch)
  } else {
    result.BranchesToSync = []string{result.InitialBranch}
    result.ShouldPushTags = true
  }
  return
}

func getSyncStepList(syncConfig SyncConfig) steps.StepList {
  stepList := steps.StepList{}
  for _, branchName := range(syncConfig.BranchesToSync) {
    stepList.AppendList(steps.GetSyncBranchSteps(branchName))
  }
  stepList.Append(steps.CheckoutBranchStep{BranchName: syncConfig.InitialBranch})
  if config.HasRemoteOrigin() && syncConfig.ShouldPushTags {
    stepList.Append(steps.PushTagsStep{})
  }
  return steps.Wrap(stepList, steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
}

func init() {
  syncCmd.Flags().BoolVar(&syncFlags.All, "all", false, "Sync all local branches")
  syncCmd.Flags().BoolVar(&syncFlags.Abort, "abort", false, "Abort a previous command that resulted in a conflict")
  syncCmd.Flags().BoolVar(&syncFlags.Continue, "continue", false, "Continue a previous command that resulted in a conflict")
  syncCmd.Flags().BoolVar(&syncFlags.Skip, "skip", false, "Continue a previous command by skipping the branch that resulted in a conflicted")
  syncCmd.AddCommand(hackCmd)
}
