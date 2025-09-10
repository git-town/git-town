package gitlab

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

// WebConnector provides connectivity to GitLab without authentication data.
type WebConnector struct {
	forgedomain.Data
}

func (self WebConnector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	url := self.NewProposalURL(data)
	browser.Open(url, data.FrontendRunner)
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

func (self WebConnector) OpenRepository(runner subshelldomain.Runner) error {
	browser.Open(self.RepositoryURL(), runner)
	return nil
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
