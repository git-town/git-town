package git

import (
  "fmt"
  "os"
  "strings"

  "github.com/Originate/gt/cmd/util"
)


func EnsureDoesNotHaveConflicts() {
  if HasConflicts() {
    util.ExitWithErrorMessage("You must resolve the conflicts before continuing.")
  }
}


func GetRootDirectory() string {
  return util.GetCommandOutput([]string{"git", "rev-parse", "--show-toplevel"})
}


func HasConflicts() bool {
  output := util.GetCommandOutput([]string{"git", "status"})
  return strings.Contains(output, "Unmerged paths")
}


func HasOpenChanges() bool {
  output := util.GetCommandOutput([]string{"git", "status", "--porcelain"})
  return output != ""
}


func IsMergeInProgress() bool {
  _, err := os.Stat(fmt.Sprintf("%s/.git/MERGE_HEAD", GetRootDirectory()))
  return err == nil
}


func IsRebaseInProgress() bool {
  status := util.GetCommandOutput([]string{"git", "status"})
  return strings.Contains(status, "rebase in progress")
}
