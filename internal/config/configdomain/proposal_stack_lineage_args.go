package configdomain

import (
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// ProposalStackLineageArgs
type ProposalStackLineageArgs struct {
	// AfterStackDisplay
	// A set of texts to appear after the main stack information is displayed.
	// The strings in this area appear in the order of the slice after the stack lineage is displayed.
	AfterStackDisplay []string
	// BeforeStackDisplay
	// A set of texts to appear before the main stack information is displayed.
	// The strings in this area appear in the order of the slice before the stack lineage is displayed.
	BeforeStackDisplay []string
	// Connector
	// The current forge connector.
	Connector Option[forgedomain.Connector]
	// CurrentBranch
	// Informs the builder which branch is the current. This is used to determine when
	// the current branch indicator should be used AND how the stack hierarchy is computed.
	CurrentBranch gitdomain.LocalBranchName
	// CurrentBranchIndicator
	// Special character used to denote the current branch's proposal (if there is one).
	CurrentBranchIndicator string
	// IndentMarker
	// Controls the marker following an indent.
	IndentMarker string
	// Lineage
	// The lineage data for the current repository
	Lineage Lineage
	// MainAndPerennialBranches
	// These branches will not be searched for having a proposal attached to them
	// when building the stack lineage.
	MainAndPerennialBranches Option[gitdomain.LocalBranchNames]
}
