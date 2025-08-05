package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func NewProposalStackLineageBuilder(args *ProposalStackLineageArgs) Option[ProposalStackLineageBuilder] {
	if args == nil || args.Connector.IsNone() {
		return None[ProposalStackLineageBuilder]()
	}

	if args.MainAndPerennialBranches.GetOrDefault().Contains(args.CurrentBranch) {
		// cannot create proposal stack lineage for main or perennial branch
		return None[ProposalStackLineageBuilder]()
	}

	findPropsalFn, hasFindPropsalFn := args.Connector.GetOrPanic().FindProposalFn().Get()
	if !hasFindPropsalFn {
		return None[ProposalStackLineageBuilder]()
	}

	lineage := args.Lineage.BranchLineage(args.CurrentBranch)
	builder := &ProposalStackLineageBuilder{
		mainAndPerennialBranches: args.MainAndPerennialBranches.GetOrDefault(),
		orderedLineage:           make([]*proposalLineage, 0, len(args.Lineage.data)),
	}
	for _, currBranch := range lineage {
		currBranchParent := args.Lineage.Parent(currBranch)
		var err error
		builder, err = builder.addBranch(currBranch, currBranchParent, findPropsalFn)
		if err != nil {
			fmt.Printf("failed to build proposal stack lineage: %s\n", err.Error())
			return None[ProposalStackLineageBuilder]()
		}
	}

	return Some(*builder)
}

type proposalLineage struct {
	branch   gitdomain.LocalBranchName
	proposal Option[forgedomain.ProposalData]
}

type ProposalStackLineageBuilder struct {
	mainAndPerennialBranches gitdomain.LocalBranchNames
	orderedLineage           []*proposalLineage
}

func (self *ProposalStackLineageBuilder) Build(args *ProposalStackLineageArgs) Option[string] {
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

		proposalData, hasProposalData := node.proposal.Get()
		if !hasProposalData {
			break
		}

		builder.WriteString(formattedDisplay(args, indent, proposalData))
	}

	for _, text := range args.AfterStackDisplay {
		builder.WriteString(text)
	}

	return Some(builder.String())
}

func (self *ProposalStackLineageBuilder) GetProposal(branch gitdomain.LocalBranchName) Option[forgedomain.ProposalData] {
	response := None[forgedomain.ProposalData]()
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
) (*ProposalStackLineageBuilder, error) {
	if self.mainAndPerennialBranches.Contains(childBranch) || parentBranch.IsNone() {
		self.orderedLineage = append(self.orderedLineage, &proposalLineage{
			branch:   childBranch,
			proposal: None[forgedomain.ProposalData](),
		})
		return self, nil
	}

	parent := parentBranch.GetOrPanic().BranchName().LocalName()
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

func formattedDisplay(args *ProposalStackLineageArgs, currentIndentLevel string, proposalData forgedomain.ProposalData) string {
	if args.CurrentBranch == proposalData.Source {
		return fmt.Sprintf("%s %s PR %s %s\n", currentIndentLevel, args.IndentMarker, proposalData.URL, args.CurrentBranchIndicator)
	}
	return fmt.Sprintf("%s %s PR %s\n", currentIndentLevel, args.IndentMarker, proposalData.URL)
}
