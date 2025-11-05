package forgedomain

import (
	"os"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const (
	// OverrideKey is the key under which the proposal API lookup override gets stored in the environment variables.
	OverrideKey = "GIT_TOWN_TEST_PROPOSAL"

	// OverrideNoProposal is the content to use in the OverrideKey environment variable to simulate the API returning that no proposal exists.
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
