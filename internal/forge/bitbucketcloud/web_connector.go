package bitbucketcloud

import (
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v22/internal/browser"
	"github.com/git-town/git-town/v22/internal/browser/browserdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

var (
	webConnector WebConnector
	_            forgedomain.Connector = webConnector
)

// WebConnector provides connectivity to Bitbucket Cloud through the web browser.
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

func (self WebConnector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return forgedomain.CommitBody(data, fmt.Sprintf("%s (#%d)", data.Title, data.Number))
}

func (self WebConnector) NewProposalURL(data forgedomain.CreateProposalArgs) string {
	return fmt.Sprintf("%s/pull-requests/new?source=%s&dest=%s%%2F%s%%3A%s",
		self.RepositoryURL(),
		url.QueryEscape(data.Branch.String()),
		url.QueryEscape(self.Organization),
		url.QueryEscape(self.Repository),
		url.QueryEscape(data.ParentBranch.String()))
}

func (self WebConnector) ProposalReference(data forgedomain.ProposalData) string {
	return forgedomain.ProposalReferenceFallback(data)
}

func (self WebConnector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}
