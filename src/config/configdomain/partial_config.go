package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/gohacks"
	"github.com/git-town/git-town/v11/src/messages"
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
		boolValue, err := gohacks.ParseBool(value)
		if err != nil {
			return true, fmt.Errorf(messages.ValueInvalid, KeyOffline, value)
		}
		token := Offline(boolValue)
		self.Offline = &token
	case KeyPerennialBranches:
		if value != "" {
			branches := domain.NewLocalBranchNames(strings.Split(value, " ")...)
			self.PerennialBranches = &branches
		}
	case KeyPushHook:
		parsed, err := gohacks.ParseBool(value)
		if err != nil {
			return true, fmt.Errorf(messages.ValueInvalid, KeyPushHook, value)
		}
		token := PushHook(parsed)
		self.PushHook = &token
	case KeyPushNewBranches:
		parsed, err := gohacks.ParseBool(value)
		if err != nil {
			return true, fmt.Errorf(messages.ValueInvalid, KeyPushNewBranches, value)
		}
		token := NewBranchPush(parsed)
		self.NewBranchPush = &token
	default:
		return false, nil
	}
	return true, nil
}

func EmptyPartialConfig() PartialConfig {
	return PartialConfig{} //nolint:exhaustruct
}

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
