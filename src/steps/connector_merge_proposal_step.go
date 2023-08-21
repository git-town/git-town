package steps

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
)

// ConnectorMergeProposalStep squash merges the branch with the given name into the current branch.
type ConnectorMergeProposalStep struct {
	Branch                    domain.LocalBranchName
	CommitMessage             string
	ProposalMessage           string
	enteredEmptyCommitMessage bool
	mergeError                error
	mergeSha                  domain.SHA
	ProposalNumber            int
	EmptyStep
}

func (step *ConnectorMergeProposalStep) CreateAbortStep() []Step {
	if step.enteredEmptyCommitMessage {
		return []Step{&DiscardOpenChangesStep{}}
	}
	return []Step{}
}

func (step *ConnectorMergeProposalStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&RevertCommitStep{Sha: step.mergeSha}}, nil
}

func (step *ConnectorMergeProposalStep) CreateAutomaticAbortError() error {
	if step.enteredEmptyCommitMessage {
		return fmt.Errorf(messages.ShipAbortedMergeError)
	}
	return step.mergeError
}

func (step *ConnectorMergeProposalStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	commitMessage := step.CommitMessage
	//nolint:nestif
	if commitMessage == "" {
		// Allow the user to enter the commit message as if shipping without a connector
		// then revert the commit since merging via the connector will perform the actual squash merge.
		step.enteredEmptyCommitMessage = true
		err := run.Frontend.SquashMerge(step.Branch)
		if err != nil {
			return err
		}
		err = run.Backend.CommentOutSquashCommitMessage(step.ProposalMessage + "\n\n")
		if err != nil {
			return fmt.Errorf(messages.SquashMessageProblem, err)
		}
		err = run.Frontend.StartCommit()
		if err != nil {
			return err
		}
		commitMessage, err = run.Backend.LastCommitMessage()
		if err != nil {
			return err
		}
		err = run.Frontend.DeleteLastCommit()
		if err != nil {
			return err
		}
		step.enteredEmptyCommitMessage = false
	}
	step.mergeSha, step.mergeError = connector.SquashMergeProposal(step.ProposalNumber, commitMessage)
	return step.mergeError
}

// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
// automatically abort if it errors.
func (step *ConnectorMergeProposalStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
