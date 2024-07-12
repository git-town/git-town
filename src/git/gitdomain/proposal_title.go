package gitdomain

// ProposalTitle is the title of a proposal
type ProposalTitle string

// String implements the fmt.Stringer interface.
func (self ProposalTitle) String() string {
	return string(self)
}
