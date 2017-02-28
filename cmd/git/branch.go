package git

import (
  "fmt"
  "io/ioutil"
  "log"
  "strings"

  "github.com/Originate/gt/cmd/util"
)


func EnsureDoesNotHaveBranch(branchName string) {
  if HasBranch(branchName) {
    message := fmt.Sprintf("A branch named '%s' already exists", branchName)
    util.ExitWithErrorMessage(message)
  }
}


func GetCurrentBranchName() string {
  if IsRebaseInProgress() {
    filename := fmt.Sprintf("%s/.git/rebase-apply/head-name", GetRootDirectory())
    content, err := ioutil.ReadFile(filename)
    if err != nil {
      log.Fatal(err)
    }
    return strings.Replace(strings.TrimSpace(string(content)), "refs/heads/", "", -1)
  } else {
    cmd := []string{"git", "rev-parse", "--abbrev-ref", "HEAD"}
    return util.GetCommandOutput(cmd)
  }
}


func GetTrackingBranchName(branchName string) string {
  return fmt.Sprintf("origin/%s", branchName)
}


func HasBranch(branchName string) bool {
  cmd := []string{"git", "branch", "-a"}
  lines := strings.Split(util.GetCommandOutput(cmd), "\n")
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



func HasTrackingBranch(branchName string) bool {
  trackingBranchName := GetTrackingBranchName(branchName)
  output := util.GetCommandOutput([]string{"git", "branch", "-r"})
  lines := strings.Split(output, "\n")
  for i := 0; i < len(lines); i++ {
    line := lines[i]
    line = strings.TrimSpace(line)
    if line == trackingBranchName {
      return true
    }
  }
  return false
}


func ShouldBranchBePushed(branchName string) bool {
  trackingBranchName := GetTrackingBranchName(branchName)
  output := util.GetCommandOutput([]string{"git", "rev-list", "--left-right", fmt.Sprintf("%s...%s", branchName, trackingBranchName)})
  return output != ""
}
