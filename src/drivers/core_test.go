package drivers_test

type mockConfig struct {
	configuredHostName    string
	codeHostingDriverName string
	giteaToken            string
	gitHubToken           string
	remoteOriginURL       string
}

func (mgc mockConfig) GetCodeHostingOriginHostname() string {
	return mgc.configuredHostName
}

func (mgc mockConfig) GetCodeHostingDriverName() string {
	return mgc.codeHostingDriverName
}

func (mgc mockConfig) GetGitHubToken() string {
	return mgc.gitHubToken
}

func (mgc mockConfig) GetGiteaToken() string {
	return mgc.giteaToken
}

func (mgc mockConfig) GetRemoteOriginURL() string {
	return mgc.remoteOriginURL
}
