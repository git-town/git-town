package forgedomain

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type ProposalCache struct {
	cache map[string]Proposal
}

func (self *ProposalCache) Get(branch, target gitdomain.LocalBranchName) Option[Proposal] {
	result, found := self.cache[self.key(branch, target)]
	if !found {
		return None[Proposal]()
	}
	return Some(result)
}

func (self *ProposalCache) Set(branch, target gitdomain.LocalBranchName, proposal Proposal) {
	self.cache[self.key(branch, target)] = proposal
}

func (self *ProposalCache) key(branch, target gitdomain.LocalBranchName) string {
	return branch.String() + "--" + target.String()
}
