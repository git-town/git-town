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
    util.ExitWithErrorMessage(fmt.Sprintf("A branch named '%s' already exists", branchName))
  }
}


func GetCurrentBranchName() string {
  if IsRebaseInProgress() {
    return getCurrentBranchNameDuringRebase()
  } else {
    return util.GetCommandOutput([]string{"git", "rev-parse", "--abbrev-ref", "HEAD"})
  }
}


func GetLocalBranchesWithMainBranchFirst() (result string[]) {


  // function local_branches {
  //   git branch | tr -d ' ' | sed 's/\*//g'
  // }
  //
  //
  // # Returns the names of local branches
  // function local_branches_with_main_first {
  //   if [ -n "$MAIN_BRANCH_NAME" ]; then
  //     echo "$MAIN_BRANCH_NAME"
  //   fi
  //   local_branches_without_main
  // }
  //
  //
  // # Returns the names of local branches without the main branch
  // function local_branches_without_main {
  //   local_branches | grep -v "^$MAIN_BRANCH_NAME\$"
  // }
}


func GetTrackingBranchName(branchName string) string {
  return "origin/" + branchName
}


func HasBranch(branchName string) bool {
  for _, branch := range(getAllBranches()) {
    if branch == branchName {
      return true
    }
  }
  return false
}


func HasTrackingBranch(branchName string) bool {
  trackingBranchName := GetTrackingBranchName(branchName)
  output := util.GetCommandOutput([]string{"git", "branch", "-r"})
  for _, line := range(strings.Split(output, "\n")) {
    if strings.TrimSpace(line) == trackingBranchName {
      return true
    }
  }
  return false
}


func ShouldBranchBePushed(branchName string) bool {
  trackingBranchName := GetTrackingBranchName(branchName)
  output := util.GetCommandOutput([]string{"git", "rev-list", "--left-right", branchName + "..." + trackingBranchName})
  return output != ""
}


// Helpers

func getAllBranches() (result []string) {
  output := util.GetCommandOutput([]string{"git", "branch", "-a"})
  for _, line := range(strings.Split(output, "\n")) {
    if (strings.Contains(line, "remotes/origin/HEAD ->")) {
      continue
    }
    line = strings.Trim(line, "* ")
    line = strings.TrimSpace(line)
    line = strings.Replace(line, "remotes/origin/", "", 1)
    result = append(result, line)
  }
  return
}


func getLocalBranches() (result []string) {
  output := util.GetCommandOutput([]string{"git", "branch"})
  for _, line := range(strings.Split(output, "\n")) {
    line = strings.Trim(line, "* ")
    line = strings.TrimSpace(line)
    result = append(result, line)
  }
  return
}


func getCurrentBranchNameDuringRebase() string {
  filename := fmt.Sprintf("%s/.git/rebase-apply/head-name", GetRootDirectory())
  rawContent, err := ioutil.ReadFile(filename)
  if err != nil {
    log.Fatal(err)
  }
  content := strings.TrimSpace(string(rawContent))
  return strings.Replace(content, "refs/heads/", "", -1)
}
