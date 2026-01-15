package azuredevops

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/giturl"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// Detect indicates whether the current repository is hosted on Azure DevOps.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "dev.azure.com" || remoteURL.Host == "ssh.dev.azure.com"
}

type NewConnectorArgs struct {
	Browser   Option[configdomain.Browser]
	RemoteURL giturl.Parts
}

// NewConnector provides the correct connector for talking to Azure DevOps.
func NewConnector(args NewConnectorArgs) WebConnector {
	return WebConnector{
		HostedRepoInfo: forgedomain.HostedRepoInfo{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
		browser: args.Browser,
	}
}
