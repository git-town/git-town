package hosting_test

type mockConfig struct {
	manualHostName string
	driverName     string
	giteaToken     string
	gitHubToken    string
	mainBranch     string
	originURL      string
}

func (mc mockConfig) HostingService() string {
	return mc.driverName
}

func (mc mockConfig) OriginOverride() string {
	return mc.manualHostName
}

func (mc mockConfig) GitHubToken() string {
	return mc.gitHubToken
}

func (mc mockConfig) GiteaToken() string {
	return mc.giteaToken
}

func (mc mockConfig) MainBranch() string {
	return mc.mainBranch
}

func (mc mockConfig) OriginURL() string {
	return mc.originURL
}
