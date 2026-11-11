package forgedomain

import (
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

type ProposalData struct {
	Active       bool // whether the proposal is open
	Body         Option[gitdomain.ProposalBody]
	MergeWithAPI bool
	Number       ProposalNumber
	Source       gitdomain.LocalBranchName
	Target       gitdomain.LocalBranchName
	Title        gitdomain.ProposalTitle
	URL          string
}

func (self ProposalData) Data() ProposalData {
	return self
}

type BitbucketCloudProposalData struct {
	CloseSourceBranch bool
	Draft             bool
	ProposalData
}

func (self BitbucketCloudProposalData) Data() ProposalData {
	return self.ProposalData
}
