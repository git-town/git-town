package github

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

// WebConnector provides connectivity to GitHub through the GitHub website.
type WebConnector struct {
	forgedomain.HostedRepoInfo
	browser Option[configdomain.Browser]
}

func (self WebConnector) BrowseRepository(runner subshelldomain.Runner) error {
	browser.Open(self.RepositoryURL(), runner, self.browser)
	return nil
}

func (self WebConnector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	browser.Open(self.NewProposalURL(data), data.FrontendRunner, self.browser)
	return nil
}

func (self WebConnector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return DefaultProposalMessage(data)
}

func (self WebConnector) NewProposalURL(data forgedomain.CreateProposalArgs) string {
	toCompare := data.Branch.String()
	if data.ParentBranch != data.MainBranch {
		toCompare = data.ParentBranch.String() + "..." + data.Branch.String()
	}
	result := fmt.Sprintf("%s/compare/%s?expand=1", self.RepositoryURL(), url.PathEscape(toCompare))
	if title, hasTitle := data.ProposalTitle.Get(); hasTitle {
		result += "&title=" + url.QueryEscape(title.String())
	}
	if body, hasBody := data.ProposalBody.Get(); hasBody {
		result += "&body=" + url.QueryEscape(body.String())
	}
	return result
}

func (self WebConnector) RepositoryURL() string {
	return RepositoryURL(self.HostnameWithStandardPort(), self.Organization, self.Repository)
}

func DefaultProposalMessage(data forgedomain.ProposalData) string {
	return forgedomain.CommitBody(data, fmt.Sprintf("%s (#%d)", data.Title, data.Number))
}
