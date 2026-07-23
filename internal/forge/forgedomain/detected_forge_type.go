package forgedomain

import . "github.com/git-town/git-town/v24/pkg/prelude"

// DetectedForgeType is the forge type that actually exists.
// Its the one the user has configured,
// or if that is "auto", the automatically detected forge type.
type DetectedForgeType ForgeType

// ForgeType converts this value into a ForgeType.
func (self DetectedForgeType) ForgeType() ForgeType { return ForgeType(self) }

func IsBitbucket(forgeTypeOpt Option[DetectedForgeType]) bool {
	if detectedForgeType, hasForgeType := forgeTypeOpt.Get(); hasForgeType {
		switch detectedForgeType.ForgeType() {
		case ForgeTypeBitbucket, ForgeTypeBitbucketDatacenter:
			return true
		case ForgeTypeAzuredevops, ForgeTypeForgejo, ForgeTypeGitea, ForgeTypeGithub, ForgeTypeGitlab:
			return false
		}
	}
	return false
}
