// Package config provides facilities to read and write the Git Town configuration.
// Git Town stores its configuration in the Git configuration under the prefix "git-town".
// It supports both the Git configuration for the local repository as well as the global Git configuration in `~/.gitconfig`.
// You can manually read the Git configuration entries for Git Town by running `git config --get-regexp git-town`.
package config

const (
	KeyCodeHostingDriver           = "git-town.code-hosting-driver"
	KeyCodeHostingOriginHostname   = "git-town.code-hosting-origin-hostname"
	KeyDeprecatedNewBranchPushFlag = "git-town.new-branch-push-flag"
	KeyDeprecatedPushVerify        = "git-town.push-verify"
	KeyGiteaToken                  = "git-town.gitea-token"  //nolint:gosec
	KeyGithubToken                 = "git-town.github-token" //nolint:gosec
	KeyGitlabToken                 = "git-town.gitlab-token" //nolint:gosec
	KeyMainBranch                  = "git-town.main-branch-name"
	KeyOffline                     = "git-town.offline"
	KeyPerennialBranches           = "git-town.perennial-branch-names"
	KeyPullBranchStrategy          = "git-town.pull-branch-strategy"
	KeyPushHook                    = "git-town.push-hook"
	KeyPushNewBranches             = "git-town.push-new-branches"
	KeyShipDeleteRemoteBranch      = "git-town.ship-delete-remote-branch"
	KeySyncUpstream                = "git-town.sync-upstream"
	KeySyncStrategy                = "git-town.sync-strategy"
	KeyTestingRemoteURL            = "git-town.testing.remote-url"
)
