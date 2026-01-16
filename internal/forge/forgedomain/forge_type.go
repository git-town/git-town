package forgedomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// ForgeType defines legal values for the "git-town.forge-type" config setting.
type ForgeType string

func (self ForgeType) String() string { return string(self) }

const (
	ForgeTypeAzuredevops         ForgeType = "azuredevops"
	ForgeTypeBitbucket           ForgeType = "bitbucket"
	ForgeTypeBitbucketDatacenter ForgeType = "bitbucket-datacenter"
	ForgeTypeForgejo             ForgeType = "forgejo"
	ForgeTypeGithub              ForgeType = "github"
	ForgeTypeGitlab              ForgeType = "gitlab"
	ForgeTypeGitea               ForgeType = "gitea"
)

// ParseForgeType provides the ForgeType enum matching the given text.
func ParseForgeType(name string, source string) (Option[ForgeType], error) {
	if name == "" {
		return None[ForgeType](), nil
	}
	nameLower := strings.ToLower(name)
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
		ForgeTypeAzuredevops,
		ForgeTypeBitbucket,
		ForgeTypeBitbucketDatacenter,
		ForgeTypeForgejo,
		ForgeTypeGithub,
		ForgeTypeGitlab,
		ForgeTypeGitea,
	}
}
