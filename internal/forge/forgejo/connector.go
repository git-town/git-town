package forgejo

import (
	"fmt"
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/git-town/git-town/v21/internal/browser"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

var (
	forgejoConnector Connector
	_                forgedomain.Connector = forgejoConnector
)

// Connector provides standardized connectivity for Forgejo-based repositories
// via the Forgejo API.
type Connector struct {
	forgedomain.Data
	APIToken Option[forgedomain.ForgejoToken]
	client   *forgejo.Client
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
	APIToken  Option[forgedomain.ForgejoToken]
	Log       print.Logger
	RemoteURL giturl.Parts
}

// NewConnector provides a new connector instance.
func NewConnector(args NewConnectorArgs) (Connector, error) {
	forgejoClient, err := forgejo.NewClient("https://"+args.RemoteURL.Host, forgejo.SetToken(args.APIToken.String()))
	return Connector{
		APIToken: args.APIToken,
		Data: forgedomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
		client: forgejoClient,
		log:    args.Log,
	}, err
}
