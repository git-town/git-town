package drivers

import (
	"errors"
	"fmt"
	"net/url"
)

// GitlabCodeHostingDriver provides tools for working with repositories
// on Gitlab.
type GitlabCodeHostingDriver struct{}

// GetNewPullRequestURL returns the URL of the page
// to create a new pull request on Gitlab
func (driver GitlabCodeHostingDriver) GetNewPullRequestURL(repository string, branch string, parentBranch string) string {
	query := url.Values{}
	query.Add("merge_request[source_branch]", branch)
	query.Add("merge_request[target_branch]", parentBranch)
	return fmt.Sprintf("https://gitlab.com/%s/merge_requests/new?%s", repository, query.Encode())
}

// GetRepositoryURL returns the URL of the given repository on Gitlab
func (driver GitlabCodeHostingDriver) GetRepositoryURL(repository string) string {
	return "https://gitlab.com/" + repository
}

func (driver GitlabCodeHostingDriver) GetPullRequestNumber(repository string, branch string, parentBranch string) (int, error) {
	return -1, errors.New("Unimplemented")
}

func (driver GitlabCodeHostingDriver) MergePullRequest(repository string, options MergePullRequestOptions) error {
	return errors.New("Unimplemented")
}
