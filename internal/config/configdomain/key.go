package configdomain

import (
	"encoding/json"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/pkg"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// Key contains all the keys used in Git Town's Git metadata configuration.
type Key string //nolint: recvcheck // MarshalJSON and UnmarshalJSON require different receiver types

// ConfigSetting contains a key-value pair for a configuration setting.
type ConfigSetting struct {
	Key   Key
	Value string
}

// ConfigUpdate contains the before and after values of a configuration setting.
type ConfigUpdate struct {
	After  ConfigSetting
	Before ConfigSetting
}

// MarshalJSON is used when serializing this LocalBranchName to JSON.
func (self Key) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.String())
}

func (self Key) String() string { return string(self) }

// UnmarshalJSON is used when de-serializing JSON into a Location.
func (self *Key) UnmarshalJSON(b []byte) error {
	value := ""
	err := json.Unmarshal(b, &value)
	*self = Key(value)
	return err
}

const (
	KeyAliasAppend                         = Key("alias.append")
	KeyAliasCompress                       = Key("alias.compress")
	KeyAliasContribute                     = Key("alias.contribute")
	KeyAliasDiffParent                     = Key("alias.diff-parent")
	KeyAliasHack                           = Key("alias.hack")
	KeyAliasDelete                         = Key("alias.delete")
	KeyAliasObserve                        = Key("alias.observe")
	KeyAliasPark                           = Key("alias.park")
	KeyAliasPrepend                        = Key("alias.prepend")
	KeyAliasPropose                        = Key("alias.propose")
	KeyAliasRename                         = Key("alias.rename")
	KeyAliasRepo                           = Key("alias.repo")
	KeyAliasSetParent                      = Key("alias.set-parent")
	KeyAliasShip                           = Key("alias.ship")
	KeyAliasSync                           = Key("alias.sync")
	KeyAutoResolve                         = Key("git-town.auto-resolve")
	KeyAutoSync                            = Key("git-town.auto-sync")
	KeyBitbucketAppPassword                = Key("git-town.bitbucket-app-password")
	KeyBitbucketUsername                   = Key("git-town.bitbucket-username")
	KeyBranchPrefix                        = Key("git-town.branch-prefix")
	KeyBrowser                             = Key("git-town.browser")
	KeyContributionRegex                   = Key("git-town.contribution-regex")
	KeyDeprecatedCodebergToken             = Key("git-town.codeberg-token")
	KeyDeprecatedCodeHostingDriver         = Key("git-town.code-hosting-driver")
	KeyDeprecatedCodeHostingOriginHostname = Key("git-town.code-hosting-origin-hostname")
	KeyDeprecatedCodeHostingPlatform       = Key("git-town.code-hosting-platform")
	KeyDeprecatedContributionBranches      = Key("git-town.contribution-branches")
	KeyDeprecatedCreatePrototypeBranches   = Key("git-town.create-prototype-branches")
	KeyDeprecatedDefaultBranchType         = Key("git-town.default-branch-type")
	KeyDeprecatedAliasKill                 = Key("alias.kill")
	KeyDeprecatedAliasRenameBranch         = Key("alias.rename-branch")
	KeyDeprecatedHostingPlatform           = Key("git-town.hosting-platform")
	KeyDeprecatedMainBranchName            = Key("git-town.main-branch-name")
	KeyDeprecatedNewBranchPushFlag         = Key("git-town.new-branch-push-flag")
	KeyDeprecatedObservedBranches          = Key("git-town.observed-branches")
	KeyDeprecatedParkedBranches            = Key("git-town.parked-branches")
	KeyDeprecatedPerennialBranchNames      = Key("git-town.perennial-branch-names")
	KeyDeprecatedProposalsShowLineage      = Key("git-town.proposals-show-lineage")
	KeyDeprecatedPrototypeBranches         = Key("git-town.prototype-branches")
	KeyDeprecatedPullBranchStrategy        = Key("git-town.pull-branch-strategy")
	KeyDeprecatedPushNewBranches           = Key("git-town.push-new-branches")
	KeyDeprecatedPushVerify                = Key("git-town.push-verify")
	KeyDeprecatedShipDeleteRemoteBranch    = Key("git-town.ship-delete-remote-branch")
	KeyDeprecatedSyncStrategy              = Key("git-town.sync-strategy")
	KeyDetached                            = Key("git-town.detached")
	KeyDevRemote                           = Key("git-town.dev-remote")
	KeyDisplayTypes                        = Key("git-town.display-types")
	KeyFeatureRegex                        = Key("git-town.feature-regex")
	KeyForgejoToken                        = Key("git-town.forgejo-token")
	KeyForgeType                           = Key("git-town.forge-type")
	KeyGiteaToken                          = Key("git-town.gitea-token")
	KeyGithubConnectorType                 = Key("git-town.github-connector")
	KeyGithubToken                         = Key(pkg.KeyGithubToken)
	KeyGitlabConnectorType                 = Key("git-town.gitlab-connector")
	KeyGitlabToken                         = Key("git-town.gitlab-token")
	KeyHostingOriginHostname               = Key("git-town.hosting-origin-hostname")
	KeyIgnoreUncommitted                   = Key("git-town.ignore-uncommitted")
	KeyMainBranch                          = Key("git-town.main-branch")
	KeyNewBranchType                       = Key("git-town.new-branch-type")
	KeyObservedRegex                       = Key("git-town.observed-regex")
	KeyObsoleteSyncBeforeShip              = Key("git-town.sync-before-ship")
	KeyOffline                             = Key("git-town.offline")
	KeyOrder                               = Key("git-town.order")
	KeyPerennialBranches                   = Key("git-town.perennial-branches")
	KeyPerennialRegex                      = Key("git-town.perennial-regex")
	KeyProposalBreadcrumb                  = Key("git-town.proposal-breadcrumb")
	KeyProposalBreadcrumbDirection         = Key("git-town.proposal-breadcrumb-direction")
	KeyPushBranches                        = Key("git-town.push-branches")
	KeyPushHook                            = Key("git-town.push-hook")
	KeyShareNewBranches                    = Key("git-town.share-new-branches")
	KeyShipDeleteTrackingBranch            = Key("git-town.ship-delete-tracking-branch")
	KeyShipStrategy                        = Key("git-town.ship-strategy")
	KeyStash                               = Key("git-town.stash")
	KeySyncFeatureStrategy                 = Key("git-town.sync-feature-strategy")
	KeySyncPerennialStrategy               = Key("git-town.sync-perennial-strategy")
	KeySyncPrototypeStrategy               = Key("git-town.sync-prototype-strategy")
	KeySyncTags                            = Key("git-town.sync-tags")
	KeySyncUpstream                        = Key("git-town.sync-upstream")
	KeyUnknownBranchType                   = Key("git-town.unknown-branch-type")
	KeyGitUserEmail                        = Key("user.email")
	KeyGitUserName                         = Key("user.name")
)

