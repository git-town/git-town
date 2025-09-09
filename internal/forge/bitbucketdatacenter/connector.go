package bitbucketdatacenter

import (
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v21/internal/browser"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Connector provides access to the API of Bitbucket installations.
type Connector struct {
	forgedomain.Data
	log      print.Logger
	token    string
	username string
}

// type-check to ensure conformance to the Connector interface
var bbdcConnector Connector
var _ forgedomain.Connector = bbdcConnector

// NewConnector provides a Bitbucket connector instance if the current repo is hosted on Bitbucket,
// otherwise nil.
func NewConnector(args NewConnectorArgs) Connector {
	return Connector{
		Data: forgedomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
		log:      args.Log,
		token:    args.AppPassword.String(),
		username: args.UserName.String(),
	}
}

type NewConnectorArgs struct {
	AppPassword Option[forgedomain.BitbucketAppPassword]
	ForgeType   Option[forgedomain.ForgeType]
	Log         print.Logger
	RemoteURL   giturl.Parts
	UserName    Option[forgedomain.BitbucketUsername]
}

func (self Connector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	browser.Open(self.NewProposalURL(data), data.FrontendRunner)
	return nil
}

func (self Connector) DefaultProposalMessage(proposalData forgedomain.ProposalData) string {
	data := proposalData.Data()
	return forgedomain.CommitBody(data, fmt.Sprintf("%s (#%d)", data.Title, data.Number))
}

func (self Connector) NewProposalURL(data forgedomain.CreateProposalArgs) string {
	return fmt.Sprintf("%s/pull-requests?create&sourceBranch=%s&targetBranch=%s",
		self.RepositoryURL(),
		url.QueryEscape(data.Branch.String()),
		url.QueryEscape(data.ParentBranch.String()))
}

func (self Connector) OpenRepository(runner subshelldomain.Runner) error {
	browser.Open(self.RepositoryURL(), runner)
	return nil
}

func (self Connector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/projects/%s/repos/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}
