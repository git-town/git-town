package forgedomain

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type ProposalData struct {
	Body         Option[string]
	MergeWithAPI bool
	Number       int
	Source       gitdomain.LocalBranchName
	Target       gitdomain.LocalBranchName
	Title        string
	URL          string
}

func (self ProposalData) Data() ProposalData {
	return self
}

type BitbucketCloudProposalData struct {
	ProposalData
	CloseSourceBranch bool
	Draft             bool
}

func (self BitbucketCloudProposalData) Data() ProposalData {
	return self.ProposalData
}
