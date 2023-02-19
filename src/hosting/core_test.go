package hosting_test

type mockRepoConfig struct {
	giteaToken     string
	gitHubToken    string
	gitLabToken    string
	hostingService string
	mainBranch     string
	originOverride string
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

func (mc mockRepoConfig) HostingService() string {
	return mc.hostingService
}

func (mc mockRepoConfig) MainBranch() string {
	return mc.mainBranch
}

func (mc mockRepoConfig) OriginOverride() string {
	return mc.originOverride
}

func (mc mockRepoConfig) OriginURL() string {
	return mc.originURL
}
