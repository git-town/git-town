package forgedomain

import (
	"os"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const (
	// OverrideKey defines the environment variable to simulate proposals at the forge during end-to-end testing.
	OverrideKey = "GIT_TOWN_TEST_PROPOSAL"

	// OverrideNoProposal is the value to use in OverrideKey to simulate the absence of proposals at the forge.
	OverrideNoProposal = "(no proposal)"
)

// ProposalOverride allows returning mock proposal data in tests.
type ProposalOverride string

func (self ProposalOverride) String() string {
	return string(self)
}

func ReadProposalOverride() Option[ProposalOverride] {
	return NewOption(ProposalOverride(os.Getenv(OverrideKey)))
}
