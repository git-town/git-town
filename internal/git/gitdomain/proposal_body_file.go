package gitdomain

// ProposalBodyFile is the body of a proposal.
type ProposalBodyFile string

// ShouldReadStdin indicates whether the body should be read from STDIN.
func (self ProposalBodyFile) ShouldReadStdin() bool {
	return self == "-"
}

// String implements the fmt.Stringer interface.
func (self ProposalBodyFile) String() string {
	return string(self)
}
