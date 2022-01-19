package drivers_test

type mockConfig struct {
	manualHostName        string
	codeHostingDriverName string
	giteaToken            string
	gitHubToken           string
	mainBranch            string
	remoteOriginURL       string
}

func (mc mockConfig) CodeHostingOriginHostname() string {
	return mc.manualHostName
}

func (mc mockConfig) CodeHostingDriverName() string {
	return mc.codeHostingDriverName
}

func (mc mockConfig) GitHubToken() string {
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
