package configdomain

type ProposalLineage string

const (
	ProposalLineageNone            ProposalLineage = "none"
	ProposalLineageComment         ProposalLineage = "comment"
	ProposalLineageTerminalDisplay ProposalLineage = "terminal-display"
)

func (self ProposalLineage) String() string {
	return string(self)
}
