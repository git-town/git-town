package forgejo

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

// WebConnector provides connectivity to Forgejo through the web browser.
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
	toCompare := data.ParentBranch.String() + "..." + data.Branch.String()
	return fmt.Sprintf("%s/compare/%s", self.RepositoryURL(), url.PathEscape(toCompare))
}

func (self WebConnector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}
