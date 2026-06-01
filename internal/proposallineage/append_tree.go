package proposallineage

import (
	"strings"

	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

const spacesPerIndent = 2

func AppendTree(builder *strings.Builder, tree TreeNodeWithProposal, currentBranch gitdomain.LocalBranchName, direction configdomain.ProposalBreadcrumbDirection, connector Option[forgedomain.Connector]) {
	switch direction {
	case configdomain.ProposalBreadcrumbDirectionDown:
		renderNodeDown(builder, tree, currentBranch, 0, false, connector)
	case configdomain.ProposalBreadcrumbDirectionUp:
		renderNodeUp(builder, tree, currentBranch, connector)
	}
}

func renderNodeDown(builder *strings.Builder, node TreeNodeWithProposal, currentBranch gitdomain.LocalBranchName, depth int, foundCurrent bool, connector Option[forgedomain.Connector]) {
	if node.BranchOrAncestorHasProposal() || !foundCurrent {
		indent := depth * spacesPerIndent
		builder.WriteString(strings.Repeat(" ", indent))
		builder.WriteString("- ")
		isCurrentBranch := node.Branch == currentBranch && !foundCurrent
		proposal, hasProposal := node.Proposal.Get()
		switch {
		case isCurrentBranch && hasProposal:
			builder.WriteString("**")
			builder.WriteString(proposal.Data.Data().Title.String())
			builder.WriteString("**")
		case hasProposal:
			builder.WriteString(proposalReference(proposal, connector))
		case isCurrentBranch:
			builder.WriteString("**")
			builder.WriteString(node.Branch.String())
			builder.WriteString("**")
		default:
			builder.WriteString(node.Branch.String())
		}
		if isCurrentBranch {
			builder.WriteString(" :point_left:")
			foundCurrent = true
		}
		builder.WriteString("\n")
	}
	for _, child := range node.Children {
		renderNodeDown(builder, child, currentBranch, depth+1, foundCurrent, connector)
	}
}

func renderNodeUp(builder *strings.Builder, node TreeNodeWithProposal, currentBranch gitdomain.LocalBranchName, connector Option[forgedomain.Connector]) bool {
	// First render children (they appear at top in up direction)
	// foundCurrent propagates UP from children to parents
	childFoundCurrent := false
	for _, child := range node.Children {
		if renderNodeUp(builder, child, currentBranch, connector) {
			childFoundCurrent = true
		}
	}
	isCurrentBranch := node.Branch == currentBranch
	// In up direction: if childFoundCurrent or isCurrentBranch, we're on the path from current to root
	onPathToRoot := childFoundCurrent || isCurrentBranch
	// Render if: has proposal/descendant with proposal, OR on path from current to root
	if node.BranchOrAncestorHasProposal() || onPathToRoot {
		builder.WriteString("- ")
		proposal, hasProposal := node.Proposal.Get()
		switch {
		case isCurrentBranch:
			builder.WriteString("**")
			builder.WriteString(node.Branch.String())
			builder.WriteString("**")
		case hasProposal:
			builder.WriteString(proposalReference(proposal, connector))
		default:
			builder.WriteString(node.Branch.String())
		}
		if isCurrentBranch {
			builder.WriteString(" :point_left:")
		}
		builder.WriteString("\n")
	}
	return onPathToRoot
}

func proposalReference(proposal forgedomain.Proposal, connector Option[forgedomain.Connector]) string {
	proposalData := proposal.Data.Data()
	if proposalConnector, hasConnector := connector.Get(); hasConnector {
		if renderedReference := proposalConnector.ProposalReference(proposalData); renderedReference != "" {
			return renderedReference
		}
	}
	return forgedomain.ProposalReferenceFallback(proposalData)
}
