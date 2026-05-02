package gitlab

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

// WebConnector provides connectivity to GitLab through the browser.
type WebConnector struct {
	forgedomain.HostedRepoInfo
	browserEnabled    browserdomain.BrowserEnabled
	browserExecutable Option[browserdomain.BrowserExecutable]
}

func (self WebConnector) BrowseRepository(runner subshelldomain.Runner) error {
	browser.Open(self.RepositoryURL(), runner, self.browserExecutable)
	return nil
}

func (self WebConnector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	proposalURL := self.NewProposalURL(data)
	if self.browserEnabled {
		browser.Open(proposalURL, data.FrontendRunner, self.browserExecutable)
	} else {
		fmt.Printf(messages.BrowserOpen, proposalURL)
	}
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

func (self WebConnector) ProposalReference(data forgedomain.ProposalData) string {
	return ProposalReference(data)
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

func ProposalReference(data forgedomain.ProposalData) string {
	if data.Number.Int() > 0 {
		return "!" + data.Number.String() + "+"
	}
	return forgedomain.ProposalReferenceFallback(data)
}
