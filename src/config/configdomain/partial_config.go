package configdomain

import (
	"fmt"

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
	default:
		return false, nil
	}
	return true, nil
}

func EmptyPartialConfig() PartialConfig {
	return PartialConfig{} //nolint:exhaustruct
}
