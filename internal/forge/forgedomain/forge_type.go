package forgedomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	"github.com/git-town/git-town/v23/internal/messages"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// ConfiguredForgeType defines legal values for the "git-town.forge-type" config setting.
type ConfiguredForgeType string

func (self ConfiguredForgeType) String() string { return string(self) }

const (
	// keep-sorted start
	ForgeTypeAzuredevops         ConfiguredForgeType = "azuredevops"
	ForgeTypeBitbucket           ConfiguredForgeType = "bitbucket"
	ForgeTypeBitbucketDatacenter ConfiguredForgeType = "bitbucket-datacenter"
	ForgeTypeForgejo             ConfiguredForgeType = "forgejo"
	ForgeTypeGitea               ConfiguredForgeType = "gitea"
	ForgeTypeGithub              ConfiguredForgeType = "github"
	ForgeTypeGitlab              ConfiguredForgeType = "gitlab"
	// keep-sorted end
)

// ParseForgeType provides the ForgeType enum matching the given text.
func ParseForgeType(name stringss.Trimmed, source string) (Option[ConfiguredForgeType], error) {
	if name == "" {
		return None[ConfiguredForgeType](), nil
	}
	nameLower := strings.ToLower(name.String())
	for _, forgeType := range forgeTypes() {
		if nameLower == forgeType.String() {
			return Some(forgeType), nil
		}
	}
	return None[ConfiguredForgeType](), fmt.Errorf(messages.ForgeTypeUnknown, source, name)
}

// forgeTypes provides all legal values for ForgeType
func forgeTypes() []ConfiguredForgeType {
	return []ConfiguredForgeType{
		// keep-sorted start
		ForgeTypeAzuredevops,
		ForgeTypeBitbucket,
		ForgeTypeBitbucketDatacenter,
		ForgeTypeForgejo,
		ForgeTypeGitea,
		ForgeTypeGithub,
		ForgeTypeGitlab,
		// keep-sorted end
	}
}

type DetectedForgeType ConfiguredForgeType
