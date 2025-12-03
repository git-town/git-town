package gitdomain

import (
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// ProposalBody is the body of a proposal
type ProposalBody string

// String implements the fmt.Stringer interface.
func (self ProposalBody) String() string {
	return string(self)
}

func NewProposalBodyOpt(text string) Option[ProposalBody] {
	if text == "" {
		return None[ProposalBody]()
	}
	return Some(ProposalBody(text))
}
