package envconfig

import (
	"os"
	"strings"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func Load(env Environment) configdomain.PartialConfig {
	return configdomain.PartialConfig{
		Aliases:                  configdomain.Aliases{}, // aliases aren't loaded from env vars
		AutoResolve:              None[configdomain.AutoResolve](),
		BitbucketAppPassword:     None[forgedomain.BitbucketAppPassword](),
		BitbucketUsername:        None[forgedomain.BitbucketUsername](),
		BranchTypeOverrides:      configdomain.BranchTypeOverrides{},
		CodebergToken:            None[forgedomain.CodebergToken](),
		ContributionRegex:        None[configdomain.ContributionRegex](),
		DevRemote:                None[gitdomain.Remote](),
		DryRun:                   None[configdomain.DryRun](),
		FeatureRegex:             None[configdomain.FeatureRegex](),
		ForgeType:                None[forgedomain.ForgeType](),
		GitHubConnectorType:      None[forgedomain.GitHubConnectorType](),
		GitHubToken:              forgedomain.ParseGitHubToken(env.Get("GITHUB_TOKEN", "GITHUB_AUTH_TOKEN")),
		GitLabConnectorType:      None[forgedomain.GitLabConnectorType](),
		GitLabToken:              None[forgedomain.GitLabToken](),
		GitUserEmail:             None[gitdomain.GitUserEmail](),
		GitUserName:              None[gitdomain.GitUserName](),
		GiteaToken:               None[forgedomain.GiteaToken](),
		HostingOriginHostname:    None[configdomain.HostingOriginHostname](),
		Lineage:                  configdomain.NewLineage(),
		MainBranch:               None[gitdomain.LocalBranchName](),
		NewBranchType:            None[configdomain.NewBranchType](),
		ObservedRegex:            None[configdomain.ObservedRegex](),
		Offline:                  None[configdomain.Offline](),
		PerennialBranches:        gitdomain.LocalBranchNames{},
		PerennialRegex:           None[configdomain.PerennialRegex](),
		ProposalsShowLineage:     None[forgedomain.ProposalsShowLineage](),
		PushHook:                 None[configdomain.PushHook](),
		ShareNewBranches:         None[configdomain.ShareNewBranches](),
		ShipDeleteTrackingBranch: None[configdomain.ShipDeleteTrackingBranch](),
		ShipStrategy:             None[configdomain.ShipStrategy](),
		SyncFeatureStrategy:      None[configdomain.SyncFeatureStrategy](),
		SyncPerennialStrategy:    None[configdomain.SyncPerennialStrategy](),
		SyncPrototypeStrategy:    None[configdomain.SyncPrototypeStrategy](),
		SyncTags:                 None[configdomain.SyncTags](),
		SyncUpstream:             None[configdomain.SyncUpstream](),
		UnknownBranchType:        None[configdomain.UnknownBranchType](),
		Verbose:                  None[configdomain.Verbose](),
	}
}

func load[T any](key configdomain.Key, parser func(string) (Option[T], error)) (Option[T], error) {
	envName := Key2Env(key)
	value := os.Getenv(envName)
	return parser(value)
}

func Key2Env(key configdomain.Key) string {
	result := strings.ToUpper(key.String())
	result = strings.ReplaceAll(result, ".", "_")
	result = strings.ReplaceAll(result, "-", "_")
	return result
}
