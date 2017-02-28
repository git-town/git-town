package cmd

import (
  "log"

  "github.com/Originate/gt/cmd/config"
  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/script"
  "github.com/Originate/gt/cmd/step"
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
    if abortFlag || continueFlag {
      runResult := step.Import("hack")
      if abortFlag {
        step.Run(append([]step.Step{runResult.AbortStep}, runResult.UndoSteps...), []step.Step{}, "hack", "")
      } else {
        git.EnsureDoesNotHaveConflicts()
        step.Run(runResult.ContinueSteps, runResult.UndoSteps, "hack", "")
      }
    } else {
      // Preconditions
      if len(args) == 0 {
        util.ExitWithErrorMessage("No branch name provided.")
      }
      targetBranchName := args[0]
      if config.HasRemote() {
        fetchCmd := []string{"git", "fetch", "--prune"}
        fetchErr := script.RunCommand(fetchCmd)
        if fetchErr != nil {
          log.Fatal(fetchErr)
        }
      }
      git.EnsureDoesNotHaveBranch(targetBranchName)
      // Build Steps
      mainBranchName := config.GetMainBranch()
      var steps []step.Step
      steps = append(steps, step.GetSyncBranchSteps(mainBranchName)...)
      steps = append(steps, step.CreateAndCheckoutBranchStep{BranchName: targetBranchName, ParentBranchName: mainBranchName})
      if config.HasRemote() && config.ShouldHackPush() {
        steps = append(steps, step.CreateTrackingBranchStep{BranchName: targetBranchName})
      }
      steps = step.Wrap(steps, step.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
      // Run Steps
      step.Run(steps, []step.Step{}, "hack", "")
    }
  },
}


func init() {
  hackCmd.Flags().BoolVar(&abortFlag, "abort", false, "Abort a previous command that resulted in a conflict")
  hackCmd.Flags().BoolVar(&continueFlag, "continue", false, "Continue a previous command that resulted in a conflict")
  RootCmd.AddCommand(hackCmd)
}
