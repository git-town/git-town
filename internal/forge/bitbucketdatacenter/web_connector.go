package bitbucketdatacenter

import (
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v22/internal/browser"
	"github.com/git-town/git-town/v22/internal/browser/browserdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// type check
var (
	webConnector WebConnector
	_            forgedomain.Connector = webConnector
)

// WebConnector provides connectivity to Bitbucket DataCenter through the web browser.
type WebConnector struct {
	forgedomain.HostedRepoInfo
	browser Option[browserdomain.Browser]
}

func (self WebConnector) BrowseRepository(runner subshelldomain.Runner) error {
	browser.Open(self.RepositoryURL(), runner, self.browser)
	return nil
}

func (self WebConnector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	proposalURL := self.NewProposalURL(data)
	if data.Headless {
		fmt.Printf(messages.BrowserOpen, proposalURL)
	} else {
		browser.Open(proposalURL, data.FrontendRunner, self.browser)
	}
	return nil
}

func (self WebConnector) DefaultProposalMessage(proposalData forgedomain.ProposalData) string {
	data := proposalData.Data()
	return forgedomain.CommitBody(data, fmt.Sprintf("%s (#%d)", data.Title, data.Number))
}

func (self WebConnector) NewProposalURL(data forgedomain.CreateProposalArgs) string {
	return fmt.Sprintf("%s/pull-requests?create&sourceBranch=%s&targetBranch=%s",
		self.RepositoryURL(),
		url.QueryEscape(data.Branch.String()),
		url.QueryEscape(data.ParentBranch.String()))
}

func (self WebConnector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/projects/%s/repos/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}
