package cmd

import (
  "fmt"
  "log"

  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/script"
  "github.com/Originate/gt/cmd/utils"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
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

// function ensure_has_target_branch {
//   if [ -z "$target_branch_name" ]; then
//     echo_error_header
//     echo_error "No branch name provided."
//     exit_with_error newline
//   fi
// }
//
//
// function preconditions {
//   target_branch_name=$1
//   ensure_has_target_branch
//
//   if [ "$(has_remote_url)" = true ]; then
//     fetch
//   fi
//
//   ensure_does_not_have_branch "$target_branch_name"
//
//   export RUN_IN_GIT_ROOT=true
//   export STASH_OPEN_CHANGES=true
// }
