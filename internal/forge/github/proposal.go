package github

import (
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

type Proposal struct {
	body         Option[string]
	mergeWithAPI bool
	number       int
	source       gitdomain.LocalBranchName
	target       gitdomain.LocalBranchName
	title        string
	url          string
}

func (self Proposal) Body() Option[string] {
	return self.body
}

func (self Proposal) MergeWithAPI() bool {
	return self.mergeWithAPI
}

func (self Proposal) Number() int {
	return self.number
}

func (self Proposal) Source() gitdomain.LocalBranchName {
	return self.source
}

func (self Proposal) Target() gitdomain.LocalBranchName {
	return self.target
}

func (self Proposal) Title() string {
	return self.title
}

func (self Proposal) URL() string {
	return self.url
}
