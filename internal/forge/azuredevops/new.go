package azuredevops

import (
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Detect indicates whether the current repository is hosted on Bitbucket Cloud.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "dev.azure.com"
}

type NewConnectorArgs struct {
	ProposalOverride Option[forgedomain.ProposalOverride]
	RemoteURL        giturl.Parts
}

// NewConnector provides the correct connector for talking to Bitbucket Cloud.
func NewConnector(args NewConnectorArgs) forgedomain.Connector { //nolint: ireturn
	return WebConnector{
		Data: forgedomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
	}
}
