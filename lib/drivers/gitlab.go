package drivers

import (
	"fmt"
	"net/url"
)

type GitlabCodeHostingDriver struct{}

func (driver GitlabCodeHostingDriver) GetNewPullRequestURL(repository string, branch string, parentBranch string) string {
	query := url.Values{}
	query.Add("merge_request[source_branch]", branch)
	query.Add("merge_request[target_branch]", parentBranch)
	return fmt.Sprintf("https://gitlab.com/%s/merge_requests/new?%s", repository, query.Encode())
}

func (driver GitlabCodeHostingDriver) GetRepositoryURL(repository string) string {
	return "https://gitlab.com/" + repository
}
