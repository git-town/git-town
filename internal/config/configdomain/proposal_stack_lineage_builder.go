package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v21/internal/cli/colors"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type ProposalStackLineageBuilder interface {
	// Adds the next branch in the lineage chain
	AddBranch(childBranch gitdomain.LocalBranchName, parentBranch Option[gitdomain.LocalBranchName]) (ProposalStackLineageBuilder, error)
	// Build - creates the proposal lineage based on the display location
	Build(cfgs ...configureProposalStackLineage) Option[string]
	// GetProposal fetches the proposal data for a branch, if there is one.
	GetProposal(branch gitdomain.LocalBranchName) Option[forgedomain.ProposalData]
}

// NewProposalStackLineageBuilder enables generating the proposal stack lineage under daferent contexts
// connector - forgedomain.Connector
// exemptBranches - the branches we do not care to fetch proposal data.
func NewProposalStackLineageBuilder(connector forgedomain.Connector, exemptBranches ...gitdomain.LocalBranchName) ProposalStackLineageBuilder {
	if _, hasFindProposal := connector.FindProposalFn().Get(); !hasFindProposal {
		// If there is no way to find proposals, use a no-op builder
		return &noopProposalStackLineageBuilder{}
	}

	return &proposalStackLineageBuilder{
		branchesExemptFromDisplayingProposalInfo: exemptBranches,
		connector:                                connector,
		orderedLineage:                           make([]*proposalLineage, 0),
	}
}

type proposalLineage struct {
	branch   gitdomain.LocalBranchName
	proposal Option[forgedomain.ProposalData]
}

type proposalStackLineageBuilder struct {
	branchesExemptFromDisplayingProposalInfo gitdomain.LocalBranchNames
	connector                                forgedomain.Connector
	orderedLineage                           []*proposalLineage
}

func (self *proposalStackLineageBuilder) AddBranch(childBranch gitdomain.LocalBranchName, parentBranch Option[gitdomain.LocalBranchName]) (ProposalStackLineageBuilder, error) {
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

func (self *proposalStackLineageBuilder) Build(cfgs ...configureProposalStackLineage) Option[string] {
	builderOptions := newProposalStackLineageBuilderOptions()
	for _, cfg := range cfgs {
		cfg(builderOptions)
	}

	var builder strings.Builder
	for _, text := range builderOptions.beforeStackDisplay {
		builder.WriteString(text)
	}

	length := len(self.orderedLineage)
	for i := len(self.orderedLineage); i > 0; i-- {
		node := self.orderedLineage[length-i]
		indent := strings.Repeat(" ", (length-i)*2)
		if self.branchesExemptFromDisplayingProposalInfo.Contains(node.branch) {
			builder.WriteString(fmt.Sprintf("%s %s %s\n", indent, builderOptions.indentMarker, node.branch.BranchName()))
			continue
		}

		proposalData, hasProposalData := node.proposal.Get()
		if !hasProposalData {
			break
		}

		builder.WriteString(formattedDisplay(builderOptions, indent, proposalData))
	}

	for _, text := range builderOptions.afterStackDisplay {
		builder.WriteString(text)
	}

	return Some(builder.String())
}

func (self *proposalStackLineageBuilder) GetProposal(branch gitdomain.LocalBranchName) Option[forgedomain.ProposalData] {
	response := None[forgedomain.ProposalData]()
	for _, curr := range self.orderedLineage {
		if curr.branch == branch {
			response = curr.proposal
		}
	}
	return response
}

func formattedDisplay(builderOptions *proposalStackLineageBuildOptions, currentIndentLevel string, proposalData forgedomain.ProposalData) string {
	if builderOptions.location == ProposalLineageInTerminal {
		if builderOptions.currentBranch.GetOrDefault() == proposalData.Source {
			return colors.Green().Styled(fmt.Sprintf("%s%s %s PR #%d %s (%s)\n", builderOptions.currentBranchIndicator, currentIndentLevel, builderOptions.indentMarker, proposalData.Number, proposalData.Title, proposalData.URL))
		}
		return fmt.Sprintf("%s %s PR #%d %s (%s)\n", currentIndentLevel, builderOptions.indentMarker, proposalData.Number, proposalData.Title, proposalData.URL)
	} else {
		if builderOptions.currentBranch.GetOrDefault() == proposalData.Source {
			return fmt.Sprintf("%s %s PR %s %s\n", currentIndentLevel, builderOptions.indentMarker, proposalData.URL, builderOptions.currentBranchIndicator)
		}
		return fmt.Sprintf("%s %s PR %s\n", currentIndentLevel, builderOptions.indentMarker, proposalData.URL)
	}
}
