package gitdomain

// ProposalBodyTemplateFile is the path to a file containing a template for generating a proposal body
type ProposalBodyTemplateFile string

// String implements the fmt.Stringer interface.
func (self ProposalBodyTemplateFile) String() string {
	return string(self)
}
