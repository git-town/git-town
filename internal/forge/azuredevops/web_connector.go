package azuredevops

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

var (
	webConnector WebConnector
	_            forgedomain.Connector = webConnector
)

// WebConnector provides connectivity to Azure DevOps through the web browser.
type WebConnector struct {
	forgedomain.HostedRepoInfo
	BrowserEnabled    browserdomain.BrowserEnabled
	BrowserExecutable Option[browserdomain.BrowserExecutable]
}

func (self WebConnector) BrowseRepository(runner subshelldomain.Runner) error {
	browser.Open(self.RepositoryURL(), runner, self.BrowserExecutable)
	return nil
}

func (self WebConnector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	proposalURL := self.NewProposalURL(data)
	if self.BrowserEnabled {
		browser.Open(proposalURL, data.FrontendRunner, self.BrowserExecutable)
	} else {
		fmt.Printf(messages.BrowserOpen, proposalURL)
	}
	return nil
}

func (self WebConnector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return forgedomain.CommitBody(data, fmt.Sprintf("%s (#%d)", data.Title, data.Number))
}

func (self WebConnector) NewProposalURL(data forgedomain.CreateProposalArgs) string {
	return fmt.Sprintf("%s/pullrequestcreate?sourceRef=%s&targetRef=%s",
		self.RepositoryURL(),
		url.QueryEscape(data.Branch.String()),
		url.QueryEscape(data.ParentBranch.String()),
	)
}

func (self WebConnector) ProposalReference(data forgedomain.ProposalData) string {
	return forgedomain.ProposalReferenceFallback(data)
}

func (self WebConnector) RepositoryURL() string {
	return fmt.Sprintf("https://dev.azure.com/%s/_git/%s", self.Organization, self.Repository)
}
