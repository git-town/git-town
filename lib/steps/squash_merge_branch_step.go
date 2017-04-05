package steps

import (
	"log"

	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/prompt"
	"github.com/Originate/git-town/lib/script"
	"github.com/Originate/git-town/lib/util"
)

type SquashMergeBranchStep struct {
	BranchName    string
	CommitMessage string
}

func (step SquashMergeBranchStep) CreateAbortStep() Step {
	return DiscardOpenChangesStep{}
}

func (step SquashMergeBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step SquashMergeBranchStep) CreateUndoStep() Step {
	return RevertCommitStep{}
}

func (step SquashMergeBranchStep) Run() error {
	err := script.RunCommand("git", "merge", "--squash", step.BranchName)
	if err != nil {
		log.Fatal("Error squash merging:", err)
	}
	commitCmd := []string{"git", "commit"}
	if step.CommitMessage != "" {
		commitCmd = append(commitCmd, "-m", step.CommitMessage)
	}
	author := prompt.GetSquashCommitAuthor(step.BranchName)
	if author != git.GetLocalAuthor() {
		commitCmd = append(commitCmd, "--author", author)
	}
	util.GetCommandOutput("sed", "-i", "-e", "s/^/# /g", ".git/SQUASH_MSG")
	return script.RunCommand(commitCmd...)
}
