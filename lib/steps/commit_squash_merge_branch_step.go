package steps

import (
	"strings"

	"github.com/Originate/git-town/lib/script"
	"github.com/Originate/git-town/lib/util"
)

type CommitSquashMergeBranchStep struct {
	BranchName    string
	CommitMessage string
}

func (step CommitSquashMergeBranchStep) CreateAbortStep() Step {
	return DiscardOpenChanges{}
}

func (step CommitSquashMergeBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step CommitSquashMergeBranchStep) CreateUndoStep() Step {
	return NoOpStep{}
}

func (step CommitSquashMergeBranchStep) Run() error {
  commitOptions := ["-m", step.CommitMessage]
	author := prompt.GetSquashCommitAuthor(step.BranchName)
	if author != git.LocalAuthor() {
		commitOptions = append(commitOptions, "--author",  "\"" + author + "\"")
	}
	util.GetCommandOutput("sed -i -e 's/^/# /g' .git/SQUASH_MSG")
	return script.RunCommand("git", "commit", commitOptions...)
}

func (step CommitSquashMergeBranchStep) ShouldAbortOnError() (bool, string) {
  return true, "Aborted because commit exited with error"
}
