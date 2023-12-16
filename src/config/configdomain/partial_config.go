package configdomain

import (
	"github.com/git-town/git-town/v11/src/domain"
)

// Data contains configuration data as it is stored in a particular configuration data source (Git, config file).
type PartialConfig struct {
	CodeHostingPlatformName *string
	GiteaToken              *GiteaToken
	GitHubToken             *GitHubToken
	GitLabToken             *GitLabToken
	MainBranch              *domain.LocalBranchName
	NewBranchPush           *NewBranchPush
	Offline                 *Offline
	PerennialBranches       *domain.LocalBranchNames
	PushHook                *PushHook
}

func (self *PartialConfig) Add(key Key, value string) (bool, error) {
	var err error
	switch key {
	case KeyCodeHostingPlatform:
		self.CodeHostingPlatformName = &value
	case KeyGiteaToken:
		self.GiteaToken = NewGiteaTokenRef(value)
	case KeyGithubToken:
		self.GitHubToken = NewGitHubTokenRef(value)
	case KeyGitlabToken:
		self.GitLabToken = NewGitLabTokenRef(value)
	case KeyMainBranch:
		self.MainBranch = domain.NewLocalBranchNameRefAllowEmpty(value)
	case KeyOffline:
		self.Offline, err = NewOfflineRef(value)
	case KeyPerennialBranches:
		self.PerennialBranches = domain.NewLocalBranchNamesRef(value)
	case KeyPushHook:
		self.PushHook, err = NewPushHookRef(value)
	case KeyPushNewBranches:
		self.NewBranchPush, err = NewNewBranchPushRef(value)
	default:
		return false, nil
	}
	return true, err
}

func EmptyPartialConfig() PartialConfig {
	return PartialConfig{} //nolint:exhaustruct
}

// PartialConfigDiff diffs the given PartialConfig instances.
func PartialConfigDiff(before, after PartialConfig) ConfigDiff {
	result := ConfigDiff{
		Added:   []Key{},
		Removed: map[Key]string{},
		Changed: map[Key]domain.Change[string]{},
	}
	CheckPtr(&result, KeyGiteaToken, before.GiteaToken, after.GiteaToken)
	CheckPtr(&result, KeyGithubToken, before.GitHubToken, after.GitHubToken)
	CheckPtr(&result, KeyGitlabToken, before.GitLabToken, after.GitLabToken)
	CheckPtr(&result, KeyMainBranch, before.MainBranch, after.MainBranch)
	CheckPtr(&result, KeyOffline, before.Offline, after.Offline)
	CheckPtr(&result, KeyPushHook, before.PushHook, after.PushHook)
	CheckPtr(&result, KeyPushNewBranches, before.NewBranchPush, after.NewBranchPush)
	CheckLocalBranchNames(&result, KeyPerennialBranches, before.PerennialBranches, after.PerennialBranches)
	CheckStringPtr(&result, KeyCodeHostingPlatform, before.CodeHostingPlatformName, after.CodeHostingPlatformName)
	return result
}
