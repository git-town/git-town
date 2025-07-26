package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v21/internal/cli/colors"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type ProposalLineageBuilder interface {
	// Adds the next branch in the lineage chain
	AddBranch(childBranch gitdomain.LocalBranchName, parentBranch Option[gitdomain.LocalBranchName]) (ProposalLineageBuilder, error)
	// Build - creates the proposal lineage based on the display location
	Build(currentBranch gitdomain.LocalBranchName, location ProposalLineageIn) Option[string]
}

func NewProposalLineageBuilder(connector forgedomain.Connector, exemptBranches ...gitdomain.LocalBranchName) ProposalLineageBuilder {
	if _, hasFindProposal := connector.FindProposalFn().Get(); !hasFindProposal {
		return &noopProposalLineageBuilder{}
	}

	return &proposalLineageBuilder{
		orderedLineage:                           make([]*proposalLineage, 0),
		connector:                                connector,
		branchesExemptFromDisplayingProposalInfo: exemptBranches,
	}
}

type proposalLineage struct {
	branch   gitdomain.LocalBranchName
	proposal Option[forgedomain.ProposalData]
}

type proposalLineageBuilder struct {
	connector                                forgedomain.Connector
	orderedLineage                           []*proposalLineage
	branchesExemptFromDisplayingProposalInfo gitdomain.LocalBranchNames
}

func (self *proposalLineageBuilder) AddBranch(childBranch gitdomain.LocalBranchName, parentBranch Option[gitdomain.LocalBranchName]) (ProposalLineageBuilder, error) {
	if self.branchesExemptFromDisplayingProposalInfo.Contains(childBranch) || parentBranch.IsNone() {
		self.orderedLineage = append(self.orderedLineage, &proposalLineage{
			branch:   childBranch,
			proposal: None[forgedomain.ProposalData](),
		})
		return self, nil
	}

	parent := parentBranch.GetOrPanic().BranchName().LocalName()
	findProposalFn, _ := self.connector.FindProposalFn().Get()

	proposal, err := findProposalFn(childBranch, parent)
	if err != nil {
		return self, fmt.Errorf("failed to find proposal for branch %s: %w", childBranch, err)
	}

	proposalData, hasProposal := proposal.Get()
	if !hasProposal {
		return self, fmt.Errorf("no proposal found branch %q", childBranch)
	}

	self.orderedLineage = append(self.orderedLineage, &proposalLineage{
		branch:   childBranch,
		proposal: Some(proposalData.Data.Data()),
	})
	return self, nil
}

func (self *proposalLineageBuilder) Build(currentBranch gitdomain.LocalBranchName, location ProposalLineageIn) Option[string] {
	var builder strings.Builder
	builder.WriteString("### This proposal is part of stack\n\n")
	length := len(self.orderedLineage)
	for i := len(self.orderedLineage) - 1; i >= 0; i-- {
		node := self.orderedLineage[length]
		indent := strings.Repeat(" ", length-i)
		if self.branchesExemptFromDisplayingProposalInfo.Contains(node.branch) {
			builder.WriteString(fmt.Sprintf("%s #%s\n", indent, node.branch.BranchName()))
			continue
		}

		proposalData, hasProposalData := node.proposal.Get()
		if !hasProposalData {
			break
		}

		builder.WriteString(locationBasedFormatting(currentBranch, indent, location, proposalData))
	}

	if builder.Len() == 0 {
		return None[string]()
	}

	return Some(builder.String())
}

func locationBasedFormatting(currentBranch gitdomain.LocalBranchName, indent string, location ProposalLineageIn, proposalData forgedomain.ProposalData) string {
	if currentBranch == proposalData.Source {
		if location == ProposalLineageInTerminal {
			return colors.Green().Styled(fmt.Sprintf("%s%s #%d [%s](%s)\n", currentBranchProposalExpression(location), indent, proposalData.Number, proposalData.Title, proposalData.URL))
		} else {
			return fmt.Sprintf("%s #%d [%s](%s) %s\n", indent, proposalData.Number, proposalData.Title, proposalData.URL, currentBranchProposalExpression(location))
		}
	}

	return fmt.Sprintf("%s #%d [%s](%s)\n", indent, proposalData.Number, proposalData.Title, proposalData.URL)
}

func currentBranchProposalExpression(location ProposalLineageIn) string {
	response := ":point_left:"
	if location == ProposalLineageInTerminal {
		response = "*"
	}

	return response
}

type noopProposalLineageBuilder struct{}

func (self *noopProposalLineageBuilder) AddBranch(childBranch gitdomain.LocalBranchName, parentBranch Option[gitdomain.LocalBranchName]) (ProposalLineageBuilder, error) {
	return self, nil
}

func (self *noopProposalLineageBuilder) Build(currentBranch gitdomain.LocalBranchName, location ProposalLineageIn) Option[string] {
	return None[string]()
}
