package hosting_test

type mockConfig struct {
	giteaToken     string
	gitHubToken    string
	gitLabToken    string
	hostingService string
	mainBranch     string
	originOverride string
	originURL      string
}

func (mc mockConfig) GiteaToken() string {
	return mc.giteaToken
}

func (mc mockConfig) GitHubToken() string {
	return mc.gitHubToken
}

func (mc mockConfig) GitLabToken() string {
	return mc.gitLabToken
}

func (mc mockConfig) HostingService() string {
	return mc.hostingService
}

func (mc mockConfig) MainBranch() string {
	return mc.mainBranch
}

func (mc mockConfig) OriginOverride() string {
	return mc.originOverride
}

func (mc mockConfig) OriginURL() string {
	return mc.originURL
}
