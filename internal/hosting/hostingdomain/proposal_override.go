package hostingdomain

import (
	"os"

	. "github.com/git-town/git-town/v15/pkg/prelude"
)

const OverrideKey = "GIT_TOWN_TEST_PROPOSAL"

func ReadProposalOverride() Option[string] {
	return ParseProposalOverride(os.Getenv(OverrideKey))
}

func ParseProposalOverride(override string) Option[string] {
	if override == "" || override == "(no proposal)" {
		return None[string]()
	}
	return Some(override)
}
