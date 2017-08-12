package steps

import (
	"log"
	"strings"

	"github.com/Originate/git-town/src/drivers"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/script"
	"github.com/Originate/git-town/src/util"
)

// DriverMergePullRequestStep squash merges the branch with the given name into the current branch
type DriverMergePullRequestStep struct {
	NoOpStep
	BranchName    string
	CommitMessage string
}

// Run executes this step.
func (step DriverMergePullRequestStep) Run() error {
	commitMessage := step.CommitMessage
	if commitMessage != "" {
		err := script.RunCommand("git", "merge", "--squash", step.BranchName)
		if err != nil {
			log.Fatal("Error squash merging:", err)
		}
		util.GetCommandOutput("sed", "-i", "-e", "s/^/# /g", ".git/SQUASH_MSG")
		err = script.RunCommand("git", "commit")
		if err != nil {
			return err
		}
		commitMessage = util.GetCommandOutput("git", "log", "-1", "--format=%B")
		err = script.RunCommand("git", "reset", "--hard", "HEAD~1")
		if err != nil {
			log.Fatal("Error resetting the main branch", err)
		}
	}
	driver := drivers.GetCodeHostingDriver()
	repository := git.GetURLRepositoryName(git.GetRemoteOriginURL())
	repositoryParts := strings.SplitN(repository, "/", 2)
	return driver.MergePullRequest(drivers.MergePullRequestOptions{
		Branch:        step.BranchName,
		CommitMessage: commitMessage,
		ParentBranch:  git.GetCurrentBranchName(),
		Owner:         repositoryParts[0],
		Repository:    repositoryParts[1],
	})
}
