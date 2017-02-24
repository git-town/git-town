package cmd

import (
  "fmt"
  "log"

  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/script"
  "github.com/Originate/gt/cmd/utils"

  "github.com/spf13/cobra"
)


var hackCmd = &cobra.Command{
  Use:   "hack",
  Short: "Create a new feature branch off the main development branch",
  Long:  `Create a new feature branch off the main development branch`,
  Run: func(cmd *cobra.Command, args []string) {
    if len(args) == 0 {
      utils.ExitWithErrorMessage("No branch name provided.")
    }
    targetBranchName := args[0]
    fetchCmd := []string{"git", "fetch", "--prune"}
    fetchErr := script.RunCommand(fetchCmd)
    if fetchErr != nil {
      log.Fatal(fetchErr)
    }
    git.EnsureDoesNotHaveBranch(targetBranchName)
    fmt.Println()
  },
}


func init() {
  RootCmd.AddCommand(hackCmd)
}
