package gitdomain

import (
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

// ProposalBody is the body of a proposal
type ProposalBody stringss.Trimmed

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
