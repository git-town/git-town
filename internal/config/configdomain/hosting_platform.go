package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v18/internal/messages"
	. "github.com/git-town/git-town/v18/pkg/prelude"
)

// ForgeType defines legal values for the "git-town.hosting-platform" config setting.
type ForgeType string

func (self ForgeType) String() string { return string(self) }

const (
	ForgeTypeBitbucket           = ForgeType("bitbucket")
	ForgeTypeBitbucketDatacenter = ForgeType("bitbucket-datacenter")
	ForgeTypeGitHub              = ForgeType("github")
	ForgeTypeGitLab              = ForgeType("gitlab")
	ForgeTypeGitea               = ForgeType("gitea")
)

// ParseForgeType provides the HostingPlatform enum matching the given text.
func ParseForgeType(platformName string) (Option[ForgeType], error) {
	if platformName == "" {
		return None[ForgeType](), nil
	}
	platformNameLower := strings.ToLower(platformName)
	for _, forgeType := range forgeTypes() {
		if platformNameLower == forgeType.String() {
			return Some(forgeType), nil
		}
	}
	return None[ForgeType](), fmt.Errorf(messages.ForgeTypeUnknown, platformName)
}

// forgeTypes provides all legal values for HostingPlatform.
func forgeTypes() []ForgeType {
	return []ForgeType{
		ForgeTypeBitbucket,
		ForgeTypeBitbucketDatacenter,
		ForgeTypeGitHub,
		ForgeTypeGitLab,
		ForgeTypeGitea,
	}
}
