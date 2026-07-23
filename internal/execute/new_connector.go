package execute

import (
	"github.com/git-town/git-town/v24/internal/cli/print"
	"github.com/git-town/git-town/v24/internal/forge"
	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/git-town/git-town/v24/internal/git/giturl"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

// NewConnector provides the forge connector for this repo and the given remote URL.
func (self OpenRepoResult) NewConnector(remoteURL Option[giturl.Parts]) (Option[forgedomain.Connector], Option[forgedomain.DetectedForgeType], error) {
	config := self.UnvalidatedConfig.NormalConfig
	return forge.NewConnector(forge.NewConnectorArgs{
		Backend:             self.Backend,
		BitbucketAPIToken:   config.BitbucketAPIToken,
		BitbucketUsername:   config.BitbucketUsername,
		BrowserEnabled:      config.BrowserEnabled,
		BrowserExecutable:   config.BrowserExecutable,
		ConfigDir:           self.ConfigDir,
		ForgeType:           config.ForgeType,
		ForgejoToken:        config.ForgejoToken,
		Frontend:            self.Frontend,
		GiteaToken:          config.GiteaToken,
		GithubConnectorType: config.GithubConnectorType,
		GithubToken:         config.GithubToken,
		GitlabConnectorType: config.GitlabConnectorType,
		GitlabToken:         config.GitlabToken,
		Log:                 print.Logger{},
		RemoteURL:           remoteURL,
	})
}
