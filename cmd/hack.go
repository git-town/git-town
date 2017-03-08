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

var abortFlag bool
var continueFlag bool

var hackCmd = &cobra.Command{
  Use:   "hack",
  Short: "Create a new feature branch off the main development branch",
  Long:  `Create a new feature branch off the main development branch`,
  Run: func(cmd *cobra.Command, args []string) {
    steps.Run(steps.RunOptions{
      CanSkip: func() bool { return false },
      Command: "hack",
      IsAbort: abortFlag,
      IsContinue: continueFlag,
      IsSkip: false,
      SkipMessageGenerator: func() string { return "" },
      StepListGenerator: func() steps.StepList {
        targetBranchName := checkPreconditions(args)
        return getStepList(targetBranchName)
      },
    })
  },
}

func checkPreconditions(args []string) string {
  if len(args) == 0 {
    util.ExitWithErrorMessage("No branch name provided.")
  }
  targetBranchName := args[0]
  if config.HasRemoteOrigin() {
    fetchErr := script.RunCommand([]string{"git", "fetch", "--prune"})
    if fetchErr != nil {
      log.Fatal(fetchErr)
    }
  }
  git.EnsureDoesNotHaveBranch(targetBranchName)
  return targetBranchName
}

func getStepList(targetBranchName string) steps.StepList {
  mainBranchName := config.GetMainBranch()
  stepList := steps.StepList{}
  stepList.AppendList(steps.GetSyncBranchSteps(mainBranchName))
  stepList.Append(steps.CreateAndCheckoutBranchStep{BranchName: targetBranchName, ParentBranchName: mainBranchName})
  if config.HasRemoteOrigin() && config.ShouldHackPush() {
    stepList.Append(steps.CreateTrackingBranchStep{BranchName: targetBranchName})
  }
  return steps.Wrap(stepList, steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
}

func init() {
  hackCmd.Flags().BoolVar(&abortFlag, "abort", false, "Abort a previous command that resulted in a conflict")
  hackCmd.Flags().BoolVar(&continueFlag, "continue", false, "Continue a previous command that resulted in a conflict")
  RootCmd.AddCommand(hackCmd)
}
