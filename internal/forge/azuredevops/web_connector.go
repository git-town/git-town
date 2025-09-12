package azuredevops

import (
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v21/internal/browser"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
)

var (
	webConnector WebConnector
	_            forgedomain.Connector = webConnector
)

// WebConnector provides connectivity to Azure DevOps through the web browser.
type WebConnector struct {
	forgedomain.HostedRepoInfo
}

func (self WebConnector) BrowseRepository(runner subshelldomain.Runner) error {
	browser.Open(self.RepositoryURL(), runner)
	return nil
}

func (self WebConnector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	browser.Open(self.NewProposalURL(data), data.FrontendRunner)
	return nil
}

func (self WebConnector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return forgedomain.CommitBody(data, fmt.Sprintf("%s (#%d)", data.Title, data.Number))
}

func (self WebConnector) NewProposalURL(data forgedomain.CreateProposalArgs) string {
	// https://dev.azure.com/kevingoslar/tikibase/_git/tikibase/pullrequestcreate?sourceRef=kg-test&targetRef=main
	return fmt.Sprintf("%s/_git/%s/pullrequestcreate?sourceRef=%s&targetRef=%s",
		self.RepositoryURL(),
		url.QueryEscape(self.Repository),
		url.QueryEscape(data.Branch.String()),
		url.QueryEscape(data.ParentBranch.String()),
	)
}

func (self WebConnector) RepositoryURL() string {
	// https://dev.azure.com/kevingoslar/_git/tikibase
	return fmt.Sprintf("https://dev.azure.com/%s/_git/%s", self.Organization, self.Repository)
}
