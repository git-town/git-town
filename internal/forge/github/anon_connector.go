package github

import (
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v21/internal/browser"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
)

var githubAnonConnector AnonConnector
var _ forgedomain.Connector = githubAnonConnector

// AnonConnector provides connectivity to GitHub without authentication.
type AnonConnector struct {
	forgedomain.Data
}

func (self AnonConnector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	browser.Open(self.NewProposalURL(data), data.FrontendRunner)
	return nil
}

func (self AnonConnector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return DefaultProposalMessage(data)
}

func (self AnonConnector) NewProposalURL(data forgedomain.CreateProposalArgs) string {
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

func (self AnonConnector) OpenRepository(runner subshelldomain.Runner) error {
	browser.Open(self.RepositoryURL(), runner)
	return nil
}
func (self AnonConnector) RepositoryURL() string {
	return RepositoryURL(self.HostnameWithStandardPort(), self.Organization, self.Repository)
}

func DefaultProposalMessage(data forgedomain.ProposalData) string {
	return forgedomain.CommitBody(data, fmt.Sprintf("%s (#%d)", data.Title, data.Number))
}
