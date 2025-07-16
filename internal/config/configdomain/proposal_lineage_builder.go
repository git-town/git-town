package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type ProposalLineageBuilder interface {
	// Adds the next branch in the lineage chain
	AddBranch(childBranch gitdomain.LocalBranchName, parentBranch Option[gitdomain.LocalBranchName]) (ProposalLineageBuilder, error)
	// Build - creates the proposal lineage based on the display location
	Build(currentBranch gitdomain.LocalBranchName, location LineageDisplayLocation) Option[string]
}

func NewProposalLineageBuilder(connector forgedomain.Connector, exemptBranches ...gitdomain.LocalBranchName) ProposalLineageBuilder {
	if _, hasFindProposal := connector.FindProposalFn().Get(); !hasFindProposal {
		return &noopProposalLineageBuilder{}
	}

	return &proposalLineageBuilder{
		orderedLineage: make([]forgedomain.ProposalData, 0),
		connector:      connector,
		trackingExempt: exemptBranches,
	}
}

type proposalLineage struct {
	branch   gitdomain.LocalBranchName
	proposal Option[forgedomain.ProposalData]
}

type proposalLineageBuilder struct {
	// map of local branch to remote branch
	orderedLineage []*proposalLineage
	connector      forgedomain.Connector
	trackingExempt gitdomain.LocalBranchNames
}

func (self *proposalLineageBuilder) AddBranch(childBranch gitdomain.LocalBranchName, parentBranch Option[gitdomain.LocalBranchName]) (ProposalLineageBuilder, error) {
	if self.trackingExempt.Contains(childBranch) || parentBranch.IsNone() {
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

func (self *proposalLineageBuilder) Build(currentBranch gitdomain.LocalBranchName, location LineageDisplayLocation) Option[string] {

	var builder strings.Builder
	builder.WriteString("### This proposal is part of stack\n\n")
	length := len(self.orderedLineage)
	currentBranchExpression := currentBranchProposalExpression(location)
	for i := len(self.orderedLineage) - 1; i >= 0; i-- {
		proposalData := self.orderedLineage[length]
		if self.trackingExempt.Contains(proposalData.Source) {
			continue
		}
		indent := strings.Repeat(" ", length-i)
		if currentBranch == proposalData.Source {
			builder.WriteString(fmt.Sprintf("%s↳ #%d [%s](%s) %s\n", indent, proposalData.Number, proposalData.Title, proposalData.URL, currentBranchExpression))
		} else {
			builder.WriteString(fmt.Sprintf("%s↳ #%d [%s](%s)\n", indent, proposalData.Number, proposalData.Title, proposalData.URL))
		}
	}

	return Some(builder.String())
}

func currentBranchProposalExpression(location LineageDisplayLocation) string {
	response := ":point_left:"
	if location == LineageDisplayLocationTerminal {
		response = "☜"
	}

	return response
}

type noopProposalLineageBuilder struct {
}

func (self *noopProposalLineageBuilder) AddBranch(childBranch gitdomain.LocalBranchName, parentBranch Option[gitdomain.LocalBranchName]) (ProposalLineageBuilder, error) {
	return self, nil
}

func (self *noopProposalLineageBuilder) Build(currentBranch gitdomain.LocalBranchName, location LineageDisplayLocation) Option[string] {
	return None[string]()
}
