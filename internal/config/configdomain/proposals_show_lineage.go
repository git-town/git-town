package configdomain

// indicates whether and how proposals should display the stack lineage of the respective branch
type ProposalsShowLineage string

const (
	ProposalsShowLineageNone = "none" // don't display lineage in proposals
	ProposalsShowLineageCI   = "ci"   // this team has set up https://github.com/git-town/action to embed the stack lineage into proposals
	ProposalsShowLineageCLI  = "cli"  // the Git Town CLI should embed the lineage into proposals
)

func (self ProposalsShowLineage) String() string {
	return string(self)
}
