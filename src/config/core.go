// Package config provides facilities to read and write the Git Town configuration.
// Git Town stores its configuration in the Git configuration under the prefix "git-town".
// It supports both the Git configuration for the local repository as well as the global Git configuration in `~/.gitconfig`.
// You can manually read the Git configuration entries for Git Town by running `git config --get-regexp git-town`.
package config

const (
	CodeHostingDriverKey         = "git-town.code-hosting-driver"
	CodeHostingOriginHostnameKey = "git-town.code-hosting-origin-hostname"
	GiteaTokenKey                = "git-town.gitea-token"  //nolint:gosec
	GithubTokenKey               = "git-town.github-token" //nolint:gosec
	GitlabTokenKey               = "git-town.gitlab-token" //nolint:gosec
	MainBranchKey                = "git-town.main-branch-name"
	NewBranchPushFlagKey         = "git-town.new-branch-push-flag"
	OfflineKey                   = "git-town.offline"
	PerennialBranchesKey         = "git-town.perennial-branch-names"
	PullBranchStrategyKey        = "git-town.pull-branch-strategy"
	PushHookKey                  = "git-town.push-hook"
	DeprecatedPushVerifyKey      = "git-town.push-verify"
	PushNewBranchesKey           = "git-town.push-new-branches"
	ShipDeleteRemoteBranchKey    = "git-town.ship-delete-remote-branch"
	SyncUpstreamKey              = "git-town.sync-upstream"
	SyncStrategyKey              = "git-town.sync-strategy"
	TestingRemoteURLKey          = "git-town.testing.remote-url"
)
