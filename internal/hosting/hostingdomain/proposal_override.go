package hostingdomain

import "os"

const (
	OverrideKey        = "GIT_TOWN_TEST_PROPOSAL"
	OverrideNoProposal = "(no proposal)"
)

func ReadProposalOverride() string {
	return os.Getenv(OverrideKey)
}