var keys = []Key{
	KeyAutoResolve,
	KeyAutoSync,
	KeyBitbucketAppPassword,
	KeyBitbucketUsername,
	KeyBranchPrefix,
	KeyBrowser,
	KeyContributionRegex,
	KeyDeprecatedAliasKill,
	KeyDeprecatedAliasRenameBranch,
	KeyDeprecatedCodebergToken,
	KeyDeprecatedCodeHostingDriver,
	KeyDeprecatedCodeHostingOriginHostname,
	KeyDeprecatedCodeHostingPlatform,
	KeyDeprecatedContributionBranches,
	KeyDeprecatedCreatePrototypeBranches,
	KeyDeprecatedDefaultBranchType,
	KeyDeprecatedHostingPlatform,
	KeyDeprecatedMainBranchName,
	KeyDeprecatedNewBranchPushFlag,
	KeyDeprecatedObservedBranches,
	KeyDeprecatedParkedBranches,
	KeyDeprecatedPerennialBranchNames,
	KeyDeprecatedProposalsShowLineage,
	KeyDeprecatedPrototypeBranches,
	KeyDeprecatedPullBranchStrategy,
	KeyDeprecatedPushNewBranches,
	KeyDeprecatedPushVerify,
	KeyDeprecatedShipDeleteRemoteBranch,
	KeyDeprecatedSyncStrategy,
	KeyDetached,
	KeyDevRemote,
	KeyDisplayTypes,
	KeyFeatureRegex,
	KeyForgejoToken,
	KeyForgeType,
	KeyGiteaToken,
	KeyGithubConnectorType,
	KeyGithubToken,
	KeyGitlabConnectorType,
	KeyGitlabToken,
	KeyGitUserEmail,
	KeyGitUserName,
	KeyHostingOriginHostname,
	KeyMainBranch,
	KeyNewBranchType,
	KeyObservedRegex,
	KeyObsoleteSyncBeforeShip,
	KeyOffline,
	KeyOrder,
	KeyPerennialBranches,
	KeyPerennialRegex,
	KeyProposalBreadcrumb,
	KeyProposalBreadcrumbDirection,
	KeyPushBranches,
	KeyPushHook,
	KeyShareNewBranches,
	KeyShipDeleteTrackingBranch,
	KeyIgnoreUncommitted,
	KeyShipStrategy,
	KeyStash,
	KeySyncFeatureStrategy,
	KeySyncPerennialStrategy,
	KeySyncPrototypeStrategy,
	KeySyncTags,
	KeySyncUpstream,
	KeyUnknownBranchType,
}

