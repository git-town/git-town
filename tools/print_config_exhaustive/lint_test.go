package main_test

import (
	"testing"

	main "github.com/git-town/git-town/tools/print_config_exhaustive"
	"github.com/shoenig/test/must"
)

func TestDefinitionFields(t *testing.T) {
	t.Parallel()

	give := `
	package configdomain

import (
	"slices"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/git-town/git-town/v21/pkg/set"
)

// configuration settings that exist in both UnvalidatedConfig and ValidatedConfig
type NormalConfigData struct {
	Aliases                  Aliases
	BitbucketAppPassword     Option[forgedomain.BitbucketAppPassword]
	BitbucketUsername        Option[forgedomain.BitbucketUsername]
	BranchTypeOverrides      BranchTypeOverrides
	CodebergToken            Option[forgedomain.CodebergToken]
	ContributionRegex        Option[ContributionRegex]
	DevRemote                gitdomain.Remote
	FeatureRegex             Option[FeatureRegex]
	ForgeType                Option[forgedomain.ForgeType] // None = auto-detect
	GitHubConnectorType      Option[forgedomain.GitHubConnectorType]
	GitHubToken              Option[forgedomain.GitHubToken]
	GitLabConnectorType      Option[forgedomain.GitLabConnectorType]
	GitLabToken              Option[forgedomain.GitLabToken]
	GiteaToken               Option[forgedomain.GiteaToken]
	HostingOriginHostname    Option[HostingOriginHostname]
	Lineage                  Lineage
	NewBranchType            Option[BranchType]
	ObservedRegex            Option[ObservedRegex]
	Offline                  Offline
	PerennialBranches        gitdomain.LocalBranchNames
	PerennialRegex           Option[PerennialRegex]
	PushHook                 PushHook
	ShareNewBranches         ShareNewBranches
	ShipDeleteTrackingBranch ShipDeleteTrackingBranch
	ShipStrategy             ShipStrategy
	SyncFeatureStrategy      SyncFeatureStrategy
	SyncPerennialStrategy    SyncPerennialStrategy
	SyncPrototypeStrategy    SyncPrototypeStrategy
	SyncTags                 SyncTags
	SyncUpstream             SyncUpstream
	UnknownBranchType        BranchType
}

func (self *NormalConfigData) NoPushHook() NoPushHook {
	return self.PushHook.Negate()
}
`
	have := main.DefinitionFields(give)
	want := []string{
		"Aliases",
		"BitbucketAppPassword",
		"BitbucketUsername",
		"BranchTypeOverrides",
		"CodebergToken",
		"ContributionRegex",
		"DevRemote",
		"FeatureRegex",
		"ForgeType",
		"GitHubConnectorType",
		"GitHubToken",
		"GitLabConnectorType",
		"GitLabToken",
		"GiteaToken",
		"HostingOriginHostname",
		"Lineage",
		"NewBranchType",
		"ObservedRegex",
		"Offline",
		"PerennialBranches",
		"PerennialRegex",
		"PushHook",
		"ShareNewBranches",
		"ShipDeleteTrackingBranch",
		"ShipStrategy",
		"SyncFeatureStrategy",
		"SyncPerennialStrategy",
		"SyncPrototypeStrategy",
		"SyncTags",
		"SyncUpstream",
		"UnknownBranchType",
	}
	must.Eq(t, want, have)
}
