package config

type Hosting struct {
	gitConfig *gitConfig
}

// HostingService provides the name of the code hosting driver to use.
func (h *Hosting) HostingService() string {
	return h.gitConfig.localOrGlobalConfigValue("git-town.code-hosting-driver")
}

// OriginOverride provides the override for the origin hostname from the Git Town configuration.
func (h *Hosting) OriginOverride() string {
	return h.gitConfig.localConfigValue("git-town.code-hosting-origin-hostname")
}

// GitHubToken provides the content of the GitHub API token stored in the local or global Git Town configuration.
func (h *Hosting) GitHubToken() string {
	return h.gitConfig.localOrGlobalConfigValue("git-town.github-token")
}

// GitLabToken provides the content of the GitLab API token stored in the local or global Git Town configuration.
func (h *Hosting) GitLabToken() string {
	return h.gitConfig.localOrGlobalConfigValue("git-town.gitlab-token")
}

// GiteaToken provides the content of the Gitea API token stored in the local or global Git Town configuration.
func (h *Hosting) GiteaToken() string {
	return h.gitConfig.localOrGlobalConfigValue("git-town.gitea-token")
}

// SetCodeHostingDriver sets the "github.code-hosting-driver" setting.
func (h *Hosting) SetCodeHostingDriver(value string) error {
	const key = "git-town.code-hosting-driver"
	h.gitConfig.localConfigCache[key] = value
	_, err := h.gitConfig.shell.Run("git", "config", key, value)
	return err
}

// SetCodeHostingOriginHostname sets the "github.code-hosting-driver" setting.
func (h *Hosting) SetCodeHostingOriginHostname(value string) error {
	const key = "git-town.code-hosting-origin-hostname"
	h.gitConfig.localConfigCache[key] = value
	_, err := h.gitConfig.shell.Run("git", "config", key, value)
	return err
}

// SetTestOrigin sets the origin to be used for testing.
func (h *Hosting) SetTestOrigin(value string) error {
	_, err := h.gitConfig.SetLocalConfigValue("git-town.testing.remote-url", value)
	return err
}
