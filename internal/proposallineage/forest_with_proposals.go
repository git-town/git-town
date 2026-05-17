package proposallineage

import (
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

type ForestWithProposals = []TreeNodeWithProposal

func AddProposalsToForest(forest Forest, connector Option[forgedomain.Connector]) ForestWithProposals {
	res := make(ForestWithProposals, 0, len(forest))

	for _, tree := range forest {
		res = append(res, AddProposalsToTree(tree, connector))
	}
	return res
}
