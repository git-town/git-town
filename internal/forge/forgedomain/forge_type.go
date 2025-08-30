package forgedomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// ForgeType defines legal values for the "git-town.forge-type" config setting.
type ForgeType string

func (self ForgeType) String() string { return string(self) }

const (
	ForgeTypeBitbucket           ForgeType = "bitbucket"
	ForgeTypeBitbucketDatacenter ForgeType = "bitbucket-datacenter"
	ForgeTypeCodeberg            ForgeType = "codeberg"
	ForgeTypeGitHub              ForgeType = "github"
	ForgeTypeGitLab              ForgeType = "gitlab"
	ForgeTypeGitea               ForgeType = "gitea"
)

// ParseForgeType provides the ForgeType enum matching the given text.
func ParseForgeType(name string) (Option[ForgeType], error) {
	if name == "" {
		return None[ForgeType](), nil
	}
	nameLower := strings.ToLower(name)
	for _, forgeType := range forgeTypes() {
		if nameLower == forgeType.String() {
			return Some(forgeType), nil
		}
	}
	return None[ForgeType](), fmt.Errorf(messages.ForgeTypeUnknown, name)
}

// forgeTypes provides all legal values for ForgeType
func forgeTypes() []ForgeType {
	return []ForgeType{
		ForgeTypeBitbucket,
		ForgeTypeBitbucketDatacenter,
		ForgeTypeCodeberg,
		ForgeTypeGitHub,
		ForgeTypeGitLab,
		ForgeTypeGitea,
	}
}
