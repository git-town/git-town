package git

import (
  "fmt"
  "strings"

  "github.com/Originate/gt/cmd/utils"
)

func EnsureDoesNotHaveBranch(branchName string) {
  if HasBranch(branchName) {
    message := fmt.Sprintf("A branch named '%s' already exists", branchName)
    utils.ExitWithErrorMessage(message)
  }
}

func GetCurrentBranchName() string {
  cmd := []string{"git", "rev-parse", "--abbrev-ref", "HEAD"}
  return utils.GetCommandOutput(cmd)
}

func HasBranch(branchName string) bool {
  cmd := []string{"git", "branch", "-a"}
  lines := strings.Split(utils.GetCommandOutput(cmd), "\n")
  for i := 0; i < len(lines); i++ {
    line := lines[i]
    line = strings.Trim(line, "* ")
    line = strings.TrimSpace(line)
    line = strings.Replace(line, "remotes/origin/", "", 1)
    if line == branchName {
      return true
    }
	}
  return false
}
