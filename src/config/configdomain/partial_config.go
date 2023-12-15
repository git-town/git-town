package configdomain

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
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
	Offline                 *Offline
	PerennialBranches       *domain.LocalBranchNames
}

func (self *PartialConfig) Add(key Key, value string) (bool, error) {
	switch key {
	case KeyCodeHostingPlatform:
		self.CodeHostingPlatformName = &value
	case KeyGiteaToken:
		token := GiteaToken(value)
		self.GiteaToken = &token
	case KeyGithubToken:
		token := GitHubToken(value)
		self.GitHubToken = &token
	case KeyGitlabToken:
		token := GitLabToken(value)
		self.GitLabToken = &token
	case KeyMainBranch:
		var token domain.LocalBranchName
		if value == "" {
			token = domain.EmptyLocalBranchName()
		} else {
			token = domain.NewLocalBranchName(value)
		}
		self.MainBranch = &token
	case KeyOffline:
		boolValue, err := gohacks.ParseBool(value)
		if err != nil {
			return false, fmt.Errorf(messages.ValueInvalid, KeyOffline, value)
		}
		token := Offline(boolValue)
		self.Offline = &token
	case KeyPerennialBranches:
		if value != "" {
			branches := domain.NewLocalBranchNames(strings.Split(value, " ")...)
			self.PerennialBranches = &branches
		}
	default:
		return false, nil
	}
	return true, nil
}

func EmptyPartialConfig() PartialConfig {
	return PartialConfig{} //nolint:exhaustruct
}

func PartialConfigDiff(before, after PartialConfig) ConfigDiff {
	fmt.Println("222222222222222222")
	spew.Dump(before, after)
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
	CheckLocalBranchNames(&result, KeyPerennialBranches, before.PerennialBranches, after.PerennialBranches)
	CheckStringPtr(&result, KeyCodeHostingPlatform, before.CodeHostingPlatformName, after.CodeHostingPlatformName)
	fmt.Println("3333333333333333333")
	spew.Dump(result)
	return result
}
