package configdomain

// indicates whether and how the user wants to embed the stack lineage of the respective branch in proposals
type DisplayLineageInProposals string

const (
	DisplayLineageInProposalsNone = "none" // don't display lineage in proposals
	DisplayLineageInProposalsCI   = "ci"   // this team has set up https://github.com/git-town/action to do this
	DisplayLineageInProposalsCLI  = "cli"  // configures the Git Town CLI to embed lineage in proposals
)

func (self DisplayLineageInProposals) String() string {
	return string(self)
}
