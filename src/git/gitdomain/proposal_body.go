package gitdomain

// ProposalBody is the body of a proposal
type ProposalBody string

// String implements the fmt.Stringer interface.
func (self ProposalBody) String() string {
	return string(self)
}
