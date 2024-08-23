package hostingdomain

import "os"

const (
	// the key under which the proposal API lookup override gets stored in the environment variables
	OverrideKey = "GIT_TOWN_TEST_PROPOSAL"

	// the content to use in the OverrideKey environment variable to simulate the API returning that no proposal exists
	OverrideNoProposal = "(no proposal)"
)

func ReadProposalOverride() string {
	return os.Getenv(OverrideKey)
}
