package azuredevops

import (
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Detect indicates whether the current repository is hosted on Azure DevOps.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "dev.azure.com" || remoteURL.Host == "ssh.dev.azure.com"
}

type NewConnectorArgs struct {
	ProposalOverride Option[forgedomain.ProposalOverride]
	RemoteURL        giturl.Parts
}

// NewConnector provides the correct connector for talking to Azure DevOps.
func NewConnector(args NewConnectorArgs) WebConnector {
	return WebConnector{
		Data: forgedomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
	}
}
