package steps

import (
	"fmt"

	"github.com/git-town/git-town/v8/src/dialog"
	"github.com/git-town/git-town/v8/src/git"
	"github.com/git-town/git-town/v8/src/hosting"
)

// SquashMergeStep squash merges the branch with the given name into the current branch.
type SquashMergeStep struct {
	EmptyStep
	Branch        string
	CommitMessage string
	Parent        string
}

func (step *SquashMergeStep) CreateAbortStep() Step {
	return &DiscardOpenChangesStep{}
}

func (step *SquashMergeStep) CreateUndoStep(backend *git.BackendCommands) (Step, error) {
	currentSHA, err := backend.CurrentSha()
	if err != nil {
		return nil, err
	}
	return &RevertCommitStep{Sha: currentSHA}, nil
}

func (step *SquashMergeStep) CreateAutomaticAbortError() error {
	return fmt.Errorf("aborted because commit exited with error")
}

func (step *SquashMergeStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	err := run.Frontend.SquashMerge(step.Branch)
	if err != nil {
		return err
	}
	branchAuthors, err := run.Backend.BranchAuthors(step.Branch, step.Parent)
	if err != nil {
		return err
	}
	author, err := dialog.SelectSquashCommitAuthor(step.Branch, branchAuthors)
	if err != nil {
		return fmt.Errorf("error getting squash commit author: %w", err)
	}
	repoAuthor, err := run.Backend.Author()
	if err != nil {
		return fmt.Errorf("cannot determine repo author: %w", err)
	}
	if err = run.Backend.CommentOutSquashCommitMessage(""); err != nil {
		return fmt.Errorf("cannot comment out the squash commit message: %w", err)
	}
	if repoAuthor == author {
		author = ""
	}
	return run.Frontend.Commit(step.CommitMessage, author)
}

func (step *SquashMergeStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
