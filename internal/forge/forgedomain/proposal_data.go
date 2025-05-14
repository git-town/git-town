package forgedomain

import (
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	. "github.com/git-town/git-town/v20/pkg/prelude"
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

func (self ProposalData) GetBody() Option[string] {
	return self.Body
}

func (self ProposalData) GetMergeWithAPI() bool {
	return self.MergeWithAPI
}

func (self ProposalData) GetNumber() int {
	return self.Number
}

func (self ProposalData) GetSource() gitdomain.LocalBranchName {
	return self.Source
}

func (self ProposalData) GetTarget() gitdomain.LocalBranchName {
	return self.Target
}

func (self ProposalData) GetTitle() string {
	return self.Title
}

func (self ProposalData) GetURL() string {
	return self.URL
}

type BitbucketCloudProposalData struct {
	ProposalData
	CloseSourceBranch bool
	DestinationCommit string
	Draft             bool
	Message           string
	Reviewers         []string
	SourceRepository  string
}
