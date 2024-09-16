package hostingdomain

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/gohacks/stringslice"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// Connector describes the activities that Git Town can perform on code hosting platforms.
// Individual implementations exist to talk to specific hosting platforms.
type Connector interface {
	// CanMakeAPICalls indicates whether this connector instance is configured to make API calls.
	// Most connectors need to have an API key to do so.
	CanMakeAPICalls() bool

	// DefaultProposalMessage provides the text that the form for creating new proposals
	// on the respective hosting platform is prepopulated with.
	DefaultProposalMessage(proposal Proposal) string

	// FindProposal provides details about the proposal for the given branch into the given target branch.
	// Returns nil if no proposal exists.
	FindProposal(branch, target gitdomain.LocalBranchName) (Option[Proposal], error)

	// SearchProposals searches for a proposal that has the given branch as its head (source) branch.
	SearchProposals(branch gitdomain.LocalBranchName) (Option[Proposal], error)

	// SquashMergeProposal squash-merges the proposal with the given number
	// using the given commit message.
	SquashMergeProposal(number int, message gitdomain.CommitMessage) error

	// NewProposalURL provides the URL of the page
	// to create a new proposal online.
	NewProposalURL(branch, parentBranch, mainBranch gitdomain.LocalBranchName, proposalTitle gitdomain.ProposalTitle, proposalBody gitdomain.ProposalBody) (string, error)

	// RepositoryURL provides the URL where the current repository can be found online.
	RepositoryURL() string

	// UpdateProposalBase updates the target branch of the given proposal.
	UpdateProposalBase(number int, target gitdomain.LocalBranchName, finalMessages stringslice.Collector) error

	// UpdateProposalBase updates the target branch of the given proposal.
	UpdateProposalHead(number int, target gitdomain.LocalBranchName, finalMessages stringslice.Collector) error
}
