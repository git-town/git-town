package gitdomain

import "github.com/git-town/git-town/v23/internal/gohacks/stringss"

// ProposalBodyFile is the body of a proposal.
type ProposalBodyFile stringss.TrimmedString

// ShouldReadStdin indicates whether the body should be read from STDIN.
func (self ProposalBodyFile) ShouldReadStdin() bool {
	return self == "-"
}

// String implements the fmt.Stringer interface.
func (self ProposalBodyFile) String() string {
	return string(self)
}
