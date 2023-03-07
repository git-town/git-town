package hosting_test

import (
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/giturl"
)

type mockRepoConfig struct {
	giteaToken     string                `exhaustruct:"optional"`
	gitHubToken    string                `exhaustruct:"optional"`
	gitLabToken    string                `exhaustruct:"optional"`
	hostingService config.HostingService `exhaustruct:"optional"`
	mainBranch     string                `exhaustruct:"optional"`
	originOverride string                `exhaustruct:"optional"`
	originURL      string
}

func (mc mockRepoConfig) GiteaToken() string {
	return mc.giteaToken
}

func (mc mockRepoConfig) GitHubToken() string {
	return mc.gitHubToken
}

func (mc mockRepoConfig) GitLabToken() string {
	return mc.gitLabToken
}

func (mc mockRepoConfig) HostingService() (config.HostingService, error) {
	return mc.hostingService, nil
}

func (mc mockRepoConfig) MainBranch() string {
	return mc.mainBranch
}

func (mc mockRepoConfig) OriginOverride() string {
	return mc.originOverride
}

func (mc mockRepoConfig) OriginURL() *giturl.Parts {
	url := giturl.Parse(mc.originURL)
	if mc.originOverride != "" {
		url.Host = mc.originOverride
	}
	return url
}
