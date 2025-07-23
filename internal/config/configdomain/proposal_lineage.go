package configdomain

type ProposalLineageIn string

const (
	ProposalLineageInNone                     ProposalLineageIn = "none"
	ProposalLineageOperationInProposalBody    ProposalLineageIn = "proposal_body"
	ProposalLineageOperationInProposalComment ProposalLineageIn = "proposal_comment"
	ProposalLineageInTerminal                 ProposalLineageIn = "terminal"
)

func (self ProposalLineageIn) String() string {
	return string(self)
}