func NewParentKey(branch gitdomain.LocalBranchName) Key {
	return Key(BranchSpecificKeyPrefix + branch + LineageKeySuffix)
}

func ParseKey(name string) Option[Key] {
	for _, configKey := range keys {
		if configKey.String() == name {
			return Some(configKey)
		}
	}
	if isLineageKey(name) || IsBranchTypeOverrideKey(name) {
		return Some(Key(name))
	}
	if aliasKey, isAliasKey := AllAliasableCommands().LookupKey(name).Get(); isAliasKey {
		return Some(aliasKey.Key())
	}
	return None[Key]()
}

// DeprecatedKeys defines the up-to-date counterparts to deprecated configuration settings.
var DeprecatedKeys = map[Key]Key{
	KeyDeprecatedCodebergToken:             KeyForgejoToken,
	KeyDeprecatedCodeHostingDriver:         KeyForgeType,
	KeyDeprecatedCodeHostingOriginHostname: KeyHostingOriginHostname,
	KeyDeprecatedCodeHostingPlatform:       KeyForgeType,
	KeyDeprecatedHostingPlatform:           KeyForgeType,
	KeyDeprecatedMainBranchName:            KeyMainBranch,
	KeyDeprecatedNewBranchPushFlag:         KeyDeprecatedPushNewBranches,
	KeyDeprecatedPerennialBranchNames:      KeyPerennialBranches,
	KeyDeprecatedProposalsShowLineage:      KeyProposalBreadcrumb,
	KeyDeprecatedPullBranchStrategy:        KeySyncPerennialStrategy,
	KeyDeprecatedPushVerify:                KeyPushHook,
	KeyDeprecatedShipDeleteRemoteBranch:    KeyShipDeleteTrackingBranch,
	KeyDeprecatedSyncStrategy:              KeySyncFeatureStrategy,
	KeyDeprecatedDefaultBranchType:         KeyUnknownBranchType,
}

// ObsoleteKeys defines the keys that are sunset and should get deleted
var ObsoleteKeys = []Key{
	KeyObsoleteSyncBeforeShip,
}

var ObsoleteBranchLists = map[Key]BranchType{
	KeyDeprecatedContributionBranches: BranchTypeContributionBranch,
	KeyDeprecatedObservedBranches:     BranchTypeObservedBranch,
	KeyDeprecatedParkedBranches:       BranchTypeParkedBranch,
	KeyDeprecatedPrototypeBranches:    BranchTypePrototypeBranch,
}

// ConfigUpdates defines the config that should have its keys and values to be updated
var ConfigUpdates = []ConfigUpdate{
	{
		Before: ConfigSetting{
			Key:   KeyProposalBreadcrumb,
			Value: "cli",
		},
		After: ConfigSetting{
			Key:   KeyProposalBreadcrumb,
			Value: string(ProposalBreadcrumbBranches),
		},
	},
	{
		Before: ConfigSetting{
			Key:   KeyDeprecatedAliasKill,
			Value: "town kill",
		},
		After: ConfigSetting{
			Key:   KeyAliasDelete,
			Value: "town delete",
		},
	},
	{
		Before: ConfigSetting{
			Key:   KeyDeprecatedAliasRenameBranch,
			Value: "town rename-branch",
		},
		After: ConfigSetting{
			Key:   KeyAliasRename,
			Value: "town rename",
		},
	},
	{
		Before: ConfigSetting{
			Key:   KeyForgeType,
			Value: "codeberg",
		},
		After: ConfigSetting{
			Key:   KeyForgeType,
			Value: "forgejo",
		},
	},
	{
		Before: ConfigSetting{
			Key:   KeyDeprecatedCreatePrototypeBranches,
			Value: "true",
		},
		After: ConfigSetting{
			Key:   KeyNewBranchType,
			Value: "prototype",
		},
	},
	{
		Before: ConfigSetting{
			Key:   KeyDeprecatedCreatePrototypeBranches,
			Value: "false",
		},
		After: ConfigSetting{
			Key:   KeyNewBranchType,
			Value: "feature",
		},
	},
	{
		Before: ConfigSetting{
			Key:   KeyDeprecatedPushNewBranches,
			Value: "true",
		},
		After: ConfigSetting{
			Key:   KeyShareNewBranches,
			Value: ShareNewBranchesPush.String(),
		},
	},
	{
		Before: ConfigSetting{
			Key:   KeyDeprecatedPushNewBranches,
			Value: "false",
		},
		After: ConfigSetting{
			Key:   KeyShareNewBranches,
			Value: ShareNewBranchesNone.String(),
		},
	},
}
