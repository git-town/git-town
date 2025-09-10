package gitlab

import (
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v21/internal/browser"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
)

var (
	gitlabAnonConnector AnonConnector
	_                   forgedomain.Connector = gitlabAnonConnector
)

// AnonConnector provides connectivity to GitLab without authentication data.
type AnonConnector struct {
	forgedomain.Data
}

func (self AnonConnector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	url := self.NewProposalURL(data)
	browser.Open(url, data.FrontendRunner)
	return nil
}

func (self AnonConnector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return DefaultProposalMessage(data)
}

func (self AnonConnector) NewProposalURL(data forgedomain.CreateProposalArgs) string {
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

func (self AnonConnector) OpenRepository(runner subshelldomain.Runner) error {
	browser.Open(self.RepositoryURL(), runner)
	return nil
}

func (self AnonConnector) RepositoryURL() string {
	return fmt.Sprintf("%s/%s", self.baseURL(), self.projectPath())
}

func (self AnonConnector) baseURL() string {
	return "https://" + self.HostnameWithStandardPort()
}

func (self AnonConnector) projectPath() string {
	return fmt.Sprintf("%s/%s", self.Organization, self.Repository)
}

func DefaultProposalMessage(data forgedomain.ProposalData) string {
	return forgedomain.CommitBody(data, fmt.Sprintf("%s (!%d)", data.Title, data.Number))
}
