package gitea

import (
	"context"
	"fmt"
	"net/url"

	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v21/internal/browser"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"golang.org/x/oauth2"
)

var giteaConnector Connector
var _ forgedomain.Connector = giteaConnector

type Connector struct {
	forgedomain.Data
	APIToken Option[forgedomain.GiteaToken]
	client   *gitea.Client
	log      print.Logger
}

func (self Connector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	toCompare := data.ParentBranch.String() + "..." + data.Branch.String()
	url := fmt.Sprintf("%s/compare/%s", self.RepositoryURL(), url.PathEscape(toCompare))
	browser.Open(url, data.FrontendRunner)
	return nil
}

func (self Connector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return forgedomain.CommitBody(data, fmt.Sprintf("%s (#%d)", data.Title, data.Number))
}

func (self Connector) OpenRepository(runner subshelldomain.Runner) error {
	browser.Open(self.RepositoryURL(), runner)
	return nil
}

func (self Connector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}

type NewConnectorArgs struct {
	APIToken  Option[forgedomain.GiteaToken]
	Log       print.Logger
	RemoteURL giturl.Parts
}

// NewGiteaConfig provides Gitea configuration data if the current repo is hosted on Gitea,
// otherwise nil.
func NewConnector(args NewConnectorArgs) Connector {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: args.APIToken.String()})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	giteaClient := gitea.NewClientWithHTTP("https://"+args.RemoteURL.Host, httpClient)
	return Connector{
		APIToken: args.APIToken,
		Data: forgedomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
		client: giteaClient,
		log:    args.Log,
	}
}
