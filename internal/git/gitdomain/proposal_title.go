package gitdomain

import "github.com/git-town/git-town/v23/internal/gohacks/stringss"

// ProposalTitle is the title of a proposal
type ProposalTitle stringss.Trimmed

// String implements the fmt.Stringer interface.
func (self ProposalTitle) String() string {
	return string(self)
}
