package gitlab

import (
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v22/internal/browser"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

var (
	webConnector WebConnector
	_            forgedomain.Connector = webConnector
)

// WebConnector provides connectivity to GitLab through the browser.
type WebConnector struct {
	forgedomain.HostedRepoInfo
	browser Option[configdomain.Browser]
}

func (self WebConnector) BrowseRepository(runner subshelldomain.Runner) error {
	browser.Open(self.RepositoryURL(), runner, self.browser)
	return nil
}

func (self WebConnector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	url := self.NewProposalURL(data)
	browser.Open(url, data.FrontendRunner, self.browser)
	return nil
}

func (self WebConnector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return DefaultProposalMessage(data)
}

func (self WebConnector) NewProposalURL(data forgedomain.CreateProposalArgs) string {
	query := url.Values{}
	query.Add("merge_request[source_branch]", data.Branch.String())
	query.Add("merge_request[target_branch]", data.ParentBranch.String())
	if title, hasTitle := data.ProposalTitle.Get(); hasTitle {
		query.Add("merge_request[title]", title.String())
	}
	if body, hasBody := data.ProposalBody.Get(); hasBody {
		query.Add("merge_request[description]", body.String())
	}
	return fmt.Sprintf("%s/-/merge_requests/new?%s", self.RepositoryURL(), query.Encode())
}

func (self WebConnector) RepositoryURL() string {
	return fmt.Sprintf("%s/%s", self.baseURL(), self.projectPath())
}

func (self WebConnector) baseURL() string {
	return "https://" + self.HostnameWithStandardPort()
}

func (self WebConnector) projectPath() string {
	return fmt.Sprintf("%s/%s", self.Organization, self.Repository)
}

func DefaultProposalMessage(data forgedomain.ProposalData) string {
	return forgedomain.CommitBody(data, fmt.Sprintf("%s (!%d)", data.Title, data.Number))
}
