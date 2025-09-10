package glab

import (
	"strconv"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
)

func (self Connector) SquashMergeProposal(number int, message gitdomain.CommitMessage) error {
	return self.Frontend.Run("glab", "mr", "merge", "--squash", "--body="+message.String(), strconv.Itoa(number))
}
