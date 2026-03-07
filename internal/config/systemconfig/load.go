package systemconfig

import (
	"errors"
	"os"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/mattn/go-isatty"
)

func Load() configdomain.PartialConfig {
	tty := hasTTY()
	displayDialogs := configdomain.DisplayDialogs(tty)
	return configdomain.PartialConfig{
		Aliases:                     configdomain.Aliases{},
		AutoResolve:                 None[configdomain.AutoResolve](),
		AutoSync:                    None[configdomain.AutoSync](),
		BitbucketAppPassword:        None[forgedomain.BitbucketAppPassword](),
		BitbucketUsername:           None[forgedomain.BitbucketUsername](),
		BranchPrefix:                None[configdomain.BranchPrefix](),
		BranchTypeOverrides:         configdomain.BranchTypeOverrides{},
		Browser:                     None[configdomain.Browser](),
		ContributionRegex:           None[configdomain.ContributionRegex](),
		Detached:                    None[configdomain.Detached](),
		DevRemote:                   None[gitdomain.Remote](),
		DisplayDialogs:              Some(displayDialogs),
		DisplayTypes:                None[configdomain.DisplayTypes](),
		DryRun:                      None[configdomain.DryRun](),
		FeatureRegex:                None[configdomain.FeatureRegex](),
		ForgeType:                   None[forgedomain.ForgeType](),
		ForgejoToken:                None[forgedomain.ForgejoToken](),
		GitUserEmail:                None[gitdomain.GitUserEmail](),
		GitUserName:                 None[gitdomain.GitUserName](),
		GiteaToken:                  None[forgedomain.GiteaToken](),
		GithubConnectorType:         None[forgedomain.GithubConnectorType](),
		GithubToken:                 None[forgedomain.GithubToken](),
		GitlabConnectorType:         None[forgedomain.GitlabConnectorType](),
		GitlabToken:                 None[forgedomain.GitlabToken](),
		HostingOriginHostname:       None[configdomain.HostingOriginHostname](),
		IgnoreUncommitted:           None[configdomain.IgnoreUncommitted](),
		Lineage:                     configdomain.NewLineage(),
		MainBranch:                  None[gitdomain.LocalBranchName](),
		NewBranchType:               None[configdomain.NewBranchType](),
		ObservedRegex:               None[configdomain.ObservedRegex](),
		Offline:                     None[configdomain.Offline](),
		Order:                       None[configdomain.Order](),
		PerennialBranches:           gitdomain.LocalBranchNames{},
		PerennialRegex:              None[configdomain.PerennialRegex](),
		ProposalBreadcrumb:          None[configdomain.ProposalBreadcrumb](),
		ProposalBreadcrumbDirection: None[configdomain.ProposalBreadcrumbDirection](),
		PushBranches:                None[configdomain.PushBranches](),
		PushHook:                    None[configdomain.PushHook](),
		ShareNewBranches:            None[configdomain.ShareNewBranches](),
		ShipDeleteTrackingBranch:    None[configdomain.ShipDeleteTrackingBranch](),
		ShipStrategy:                None[configdomain.ShipStrategy](),
		Stash:                       None[configdomain.Stash](),
		SyncFeatureStrategy:         None[configdomain.SyncFeatureStrategy](),
		SyncPerennialStrategy:       None[configdomain.SyncPerennialStrategy](),
		SyncPrototypeStrategy:       None[configdomain.SyncPrototypeStrategy](),
		SyncTags:                    None[configdomain.SyncTags](),
		SyncUpstream:                None[configdomain.SyncUpstream](),
		UnknownBranchType:           None[configdomain.UnknownBranchType](),
		Verbose:                     None[configdomain.Verbose](),
	}
}

// hasTTY reports whether an interactive terminal is available.
func hasTTY() bool {
	fd := os.Stdin.Fd()
	if isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd) {
		return true
	}
	return canOpenTTY()
}

// ErrNoTTY indicates that an interactive terminal is required but not available.
var ErrNoTTY = errors.New("no interactive terminal available")
