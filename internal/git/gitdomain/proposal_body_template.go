package gitdomain

// ProposalBodyTemplate is the template text for generating a proposal body
type ProposalBodyTemplate string

// String implements the fmt.Stringer interface.
func (self ProposalBodyTemplate) String() string {
	return string(self)
}
