package step

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
)

// ConnectorMergeProposal squash merges the branch with the given name into the current branch.
type ConnectorMergeProposal struct {
	Branch                    domain.LocalBranchName
	CommitMessage             string
	ProposalMessage           string
	enteredEmptyCommitMessage bool
	mergeError                error
	ProposalNumber            int
	Empty
}

func (step *ConnectorMergeProposal) CreateAbortSteps() []Step {
	if step.enteredEmptyCommitMessage {
		return []Step{&DiscardOpenChanges{}}
	}
	return []Step{}
}

func (step *ConnectorMergeProposal) CreateAutomaticAbortError() error {
	if step.enteredEmptyCommitMessage {
		return fmt.Errorf(messages.ShipAbortedMergeError)
	}
	return step.mergeError
}

func (step *ConnectorMergeProposal) Run(args RunArgs) error {
	commitMessage := step.CommitMessage
	//nolint:nestif
	if commitMessage == "" {
		// Allow the user to enter the commit message as if shipping without a connector
		// then revert the commit since merging via the connector will perform the actual squash merge.
		step.enteredEmptyCommitMessage = true
		err := args.Runner.Frontend.SquashMerge(step.Branch)
		if err != nil {
			return err
		}
		err = args.Runner.Backend.CommentOutSquashCommitMessage(step.ProposalMessage + "\n\n")
		if err != nil {
			return fmt.Errorf(messages.SquashMessageProblem, err)
		}
		err = args.Runner.Frontend.StartCommit()
		if err != nil {
			return err
		}
		commitMessage, err = args.Runner.Backend.LastCommitMessage()
		if err != nil {
			return err
		}
		err = args.Runner.Frontend.DeleteLastCommit()
		if err != nil {
			return err
		}
		step.enteredEmptyCommitMessage = false
	}
	_, step.mergeError = args.Connector.SquashMergeProposal(step.ProposalNumber, commitMessage)
	return step.mergeError
}

// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
// automatically abort if it errors.
func (step *ConnectorMergeProposal) ShouldAutomaticallyAbortOnError() bool {
	return true
}
