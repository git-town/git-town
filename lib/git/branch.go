package git

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/Originate/git-town/lib/gitconfig"
	"github.com/Originate/git-town/lib/util"
)

func DoesBranchHaveUnmergedCommits(branchName string) bool {
	return util.GetCommandOutput("git", "log", gitconfig.GetMainBranch()+".."+branchName) != ""
}

func EnsureDoesNotHaveBranch(branchName string) {
	if HasBranch(branchName) {
		util.ExitWithErrorMessage(fmt.Sprintf("A branch named '%s' already exists", branchName))
	}
}

func EnsureHasBranch(branchName string) {
	if !HasBranch(branchName) {
		util.ExitWithErrorMessage(fmt.Sprintf("There is no branch named '%s'", branchName))
	}
}

func GetCurrentBranchName() string {
	if IsRebaseInProgress() {
		return getCurrentBranchNameDuringRebase()
	} else {
		return util.GetCommandOutput("git", "rev-parse", "--abbrev-ref", "HEAD")
	}
}

func GetLocalBranches() (result []string) {
	output := util.GetCommandOutput("git", "branch")
	for _, line := range strings.Split(output, "\n") {
		line = strings.Trim(line, "* ")
		line = strings.TrimSpace(line)
		result = append(result, line)
	}
	return
}

func GetLocalBranchesWithDeletedTrackingBranches() (result []string) {
	output := util.GetCommandOutput("git", "branch", "-vv")
	for _, line := range strings.Split(output, "\n") {
		line = strings.Trim(line, "* ")
		parts := strings.SplitN(line, " ", 2)
		branchName := parts[0]
		deleteTrackingBranchStatus := fmt.Sprintf("[%s: gone]", GetTrackingBranchName(branchName))
		if strings.Contains(parts[1], deleteTrackingBranchStatus) {
			result = append(result, branchName)
		}
	}
	return
}

func GetLocalBranchesWithMainBranchFirst() (result []string) {
	mainBranch := gitconfig.GetMainBranch()
	result = append(result, mainBranch)
	for _, branch := range GetLocalBranches() {
		if branch != mainBranch {
			result = append(result, branch)
		}
	}
	return
}

func GetTrackingBranchName(branchName string) string {
	return "origin/" + branchName
}

func HasBranch(branchName string) bool {
	output := util.GetCommandOutput("git", "branch", "-a")
	for _, line := range strings.Split(output, "\n") {
		line = strings.Trim(line, "* ")
		line = strings.TrimSpace(line)
		line = strings.Replace(line, "remotes/origin/", "", 1)
		if line == branchName {
			return true
		}
	}
	return false
}

func HasLocalBranch(branchName string) bool {
	return util.DoesStringArrayContain(GetLocalBranches(), branchName)
}

func HasTrackingBranch(branchName string) bool {
	trackingBranchName := GetTrackingBranchName(branchName)
	output := util.GetCommandOutput("git", "branch", "-r")
	for _, line := range strings.Split(output, "\n") {
		if strings.TrimSpace(line) == trackingBranchName {
			return true
		}
	}
	return false
}

func ShouldBranchBePushed(branchName string) bool {
	trackingBranchName := GetTrackingBranchName(branchName)
	output := util.GetCommandOutput("git", "rev-list", "--left-right", branchName+"..."+trackingBranchName)
	return output != ""
}

// Helpers

func getCurrentBranchNameDuringRebase() string {
	filename := fmt.Sprintf("%s/.git/rebase-apply/head-name", GetRootDirectory())
	rawContent, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	content := strings.TrimSpace(string(rawContent))
	return strings.Replace(content, "refs/heads/", "", -1)
}
