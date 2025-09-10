package gitlab

import (
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type NewConnectorArgs struct {
	APIToken  Option[forgedomain.GitLabToken]
	Log       print.Logger
	RemoteURL giturl.Parts
}

func NewConnector(args NewConnectorArgs) (forgedomain.Connector, error) {
	anonConnector := AnonConnector{
		Data: forgedomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
	}
	apiToken, hasAPIToken := args.APIToken.Get()
	if !hasAPIToken {
		return anonConnector, nil
	}
	client, err := gitlab.NewClient(apiToken.String(), gitlab.WithBaseURL(anonConnector.baseURL()))
	if err != nil {
		return anonConnector, err
	}
	return AuthConnector{
		AnonConnector: anonConnector,
		APIToken:      apiToken,
		client:        client,
		log:           args.Log,
	}, nil
}
