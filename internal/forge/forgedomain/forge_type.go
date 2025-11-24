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
	ForgeTypeAzureDevOps         ForgeType = "azuredevops"
	ForgeTypeBitbucket           ForgeType = "bitbucket"
	ForgeTypeBitbucketDatacenter ForgeType = "bitbucket-datacenter"
	ForgeTypeForgejo             ForgeType = "forgejo"
	ForgeTypeGitHub              ForgeType = "github"
	ForgeTypeGitLab              ForgeType = "gitlab"
	ForgeTypeGitea               ForgeType = "gitea"
)

func ParseForgeType(name string, source string) (ForgeType, error) {
	forgeTypeOpt, err := ParseForgeTypeOpt(name, source)
	if err != nil {
		return "", err
	}
	if forgeType, hasForgeType := forgeTypeOpt.Get(); hasForgeType {
		return forgeType, nil
	}
	return "", fmt.Errorf(messages.ForgeTypeUnknown, source, name)
}

// ParseForgeTypeOpt provides the ForgeType enum matching the given text.
func ParseForgeTypeOpt(name string, source string) (Option[ForgeType], error) {
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
		ForgeTypeAzureDevOps,
		ForgeTypeBitbucket,
		ForgeTypeBitbucketDatacenter,
		ForgeTypeForgejo,
		ForgeTypeGitHub,
		ForgeTypeGitLab,
		ForgeTypeGitea,
	}
}
