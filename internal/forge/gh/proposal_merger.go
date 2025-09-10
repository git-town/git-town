package gh

import (
	"strconv"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
)

var _ forgedomain.ProposalMerger = ghConnector

func (self Connector) SquashMergeProposal(number int, message gitdomain.CommitMessage) error {
	return self.Frontend.Run("gh", "pr", "merge", "--squash", "--body="+message.String(), strconv.Itoa(number))
}
