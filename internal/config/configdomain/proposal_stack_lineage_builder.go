package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func NewProposalStackLineageBuilder(args ProposalStackLineageArgs) Option[ProposalStackLineageBuilder] {
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return None[ProposalStackLineageBuilder]()
	}

	if args.MainAndPerennialBranches.Contains(args.CurrentBranch) {
		// cannot create proposal stack lineage for main or perennial branch
		return None[ProposalStackLineageBuilder]()
	}

	findPropsalFn, hasFindPropsalFn := connector.FindProposalFn().Get()
	if !hasFindPropsalFn {
		return None[ProposalStackLineageBuilder]()
	}

	lineage := args.Lineage.BranchLineage(args.CurrentBranch)
	builder := &ProposalStackLineageBuilder{
		mainAndPerennialBranches: args.MainAndPerennialBranches,
		orderedLineage:           make([]*proposalLineage, 0, len(args.Lineage.data)),
	}
	for _, currBranch := range lineage {
		currBranchParent := args.Lineage.Parent(currBranch)
		if err := builder.addBranch(currBranch, currBranchParent, findPropsalFn); err != nil {
			fmt.Printf("failed to build proposal stack lineage: %s\n", err.Error())
			return None[ProposalStackLineageBuilder]()
		}
	}

	return Some(*builder)
}

type proposalLineage struct {
	branch   gitdomain.LocalBranchName
	proposal Option[forgedomain.Proposal]
}

type ProposalStackLineageBuilder struct {
	mainAndPerennialBranches gitdomain.LocalBranchNames
	orderedLineage           []*proposalLineage
}

func (self *ProposalStackLineageBuilder) Build(args ProposalStackLineageArgs) Option[string] {
	var builder strings.Builder
	for _, text := range args.BeforeStackDisplay {
		builder.WriteString(text)
	}

	length := len(self.orderedLineage)
	for i := len(self.orderedLineage); i > 0; i-- {
		node := self.orderedLineage[length-i]
		indent := strings.Repeat(" ", (length-i)*2)
		if self.mainAndPerennialBranches.Contains(node.branch) {
			builder.WriteString(fmt.Sprintf("%s %s %s\n", indent, args.IndentMarker, node.branch.BranchName()))
			continue
		}

		proposal, hasProposal := node.proposal.Get()
		if !hasProposal {
			break
		}

		builder.WriteString(formattedDisplay(args, indent, proposal))
	}

	for _, text := range args.AfterStackDisplay {
		builder.WriteString(text)
	}

	return Some(builder.String())
}

func (self *ProposalStackLineageBuilder) GetProposal(branch gitdomain.LocalBranchName) Option[forgedomain.Proposal] {
	response := None[forgedomain.Proposal]()
	for _, curr := range self.orderedLineage {
		if curr.branch == branch {
			response = curr.proposal
		}
	}
	return response
}

func (self *ProposalStackLineageBuilder) addBranch(
	childBranch gitdomain.LocalBranchName,
	parentBranch Option[gitdomain.LocalBranchName],
	findProposalFn func(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error),
) error {
	parent, hasParentBranch := parentBranch.Get()
	if self.mainAndPerennialBranches.Contains(childBranch) || !hasParentBranch {
		self.orderedLineage = append(self.orderedLineage, &proposalLineage{
			branch:   childBranch,
			proposal: None[forgedomain.Proposal](),
		})
		return nil
	}

	proposal, err := findProposalFn(childBranch, parent)
	if err != nil {
		return fmt.Errorf("failed to find proposal for branch %s: %w", childBranch, err)
	}

	proposalData, hasProposal := proposal.Get()
	if !hasProposal {
		self.orderedLineage = append(self.orderedLineage, &proposalLineage{
			branch:   childBranch,
			proposal: None[forgedomain.Proposal](),
		})
		return nil
	}

	self.orderedLineage = append(self.orderedLineage, &proposalLineage{
		branch:   childBranch,
		proposal: Some(proposalData),
	})
	return nil
}

func formattedDisplay(args ProposalStackLineageArgs, currentIndentLevel string, proposal forgedomain.Proposal) string {
	proposalData := proposal.Data
	if args.CurrentBranch == proposalData.Data().Source {
		return fmt.Sprintf("%s %s PR %s %s\n", currentIndentLevel, args.IndentMarker, proposalData.Data().URL, args.CurrentBranchIndicator)
	}
	return fmt.Sprintf("%s %s PR %s\n", currentIndentLevel, args.IndentMarker, proposalData.Data().URL)
}
