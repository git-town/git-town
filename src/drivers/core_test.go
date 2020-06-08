package drivers_test

type mockConfig struct {
	configuredHostName    string
	codeHostingDriverName string
	giteaToken            string
	gitHubToken           string
	mainBranch            string
	remoteOriginURL       string
}

func (mc mockConfig) GetCodeHostingOriginHostname() string {
	return mc.configuredHostName
}

func (mc mockConfig) GetCodeHostingDriverName() string {
	return mc.codeHostingDriverName
}

func (mc mockConfig) GetGitHubToken() string {
	return mc.gitHubToken
}

func (mc mockConfig) GetGiteaToken() string {
	return mc.giteaToken
}

func (mc mockConfig) GetMainBranch() string {
	return mc.mainBranch
}

func (mc mockConfig) GetRemoteOriginURL() string {
	return mc.remoteOriginURL
}
