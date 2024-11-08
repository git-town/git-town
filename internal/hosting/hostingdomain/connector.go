package hostingdomain

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/gohacks/stringslice"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// Connector describes the activities that Git Town can perform on code hosting platforms.
// Individual implementations exist to talk to specific hosting platforms.
// Functions that might or might not be supported by a connector are implemented as higher-level functions,
// i.e. they return an option of the function to call.
// A `None“ value implies that the respective functionality isn't supported by this connector implementation.
type Connector interface {
	// DefaultProposalMessage provides the text that the form for creating new proposals
	// on the respective hosting platform is prepopulated with.
	DefaultProposalMessage(proposal Proposal) string

	// provides a function to load details about the proposal for the given branch into the given target branch.
	FindProposalFn() Option[func(branch, target gitdomain.LocalBranchName) (Option[Proposal], error)]

	// SearchProposals searches for a proposal that has the given branch as its head (source) branch.
	SearchProposalFn() Option[func(branch gitdomain.LocalBranchName) (Option[Proposal], error)]

	// SquashMergeProposal squash-merges the proposal with the given number
	// using the given commit message.
	SquashMergeProposalFn() Option[func(number int, message gitdomain.CommitMessage) error]

	// NewProposalURL provides the URL of the page
	// to create a new proposal online.
	NewProposalURL(branch, parentBranch, mainBranch gitdomain.LocalBranchName, proposalTitle gitdomain.ProposalTitle, proposalBody gitdomain.ProposalBody) (string, error)

	// RepositoryURL provides the URL where the current repository can be found online.
	RepositoryURL() string

	// UpdateProposalBase provides a function to update the source branch of proposal.
	UpdateProposalSourceFn() Option[func(number int, newSource gitdomain.LocalBranchName, finalMessages stringslice.Collector) error]

	// UpdateProposalBase provides a function to update the target branch of proposal.
	UpdateProposalTargetFn() Option[func(number int, newTarget gitdomain.LocalBranchName, finalMessages stringslice.Collector) error]
}
