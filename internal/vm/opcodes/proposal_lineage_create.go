package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type ProposalCreateLineage struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ProposalCreateLineage) Run(args shared.RunArgs) error {
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return forgedomain.UnsupportedServiceError()
	}

	builder := configdomain.NewProposalLineageBuilder(connector, args.Config.Value.MainAndPerennials()...)

	// Get Parent Branch for finding the proposal
	parentBranch, hasParentBranch := args.Config.Value.NormalConfig.Lineage.Parent(self.Branch).Get()
	if !hasParentBranch {
		return fmt.Errorf("branch %s has no parent", self.Branch)
	}

	findProposalFn, hasFindProposal := connector.FindProposalFn().Get()
	if !hasFindProposal {
		return fmt.Errorf("connector does not support finding proposals")
	}

	proposal, err := findProposalFn(self.Branch, parentBranch)
	if err != nil {
		return fmt.Errorf("failed to find proposal for branch %s: %w", self.Branch, err)
	}

	proposalData, hasProposal := proposal.Get()
	if !hasProposal {
		return fmt.Errorf("no proposal found branch %q", self.Branch)
	}

}

func (self *ProposalCreateLineage) buildProposalLineage(lineage configdomain.Lineage, currentBranch gitdomain.LocalBranchName, connector forgedomain.Connector) (*proposalLineage, error) {
	ancestors := lineage.Ancestors(currentBranch)
	rootBranch := ancestors[0]
	descendants := lineage.Descendants(currentBranch)

	findProposalFn, hasFindProposal := connector.FindProposalFn().Get()
	if !hasFindProposal {
		return nil, fmt.Errorf("connector does not support finding proposals")
	}

	builder := NewProposalLineageBuilder(rootBranch.BranchName())
	// from root to current branch
	for i := len(ancestors) - 1; i > 0; i-- {
		branch := ancestors[i]
		parent, hasParent := lineage.Parent(branch).Get()
		if hasParent {
			if proposal, err := findProposalFn(branch, parent); err == nil {
				if proposal, hasProposal := proposal.Get(); hasProposal {
					builder.Add(branch, Some(proposal.Data.Data()))
				}
			}
		}
	}

	currentBranchParent, hasCurrentBranchParent := lineage.Parent(currentBranch).Get()

	if hasCurrentBranchParent {
		if proposal, err := findProposalFn(currentBranch, currentBranchParent); err == nil {
			if proposal, hasProposal := proposal.Get(); hasProposal {
				builder.Add(currentBranch, Some(proposal.Data.Data()))

			}
		}
	}
}
