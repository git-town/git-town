package cmd

import (
  "fmt"
  "log"

  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/script"
  "github.com/Originate/gt/cmd/util"
  "github.com/Originate/gt/cmd/step"

  "github.com/spf13/cobra"
)


var hackCmd = &cobra.Command{
  Use:   "hack",
  Short: "Create a new feature branch off the main development branch",
  Long:  `Create a new feature branch off the main development branch`,
  Run: func(cmd *cobra.Command, args []string) {
    if len(args) == 0 {
      util.ExitWithErrorMessage("No branch name provided.")
    }
    targetBranchName := args[0]
    fetchCmd := []string{"git", "fetch", "--prune"}
    fetchErr := script.RunCommand(fetchCmd)
    if fetchErr != nil {
      log.Fatal(fetchErr)
    }
    git.EnsureDoesNotHaveBranch(targetBranchName)
    fmt.Println()
    var p []string
    config := util.Config{
      HasRemote: true,
      MainBranchName: "master",
      PerennialBranchNames: p,
      PullBranchStrategy: "rebase",
    }
    var steps []step.Step
    steps = append(steps, step.GetSyncBranchSteps("master", config)...)
    for i := 0; i < len(steps); i++ {
      err := steps[i].Run()
      if err != nil {
        log.Fatal(err)
      }
    }
    // echo "create_and_checkout_feature_branch $target_branch_name $MAIN_BRANCH_NAME"
    // echo_if_all_true "create_tracking_branch $target_branch_name" "$HAS_REMOTE" "$(hack_should_push)"
  },
}


func init() {
  RootCmd.AddCommand(hackCmd)
}
