package configdomain

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type proposalStackLineageBuildOptions struct {
	afterStackDisplay      []string
	beforeStackDisplay     []string
	currentBranch          Option[gitdomain.LocalBranchName]
	currentBranchIndicator string
	indentMarker           string
	location               ProposalLineageIn
}

type configureProposalStackLineage func(opts *proposalStackLineageBuildOptions)

// WithCurrentBranch
// Informs the builder which branch is the current. This is used to determine when the
// `WithCurrentBranchIndicator` character appears.
func WithCurrentBranch(branch gitdomain.LocalBranchName) configureProposalStackLineage {
	return func(opts *proposalStackLineageBuildOptions) {
		opts.currentBranch = Some(branch)
	}
}

// WithCurrentBranchIndicator
// Special character used to denote the current branch's proposal (if there is one).
func WithCurrentBranchIndicator(indicator string) configureProposalStackLineage {
	return func(opts *proposalStackLineageBuildOptions) {
		opts.currentBranchIndicator = indicator
	}
}

// WithIndentMarker
// Controls the marker following an indent.
func WithIndentMarker(marker string) configureProposalStackLineage {
	return func(opts *proposalStackLineageBuildOptions) {
		opts.indentMarker = marker
	}
}

// WithProposalLineageIn
// Determines the context the proposal stack lineage is displayed.
func WithProposalLineageIn(location ProposalLineageIn) configureProposalStackLineage {
	return func(opts *proposalStackLineageBuildOptions) {
		opts.location = location
	}
}

// WithStringAfterStackDisplay
// A set of texts to appear after the main stack information is displayed.
// If this method is called more than once, the text appear in FIFO order.
func WithStringAfterStackDisplay(text string) configureProposalStackLineage {
	return func(opts *proposalStackLineageBuildOptions) {
		opts.afterStackDisplay = append(opts.afterStackDisplay, text)
	}
}

// WithStringBeforeStackDisplay
// A set of texts to appear before the main stack information is displayed.
// If this method is called more than once, the text appear in FIFO order.
func WithStringBeforeStackDisplay(text string) configureProposalStackLineage {
	return func(opts *proposalStackLineageBuildOptions) {
		opts.beforeStackDisplay = append(opts.beforeStackDisplay, text)
	}
}

func newProposalStackLineageBuilderOptions() *proposalStackLineageBuildOptions {
	return &proposalStackLineageBuildOptions{
		afterStackDisplay:      make([]string, 0),
		beforeStackDisplay:     make([]string, 0),
		currentBranch:          None[gitdomain.LocalBranchName](),
		currentBranchIndicator: "point_left",
		indentMarker:           "-",
		location:               ProposalLineageOperationInProposalBody,
	}
}
