package bitbucketcloud

import (
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

type BitbucketCloudProposal struct {
	body         Option[string]
	mergeWithAPI bool
	number       int
	source       gitdomain.LocalBranchName
	target       gitdomain.LocalBranchName
	title        string
	url          string
}

func (self BitbucketCloudProposal) Body() Option[string] {
	return self.body
}

func (self BitbucketCloudProposal) MergeWithAPI() bool {
	return self.mergeWithAPI
}

func (self BitbucketCloudProposal) Number() int {
	return self.number
}

func (self BitbucketCloudProposal) Source() gitdomain.LocalBranchName {
	return self.source
}

func (self BitbucketCloudProposal) Target() gitdomain.LocalBranchName {
	return self.target
}

func (self BitbucketCloudProposal) Title() string {
	return self.title
}

func (self BitbucketCloudProposal) URL() string {
	return self.url
}
