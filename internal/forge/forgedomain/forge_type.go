package forgedomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	"github.com/git-town/git-town/v23/internal/messages"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// ForgeType defines legal values for the "git-town.forge-type" config setting.
type ForgeType string

func (self ForgeType) String() string { return string(self) }

const (
	// keep-sorted start
	ForgeTypeAzuredevops         ForgeType = "azuredevops"
	ForgeTypeBitbucket           ForgeType = "bitbucket"
	ForgeTypeBitbucketDatacenter ForgeType = "bitbucket-datacenter"
	ForgeTypeForgejo             ForgeType = "forgejo"
	ForgeTypeGitea               ForgeType = "gitea"
	ForgeTypeGithub              ForgeType = "github"
	ForgeTypeGitlab              ForgeType = "gitlab"
	// keep-sorted end
)

// ParseForgeType provides the ForgeType enum matching the given text.
func ParseForgeType(name stringss.Trimmed, source string) (Option[ForgeType], error) {
	if name == "" {
		return None[ForgeType](), nil
	}
	nameLower := strings.ToLower(name.String())
	for _, forgeType := range forgeTypes() {
		if nameLower == forgeType.String() {
			return Some(forgeType), nil
		}
	}
	return None[ForgeType](), fmt.Errorf(messages.ForgeTypeUnknown, source, name)
}

// forgeTypes provides all legal values for ForgeType
func forgeTypes() []ForgeType {
	return []ForgeType{
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
