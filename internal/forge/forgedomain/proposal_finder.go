package forgedomain

import (
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

func ProposalFinderFromConnector(connector Option[Connector]) Option[ProposalFinder] {
	conn, hasConn := connector.Get()
	if !hasConn {
		return None[ProposalFinder]()
	}
	if proposalFinder, hasProposalFinder := conn.(ProposalFinder); hasProposalFinder {
		return Some(proposalFinder)
	}

	return None[ProposalFinder]()
}
