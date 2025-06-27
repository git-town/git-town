package gitlab

import (
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type Data struct {
	forgedomain.Data
	APIToken Option[configdomain.GitLabToken]
}

func (self Data) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return DefaultProposalMessage(data)
}

func (self Data) NewProposalURL(data forgedomain.CreateProposalArgs) string {
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

func (self Data) RepositoryURL() string {
	return fmt.Sprintf("%s/%s", self.baseURL(), self.projectPath())
}

func (self Data) baseURL() string {
	return "https://" + self.HostnameWithStandardPort()
}

func (self Data) projectPath() string {
	return fmt.Sprintf("%s/%s", self.Organization, self.Repository)
}

func DefaultProposalMessage(data forgedomain.ProposalData) string {
	return forgedomain.CommitBody(data, fmt.Sprintf("%s (!%d)", data.Title, data.Number))
}
