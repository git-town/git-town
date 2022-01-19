package drivers_test

type mockConfig struct {
	manualHostName        string
	codeHostingDriverName string
	giteaToken            string
	gitHubToken           string
	mainBranch            string
	remoteOriginURL       string
}

func (mc mockConfig) CodeHostingDriverName() string {
	return mc.codeHostingDriverName
}

func (mc mockConfig) CodeHostingOriginHostname() string {
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

func (mc mockConfig) RemoteOriginURL() string {
	return mc.remoteOriginURL
}
