package bitbucketdatacenter

import (
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v24/internal/browser"
	"github.com/git-town/git-town/v24/internal/browser/browserdomain"
	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/git-town/git-town/v24/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

// type check
var (
	webConnector WebConnector
	_            forgedomain.Connector = webConnector
)

// WebConnector provides connectivity to Bitbucket DataCenter through the web browser.
type WebConnector struct {
	forgedomain.HostedRepoInfo
	browserEnabled    browserdomain.BrowserEnabled
	browserExecutable Option[browserdomain.BrowserExecutable]
}

func (self WebConnector) BrowseRepository(runner subshelldomain.Runner) error {
	browser.Open(self.RepositoryURL(), runner, self.browserExecutable, self.browserEnabled)
	return nil
}

func (self WebConnector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	browser.Open(self.NewProposalURL(data), data.FrontendRunner, self.browserExecutable, self.browserEnabled)
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

func (self WebConnector) ProposalReference(data forgedomain.ProposalData) string {
	return forgedomain.ProposalReferenceFallback(data)
}

func (self WebConnector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/projects/%s/repos/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}
